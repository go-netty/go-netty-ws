package nettyws

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty-transport/websocket"
	"github.com/go-netty/go-netty/transport"
)

type OnOpenFunc func(conn Conn)
type OnDataFunc func(conn Conn, data []byte)
type OnCloseFunc func(conn Conn, err error)

type Websocket struct {
	options   *wsOptions
	ctx       context.Context
	cancel    context.CancelFunc
	listeners sync.Map // map<url , netty.Listener>
	upgrader  websocket.HTTPUpgrader
	holder    netty.ChannelHolder

	OnOpen  OnOpenFunc
	OnData  OnDataFunc
	OnClose OnCloseFunc
}

// NewWebsocket create websocket instance with url and options
func NewWebsocket(options ...Option) *Websocket {
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

	ws := &Websocket{options: opts, holder: NewChannelHolder(1024)}
	ws.ctx, ws.cancel = context.WithCancel(opts.engine.Context())
	ws.upgrader = websocket.NewHTTPUpgrader(opts.engine, transport.WithAttachment(ws), transport.WithContext(ws.ctx), websocket.WithOptions(ws.options.ToOptions()))
	ws.upgrader.Upgrader = opts.upgrader
	return ws
}

// Open websocket client
func (ws *Websocket) Open(addr string) error {
	_, err := ws.options.engine.Connect(addr, transport.WithAttachment(ws), transport.WithContext(ws.ctx), websocket.WithOptions(ws.options.ToOptions()))
	return err
}

// Listen serve addr on this server
func (ws *Websocket) Listen(addr string) error {
	// create listener
	listener := ws.options.engine.Listen(addr, transport.WithAttachment(ws), transport.WithContext(ws.ctx), websocket.WithOptions(ws.options.ToOptions()))
	ws.listeners.Store(addr, listener)

	defer func() {
		if _, loaded := ws.listeners.LoadAndDelete(addr); loaded {
			_ = listener.Close()
		}
	}()

	// listen connections
	return listener.Sync()
}

// Close the listeners and connections
func (ws *Websocket) Close() error {
	// all child or client connections to canceled
	ws.cancel()

	// close all listeners
	ws.listeners.Range(func(key, value interface{}) bool {
		ws.listeners.Delete(key)
		value.(netty.Listener).Close()
		return true
	})

	// close all connections
	ws.holder.CloseAll(ClosedError{Code: 1000, Reason: "websocket shutdown"})
	return nil
}

// UpgradeHTTP Upgrade upgrades http connection to the websocket connection
func (ws *Websocket) UpgradeHTTP(writer http.ResponseWriter, request *http.Request) (conn Conn, err error) {

	select {
	case <-ws.ctx.Done():
		return nil, fmt.Errorf("websocket closed")
	default:
	}

	channel, err := ws.upgrader.Upgrade(writer, request)
	if nil != err {
		return nil, err
	}

	channel.Pipeline().IndexOf(func(handler netty.Handler) bool {
		var ok bool
		conn, ok = handler.(Conn)
		return ok
	})

	if nil == conn {
		err = fmt.Errorf("not found `Conn` Handler in pipleine")
		channel.Close(err)
	}
	return
}
