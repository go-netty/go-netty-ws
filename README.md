# go-netty-ws

[![GoDoc][1]][2] [![license-Apache 2][3]][4]

<!--[![Downloads][7]][8]-->

[1]: https://godoc.org/github.com/go-netty/go-netty-ws?status.svg
[2]: https://godoc.org/github.com/go-netty/go-netty-ws
[3]: https://img.shields.io/badge/license-Apache%202-blue.svg
[4]: LICENSE

An Websocket server & client written by [go-netty](https://github.com/go-netty/go-netty)

## Install
```
go get github.com/go-netty/go-netty-ws@latest
```

## Example

server :
```go
// create websocket instance
var ws = nettyws.NewWebsocket("ws://127.0.0.1:9527")

// setup OnOpen handler
ws.OnOpen = func(conn nettyws.Conn) {
    fmt.Println("OnOpen: ", conn.RemoteAddr())
}

// setup OnData handler
ws.OnData = func(conn nettyws.Conn, data []byte) {
    fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data))
}

// setup OnClose handler
ws.OnClose = func(conn nettyws.Conn, err error) {
    fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
}

fmt.Println("listening websocket connections ....")

// listening websocket server
if err := ws.Listen(); nil != err {
    panic(err)
}

```

client :
```go
// create websocket instance
var ws = nettyws.NewWebsocket("ws://127.0.0.1:9527")

// setup OnOpen handler
ws.OnOpen = func(conn nettyws.Conn) {
    fmt.Println("OnOpen: ", conn.RemoteAddr())
}

// setup OnData handler
ws.OnData = func(conn nettyws.Conn, data []byte) {
    fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data))
}

// setup OnClose handler
ws.OnClose = func(conn nettyws.Conn, err error) {
    fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
}

fmt.Println("open websocket connection ...")

// connect to websocket server
if err := ws.Open(); nil != err {
    panic(err)
}
```