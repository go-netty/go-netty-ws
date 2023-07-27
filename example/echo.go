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
