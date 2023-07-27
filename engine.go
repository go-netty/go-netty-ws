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

package nettyws

import (
	"strconv"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty-transport/websocket"
)

// ClosedError returned when peer has closed the connection with appropriate
// code and a textual reason.
type ClosedError struct {
	Code   int
	Reason string
}

// Error implements error interface.
func (err ClosedError) Error() string {
	return "ws closed: " + strconv.FormatUint(uint64(err.Code), 10) + " " + err.Reason
}

var defaultEngine = netty.NewBootstrap(
	netty.WithTransport(websocket.New()),
	netty.WithChannel(netty.NewChannel()),
	netty.WithChannelHolder(nil),
	netty.WithClientInitializer(makeInitializer(true)),
	netty.WithChildInitializer(makeInitializer(false)),
)

func makeInitializer(client bool) netty.ChannelInitializer {
	return func(channel netty.Channel) {
		ws := channel.Attachment().(*Websocket)
		channel.Pipeline().
			AddLast(ws.holder).
			AddLast(newConn(ws, channel, client))
	}
}
