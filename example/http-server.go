package main

import (
	"fmt"
	"net/http"

	nettyws "github.com/go-netty/go-netty-ws"
)

func main() {

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
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		ws.UpgradeHTTP(writer, request)
	})

	// listen http server
	if err := http.ListenAndServe(":9527", serveMux); nil != err {
		panic(err)
	}
}
