package nettyws

import (
	"bytes"
	"context"
	"net/http"
	"time"

	"github.com/go-netty/go-netty"
	"github.com/gobwas/ws/wsutil"
)

// Conn is a websocket connection.
type Conn interface {
	// Context returns the context of the connection.
	Context() context.Context
	// LocalAddr returns the local network address.
	LocalAddr() string
	// RemoteAddr returns the remote network address.
	RemoteAddr() string
	// Header returns the HTTP header on handshake request.
	Header() http.Header
	// SetDeadline sets the read and write deadlines associated
	// with the connection. It is equivalent to calling both
	// SetReadDeadline and SetWriteDeadline.
	SetDeadline(t time.Time) error
	// SetReadDeadline sets the deadline for future Read calls
	// and any currently-blocked Read call.
	// A zero value for t means Read will not time out.
	SetReadDeadline(t time.Time) error
	// SetWriteDeadline sets the deadline for future Write calls
	// and any currently-blocked Write call.
	// Even if write times out, it may return n > 0, indicating that
	// some of the data was successfully written.
	// A zero value for t means Write will not time out.
	SetWriteDeadline(t time.Time) error
	// Write writes a message to the connection.
	Write(message []byte) error
	// WriteClose write websocket close frame with code and close reason.
	WriteClose(code int, reason string) error
	// Close closes the connection.
	Close() error
	// Userdata returns the user-data.
	Userdata() interface{}
	// SetUserdata sets the user-data.
	SetUserdata(userdata interface{})
}

type wsc interface {
	WriteClose(code int, reason string) error
}

type wsh interface {
	Route() string
	Header() http.Header
}

type wsConn struct {
	ws       *Websocket
	channel  netty.Channel
	client   bool
	userdata interface{}
}

// newConn create a websocket connection.
func newConn(ws *Websocket, channel netty.Channel, client bool) Conn {
	return &wsConn{ws: ws, channel: channel, client: client}
}

// Context returns the context of the connection.
func (c *wsConn) Context() context.Context {
	return c.channel.Context()
}

// LocalAddr returns the local network address.
func (c *wsConn) LocalAddr() string {
	return c.channel.LocalAddr()
}

// RemoteAddr returns the remote network address.
func (c *wsConn) RemoteAddr() string {
	return c.channel.RemoteAddr()
}

// Header returns the HTTP header on handshake request.
func (c *wsConn) Header() http.Header {
	return c.channel.Transport().(wsh).Header()
}

// SetDeadline sets the read and write deadlines associated
// with the connection. It is equivalent to calling both
// SetReadDeadline and SetWriteDeadline.
func (c *wsConn) SetDeadline(t time.Time) error {
	return c.channel.Transport().SetDeadline(t)
}

// SetReadDeadline sets the deadline for future Read calls
// and any currently-blocked Read call.
// A zero value for t means Read will not time out.
func (c *wsConn) SetReadDeadline(t time.Time) error {
	return c.channel.Transport().SetReadDeadline(t)
}

// SetWriteDeadline sets the deadline for future Write calls
// and any currently-blocked Write call.
// Even if write times out, it may return n > 0, indicating that
// some of the data was successfully written.
// A zero value for t means Write will not time out.
func (c *wsConn) SetWriteDeadline(t time.Time) error {
	return c.channel.Transport().SetWriteDeadline(t)
}

// Write writes a message to the connection.
func (c *wsConn) Write(message []byte) error {
	_, err := c.channel.Write1(message)
	return err
}

// WriteClose write websocket close frame with code and close reason.
func (c *wsConn) WriteClose(code int, reason string) error {
	return c.channel.Transport().(wsc).WriteClose(code, reason)
}

// Close closes the connection.
func (c *wsConn) Close() error {
	c.channel.Close(nil)
	return nil
}

// Userdata returns the user-data.
func (c *wsConn) Userdata() interface{} {
	return c.userdata
}

// SetUserdata sets the user-data.
func (c *wsConn) SetUserdata(userdata interface{}) {
	c.userdata = userdata
}

func (c *wsConn) HandleActive(ctx netty.ActiveContext) {
	if onOpen := c.ws.OnOpen; nil != onOpen {
		onOpen(c)
		return
	}
	ctx.HandleActive()
}

func (c *wsConn) HandleRead(ctx netty.InboundContext, message netty.Message) {
	if onData := c.ws.OnData; onData != nil {
		onData(c, message.(*bytes.Buffer).Bytes())
		return
	}
	ctx.HandleRead(message)
}

func (c *wsConn) HandleException(ctx netty.ExceptionContext, ex netty.Exception) {
	ctx.Close(ex)
}

func (c *wsConn) HandleInactive(ctx netty.InactiveContext, ex netty.Exception) {
	// covert error
	if closeErr, ok := ex.(wsutil.ClosedError); ok {
		ex = ClosedError{Code: int(closeErr.Code), Reason: closeErr.Reason}
	}

	if onClose := c.ws.OnClose; nil != onClose {
		onClose(c, ex)
		return
	}
	ctx.HandleInactive(ex)
}
