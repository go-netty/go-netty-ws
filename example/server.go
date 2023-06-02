package main

import (
	"fmt"

	nettyws "github.com/go-netty/go-netty-ws"
)

func main() {
	// create websocket instance
	var ws = nettyws.NewWebsocket("ws://127.0.0.1:9527")

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
	if err := ws.Listen(); nil != err {
		panic(err)
	}
}
