package nettyws

import (
	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty-transport/websocket"
	"github.com/go-netty/go-netty/codec/frame"
)

var defaultEngine = netty.NewBootstrap(
	netty.WithTransport(websocket.New()),
	netty.WithChannel(netty.NewChannel()),
	netty.WithClientInitializer(clientInitializer),
	netty.WithChildInitializer(childInitializer),
)

func clientInitializer(channel netty.Channel) {
	ws := channel.Attachment().(*Websocket)
	channel.Pipeline().
		AddLast(frame.PacketCodec(1024)).
		AddLast(NewConn(channel, true, ws.OnOpen, ws.OnData, ws.OnClose))
}

func childInitializer(channel netty.Channel) {
	ws := channel.Attachment().(*Websocket)
	channel.Pipeline().
		AddLast(frame.PacketCodec(1024)).
		AddLast(NewConn(channel, false, ws.OnOpen, ws.OnData, ws.OnClose))
}
