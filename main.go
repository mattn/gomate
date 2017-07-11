package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

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

	fmt.Fprintf(conn, "open\n")
	fmt.Fprintf(conn, "display-name: %v\n", filename)
	fmt.Fprintf(conn, "real-path: %v\n", filename)
	fmt.Fprintf(conn, "data-on-save: yes\n")
	fmt.Fprintf(conn, "re-activate: yes\n")
	fmt.Fprintf(conn, "token: %v\n", filename)
	fmt.Fprintf(conn, "data: %v\n", st.Size())
	io.Copy(conn, f)
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
		var filename string
		var size int64
		for {
			b, _, err = buf.ReadLine()
			if err != nil {
				return err
			}
			cmd = string(b)
			if strings.HasPrefix(cmd, "token:") {
				filename = strings.TrimSpace(cmd[6:])
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
		return os.Rename(f.Name(), filename)
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
	log.Println(string(b))
	for {
		handleCmds(buf)
	}
}
