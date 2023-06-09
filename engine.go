package nettyws

import (
	"strconv"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty-transport/websocket"
	"github.com/go-netty/go-netty/codec/frame"
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
			AddLast(frame.PacketCodec(1024)).
			AddLast(NewConn(channel, client, ws.OnOpen, ws.OnData, ws.OnClose))
	}
}
