/*
 * Copyright 2023 the go-netty project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"fmt"
	"net/http"
	"time"

	nettyws "github.com/go-netty/go-netty-ws"
)

func main() {
	header := http.Header{}
	header.Add("server-time", time.Now().String())

	// create websocket instance
	var ws = nettyws.NewWebsocket(nettyws.WithServerHeader(header))

	// setup OnOpen handler
	ws.OnOpen = func(conn nettyws.Conn) {
		fmt.Println("OnOpen: ", conn.RemoteAddr(), ", header: ", conn.Header())
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
}
