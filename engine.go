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
	netty.WithClientInitializer(clientInitializer),
	netty.WithChildInitializer(childInitializer),
)

func clientInitializer(channel netty.Channel) {
	ws := channel.Attachment().(*Websocket)
	channel.Pipeline().
		AddLast(ws.holder).
		AddLast(frame.PacketCodec(1024)).
		AddLast(NewConn(channel, true, ws.OnOpen, ws.OnData, ws.OnClose))
}

func childInitializer(channel netty.Channel) {
	ws := channel.Attachment().(*Websocket)
	channel.Pipeline().
		AddLast(ws.holder).
		AddLast(frame.PacketCodec(1024)).
		AddLast(NewConn(channel, false, ws.OnOpen, ws.OnData, ws.OnClose))
}
