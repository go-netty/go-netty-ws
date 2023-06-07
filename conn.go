package nettyws

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/go-netty/go-netty"
	"github.com/gobwas/ws/wsutil"
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

type Message = netty.Message

type Conn interface {
	Context() context.Context
	LocalAddr() string
	RemoteAddr() string
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	Write(message []byte) error
	WriteClose(code int, reason string) error
	Close() error
	Userdata() interface{}
	SetUserdata(userdata interface{})
}

func NewConn(channel netty.Channel, client bool, onOpen OnOpenFunc, onData OnDataFunc, onClose OnCloseFunc) Conn {
	return &wsConn{channel: channel, client: client, onOpen: onOpen, onData: onData, onClose: onClose}
}

type wsConn struct {
	channel  netty.Channel
	client   bool
	onOpen   OnOpenFunc
	onData   OnDataFunc
	onClose  OnCloseFunc
	userdata interface{}
}

func (c *wsConn) Context() context.Context {
	return c.channel.Context()
}

func (c *wsConn) LocalAddr() string {
	return c.channel.LocalAddr()
}

func (c *wsConn) RemoteAddr() string {
	return c.channel.RemoteAddr()
}

func (c *wsConn) SetDeadline(t time.Time) error {
	return c.channel.Transport().SetDeadline(t)
}

func (c *wsConn) SetReadDeadline(t time.Time) error {
	return c.channel.Transport().SetReadDeadline(t)
}

func (c *wsConn) SetWriteDeadline(t time.Time) error {
	return c.channel.Transport().SetWriteDeadline(t)
}

func (c *wsConn) Write(message []byte) error {
	_, err := c.channel.Write1(message)
	return err
}

func (c *wsConn) WriteClose(code int, reason string) error {
	type wst interface {
		WriteClose(code int, reason string) error
	}
	return c.channel.Transport().(wst).WriteClose(code, reason)
}

func (c *wsConn) Close() error {
	c.channel.Close(nil)
	return nil
}

func (c *wsConn) Userdata() interface{} {
	return c.userdata
}

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
