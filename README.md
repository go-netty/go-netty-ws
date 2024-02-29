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

## API overview
```
type Websocket
    func NewWebsocket(options ...Option) *Websocket
    func (ws *Websocket) Close() error
    func (ws *Websocket) Listen(addr string) error
    func (ws *Websocket) Open(addr string) (Conn, error)
    func (ws *Websocket) ServeHTTP(w http.ResponseWriter, r *http.Request)
    func (ws *Websocket) UpgradeHTTP(w http.ResponseWriter, r *http.Request) (Conn, error)

type Option
    func WithAsyncWrite(writeQueueSize int, writeForever bool) Option
    func WithBinary() Option
    func WithBufferSize(readBufferSize, writeBufferSize int) Option
    func WithCompress(compressLevel int, compressThreshold int64) Option
    func WithClientHeader(header http.Header) Option
    func WithDialer(dialer Dialer) Option
    func WithMaxFrameSize(maxFrameSize int64) Option
    func WithNoDelay(noDelay bool) Option
    func WithServerHeader(header http.Header) Option
    func WithServeMux(serveMux *http.ServeMux) Option
    func WithServeTLS(tls *tls.Config) Option
    func WithValidUTF8() Option
```

## Easy to use

> Note: `nettyws` does not support mixed text messages and binary messages, use the `WithBinary` option to switch to binary message mode.

### server :
```go
// create websocket instance
var ws = nettyws.NewWebsocket()

// setup OnOpen handler
ws.OnOpen = func(conn nettyws.Conn) {
    fmt.Println("OnOpen: ", conn.RemoteAddr())
}

// setup OnData handler
ws.OnData = func(conn nettyws.Conn, data []byte) {
    fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data))
    conn.Write(data)
}

// setup OnClose handler
ws.OnClose = func(conn nettyws.Conn, err error) {
    fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
}

fmt.Println("listening websocket connections ....")

// listen websocket server
if err := ws.Listen("ws://127.0.0.1:9527/ws"); nil != err {
    panic(err)
}
```

### client :
```go
// create websocket instance
var ws = nettyws.NewWebsocket()

// setup OnOpen handler
ws.OnOpen = func(conn nettyws.Conn) {
    fmt.Println("OnOpen: ", conn.RemoteAddr())
    conn.Write([]byte("hello world"))
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
if _, err := ws.Open("ws://127.0.0.1:9527/ws"); nil != err {
    panic(err)
}
```

### upgrade from http server:
```go
// create websocket instance
var ws = nettyws.NewWebsocket()

// setup OnOpen handler
ws.OnOpen = func(conn nettyws.Conn) {
    fmt.Println("OnOpen: ", conn.RemoteAddr())
}

// setup OnData handler
ws.OnData = func(conn nettyws.Conn, data []byte) {
    fmt.Println("OnData: ", conn.RemoteAddr(), ", message: ", string(data))
    conn.Write(data)
}

// setup OnClose handler
ws.OnClose = func(conn nettyws.Conn, err error) {
    fmt.Println("OnClose: ", conn.RemoteAddr(), ", error: ", err)
}

fmt.Println("upgrade websocket connections ....")

// upgrade websocket connection from http server
http.Handle("/ws", ws)

// listen http server
if err := http.ListenAndServe(":9527", nil); nil != err {
    panic(err)
}
```

## Associated
* https://github.com/go-netty/go-netty
* https://github.com/go-netty/go-netty-transport
* https://github.com/gobwas/ws
* https://github.com/lesismal/go-websocket-benchmark
* https://github.com/crossbario/autobahn-testsuite