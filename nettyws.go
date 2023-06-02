package nettyws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty-transport/websocket"
	"github.com/go-netty/go-netty/codec/frame"
	"github.com/go-netty/go-netty/transport"
)

var engine = netty.NewBootstrap(
	netty.WithTransport(websocket.New()),
	netty.WithClientInitializer(func(channel netty.Channel) {
		ws := channel.Attachment().(*Websocket)
		channel.Pipeline().
			AddLast(frame.PacketCodec(256)).
			AddLast(NewConn(channel, true, ws.OnOpen, ws.OnData, ws.OnClose))
	}),
	netty.WithChildInitializer(func(channel netty.Channel) {
		ws := channel.Attachment().(*Websocket)
		channel.Pipeline().
			AddLast(frame.PacketCodec(256)).
			AddLast(NewConn(channel, false, ws.OnOpen, ws.OnData, ws.OnClose))
	}),
)

type OnOpenFunc func(conn Conn)
type OnDataFunc func(conn Conn, data []byte)
type OnCloseFunc func(conn Conn, err error)

type Websocket struct {
	url      string
	options  *wsOptions
	ctx      context.Context
	cancel   context.CancelFunc
	listener netty.Listener

	OnOpen  OnOpenFunc
	OnData  OnDataFunc
	OnClose OnCloseFunc
}

func NewWebsocket(url string, options ...Option) *Websocket {
	opts := &wsOptions{
		engine:      engine,
		serveMux:    http.NewServeMux(),
		messageType: MsgText,
	}
	for _, op := range options {
		op(opts)
	}

	ctx, cancel := context.WithCancel(opts.engine.Context())
	return &Websocket{url: url, options: opts, ctx: ctx, cancel: cancel}
}

func (ws *Websocket) Open() error {
	_, err := engine.Connect(ws.url, transport.WithAttachment(ws), transport.WithContext(ws.ctx), websocket.WithOptions(ws.options.ToOptions()))
	return err
}

func (ws *Websocket) Listen() error {
	if nil != ws.listener {
		return fmt.Errorf("duplicate listen")
	}
	ws.listener = engine.Listen(ws.url, transport.WithAttachment(ws), transport.WithContext(ws.ctx), websocket.WithOptions(ws.options.ToOptions()))
	return ws.listener.Sync()
}

func (ws *Websocket) Close() error {
	ws.cancel()
	if nil != ws.listener {
		return ws.listener.Close()
	}
	return nil
}
