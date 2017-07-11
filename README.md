# gomate

Edit files from an ssh session in TextMate/VSCode

## Usage

logged-in server and run gomate with port-forwarding.

```
$ ssh -R 52698:127.0.0.1:52698 user@example.org gomate /path/to/the/file.txt
```

## Installation

Install [Remote VSCode](https://marketplace.visualstudio.com/items?itemName=rafaelmaiolla.remote-vscode)

```
go get github.com/mattn/gomate
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
