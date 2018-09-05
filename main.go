package main

import (
	"bufio"
	"crypto/md5"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var hashes = make(map[string]string)

func sendFile(conn net.Conn, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return err
	}

	hash := fmt.Sprintf("%x", md5.Sum([]byte(filename)))
	hashes[hash] = filename

	fmt.Fprintf(conn, "open\n")
	fmt.Fprintf(conn, "display-name: %v\n", filepath.Base(filename))
	//fmt.Fprintf(conn, "real-path: %v\n", filename)
	fmt.Fprintf(conn, "data-on-save: yes\n")
	fmt.Fprintf(conn, "re-activate: yes\n")
	fmt.Fprintf(conn, "token: %v\n", hash)
	fmt.Fprintf(conn, "data: %v\n", st.Size())
	if _, err := io.Copy(conn, f); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(conn, "\n.\n")

	return nil
}

func handleCmds(buf *bufio.Reader) error {
	b, _, err := buf.ReadLine()
	if err != nil {
		return err
	}
	cmd := strings.TrimSpace(string(b))
	log.Println(cmd)
	switch cmd {
	case "close":
		return nil
	case "save":
		var token string
		var size int64
		for {
			b, _, err = buf.ReadLine()
			if err != nil {
				return err
			}
			cmd = strings.TrimSpace(string(b))
			log.Printf(cmd)
			if strings.HasPrefix(cmd, "token:") {
				token = strings.TrimSpace(cmd[6:])
			} else if strings.HasPrefix(cmd, "data:") {
				size, err = strconv.ParseInt(cmd[6:], 10, 64)
				if err != nil {
					return err
				}
				break
			}
		}
		f, err := ioutil.TempFile("", "")
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.CopyN(f, buf, int64(size))
		if err != nil {
			return err
		}
		f.Close()
		if filename, ok := hashes[token]; ok {
			return os.Rename(f.Name(), filename)
		}
		return errors.New("unknown token: " + token)
	}
	return nil
}

func main() {
	var hostname string
	var port int

	flag.StringVar(&hostname, "hostname", "localhost", "hostname")
	flag.IntVar(&port, "port", 52698, "port")
	flag.Parse()

	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", hostname, port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for _, f := range flag.Args() {
		if err = sendFile(conn, f); err != nil {
			log.Fatal(err)
		}
	}

	buf := bufio.NewReader(conn)

	b, _, err := buf.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(strings.TrimSpace(string(b)))
	for {
		handleCmds(buf)
	}
}
