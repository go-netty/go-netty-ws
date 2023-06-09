package nettyws

import (
	"bytes"
	"context"
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

// NewConn create a websocket connection.
func NewConn(channel netty.Channel, client bool, onOpen OnOpenFunc, onData OnDataFunc, onClose OnCloseFunc) Conn {
	return &wsConn{channel: channel, client: client, onOpen: onOpen, onData: onData, onClose: onClose}
}

type wsc interface {
	WriteClose(code int, reason string) error
}

type wsConn struct {
	channel  netty.Channel
	client   bool
	onOpen   OnOpenFunc
	onData   OnDataFunc
	onClose  OnCloseFunc
	userdata interface{}
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
	if nil != c.onOpen {
		c.onOpen(c)
		return
	}
	ctx.HandleActive()
}

func (c *wsConn) HandleRead(ctx netty.InboundContext, message netty.Message) {
	if c.onData != nil {
		buffer := message.(*bytes.Buffer)
		c.onData(c, buffer.Bytes())
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

	if nil != c.onClose {
		c.onClose(c, ex)
		return
	}
	ctx.HandleInactive(ex)
}
