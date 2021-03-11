# go-socks
[![CI](https://github.com/nwtgck/go-socks/actions/workflows/ci.yml/badge.svg)](https://github.com/nwtgck/go-socks/actions/workflows/ci.yml)

SOCKS4, SOCKS4a and SOCKS5 proxy server in Go

## Thanks
This project is a fork of [armon/go-socks5](https://github.com/armon/go-socks5). The SOCKS5 implementation was written in that project. Thanks!


## Feature

The package has the following features:
* "No Auth" mode
* User/Password authentication
* Support for the CONNECT command
* Rules to do granular filtering of commands
* Custom DNS resolution
* Unit tests

## TODO

The package still needs the following:
* Support for the BIND command
* Support for the ASSOCIATE command


## Example

Below is a simple example of usage

```go
socksConf := &socks.Config{}
socksServer, err := socks.New(socksConf)
if err != nil {
    panic(err)
}

l, err := net.Listen("tcp", "127.0.0.1:1080")
if err != nil {
    panic(err)
}
for {
    conn, err := l.Accept()
    if err != nil {
        panic(err)
    }
    fmt.Println("accepted")
    go socksServer.ServeConn(conn)
}
```
