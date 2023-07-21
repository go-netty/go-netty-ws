package main

import (
	nettyws "github.com/go-netty/go-netty-ws"
)

func main() {
	ws := nettyws.NewWebsocket(nettyws.WithBufferSize(4096, 0))

	ws.OnData = func(conn nettyws.Conn, data []byte) {
		_ = conn.Write(data)
	}

	if err := ws.Listen(":8000"); nil != err {
		panic(err)
	}
}
