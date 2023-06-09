package nettyws

import (
	"net/http"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty-transport/websocket"
	"github.com/gobwas/ws"
)

// MessageType websocket message type
type MessageType int

const (
	// MsgText text message
	MsgText MessageType = iota
	// MsgBinary binary message
	MsgBinary
)

// HTTPUpgrader contains options for upgrading connection to websocket from net/http Handler arguments.
type HTTPUpgrader = ws.HTTPUpgrader

type options struct {
	engine          netty.Bootstrap
	serveMux        *http.ServeMux
	upgrader        HTTPUpgrader
	certFile        string
	keyFile         string
	checkUTF8       bool
	maxFrameSize    int64
	readBufferSize  int
	writeBufferSize int
	messageType     MessageType
}

func parseOptions(opt ...Option) *options {
	opts := &options{
		engine:          defaultEngine,
		serveMux:        http.NewServeMux(),
		messageType:     MsgText,
		readBufferSize:  1024,
		writeBufferSize: 0,
	}
	for _, op := range opt {
		op(opts)
	}
	return opts
}

func (wso *options) wsOptions() *websocket.Options {
	opCode := ws.OpText
	if MsgBinary == wso.messageType {
		opCode = ws.OpBinary
	}
	return &websocket.Options{
		Cert:            wso.certFile,
		Key:             wso.keyFile,
		OpCode:          opCode,
		CheckUTF8:       wso.checkUTF8,
		MaxFrameSize:    wso.maxFrameSize,
		ReadBufferSize:  wso.readBufferSize,
		WriteBufferSize: wso.writeBufferSize,
		Backlog:         256,
		Dialer:          ws.DefaultDialer,
		Upgrader:        ws.DefaultHTTPUpgrader,
		ServeMux:        wso.serveMux,
	}
}

type Option func(*options)

// WithServeMux overwrite default http.ServeMux
func WithServeMux(serveMux *http.ServeMux) Option {
	return func(options *options) {
		options.serveMux = serveMux
	}
}

// WithServeTLS serve port with TLS
func WithServeTLS(certFile, keyFile string) Option {
	return func(options *options) {
		options.certFile, options.keyFile = certFile, keyFile
	}
}

// WithBinary switch to binary message mode
func WithBinary() Option {
	return func(options *options) {
		options.messageType = MsgBinary
	}
}

// WithValidUTF8 enable UTF-8 checks for text frames payload
func WithValidUTF8() Option {
	return func(options *options) {
		options.checkUTF8 = true
	}
}

// WithMaxFrameSize set the maximum frame size
func WithMaxFrameSize(maxFrameSize int64) Option {
	return func(options *options) {
		options.maxFrameSize = maxFrameSize
	}
}

// WithBufferSize set the read/write buffer size
func WithBufferSize(readBufferSize, writeBufferSize int) Option {
	return func(options *options) {
		options.readBufferSize, options.writeBufferSize = readBufferSize, writeBufferSize
	}
}

// WithAsyncWrite enable async write
func WithAsyncWrite(writeQueueSize int, writeForever bool) Option {
	return func(options *options) {
		options.engine = netty.NewBootstrap(
			netty.WithTransport(websocket.New()),
			netty.WithChannel(netty.NewAsyncWriteChannel(writeQueueSize, writeForever)),
			netty.WithChannelHolder(nil),
			netty.WithClientInitializer(makeInitializer(true)),
			netty.WithChildInitializer(makeInitializer(false)),
		)
	}
}

// WithUpgrader set the HTTPUpgrader
func WithUpgrader(upgrader HTTPUpgrader) Option {
	return func(options *options) {
		options.upgrader = upgrader
	}
}
