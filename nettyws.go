package nettyws

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty-transport/websocket"
	"github.com/go-netty/go-netty/transport"
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

// NewWebsocket create websocket instance with url and options
func NewWebsocket(url string, options ...Option) *Websocket {
	opts := &wsOptions{
		engine:          defaultEngine,
		serveMux:        http.NewServeMux(),
		messageType:     MsgText,
		readBufferSize:  1024,
		writeBufferSize: 0,
	}
	for _, op := range options {
		op(opts)
	}

	ctx, cancel := context.WithCancel(opts.engine.Context())
	return &Websocket{url: url, options: opts, ctx: ctx, cancel: cancel}
}

// Open websocket client
func (ws *Websocket) Open() error {
	_, err := ws.options.engine.Connect(ws.url, transport.WithAttachment(ws), transport.WithContext(ws.ctx), websocket.WithOptions(ws.options.ToOptions()))
	return err
}

// Listen serve port on this server
func (ws *Websocket) Listen() error {
	if nil != ws.listener {
		return fmt.Errorf("duplicate listen")
	}
	ws.listener = ws.options.engine.Listen(ws.url, transport.WithAttachment(ws), transport.WithContext(ws.ctx), websocket.WithOptions(ws.options.ToOptions()))
	return ws.listener.Sync()
}

// Close the client or server child connections
func (ws *Websocket) Close() error {
	ws.cancel()
	if nil != ws.listener {
		return ws.listener.Close()
	}
	return nil
}
