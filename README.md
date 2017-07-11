# gomate

Edit files from an ssh session in TextMate/VSCode

## Usage

logged-in server with port-forwarding.

```
$ ssh -R 52698:127.0.0.1:52698 user@example.org rmate /path/to/the/file.txt
```

then, run gomate on the server

```
$ gomate /path/to/the/file.html /path/to/the/file.css
```

## Installation

```
go get github.com/mattn/gomate
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
