package main

import (
	"compress/flate"
	"fmt"
	"log"
	"time"

	nettyws "github.com/go-netty/go-netty-ws"
)

const remoteAddr = "127.0.0.1:9001"

func main() {
	const count = 517
	for i := 1; i <= count; i++ {
		testCase(i)
	}
	updateReports()
}

func testCase(id int) {
	var url = fmt.Sprintf("ws://%s/runCase?case=%d&agent=nettyws/client", remoteAddr, id)
	var onexit = make(chan struct{})

	ws := nettyws.NewWebsocket(
		nettyws.WithBufferSize(16*1024, 0),
		nettyws.WithMaxFrameSize(32*1024*1024),
		nettyws.WithCompress(flate.BestSpeed, 512),
		nettyws.WithValidUTF8(),
	)

	ws.OnOpen = func(conn nettyws.Conn) {
		conn.SetDeadline(time.Now().Add(30 * time.Second))
	}

	ws.OnData = func(conn nettyws.Conn, data []byte) {
		conn.Write(data)
	}

	ws.OnClose = func(conn nettyws.Conn, err error) {
		onexit <- struct{}{}
	}

	if _, err := ws.Open(url); nil != err {
		log.Println("ws.Open(", url, ") =>", err)
	}

	<-onexit
}

func updateReports() {
	var url = fmt.Sprintf("ws://%s/updateReports?agent=nettyws/client", remoteAddr)
	var onexit = make(chan struct{})

	var ws = nettyws.NewWebsocket(
		nettyws.WithCompress(flate.BestSpeed, 512),
		nettyws.WithValidUTF8(),
	)

	ws.OnOpen = func(conn nettyws.Conn) {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
	}

	ws.OnData = func(conn nettyws.Conn, data []byte) {
	}

	ws.OnClose = func(conn nettyws.Conn, err error) {
		onexit <- struct{}{}
	}

	if _, err := ws.Open(url); nil != err {
		log.Println("ws.Open(", url, ") =>", err)
	}

	<-onexit
}
