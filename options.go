package nettyws

import (
	"net/http"

	"github.com/go-netty/go-netty"
	"github.com/go-netty/go-netty-transport/websocket"
	"github.com/gobwas/ws"
)

type MessageType int

const (
	MsgText MessageType = iota
	MsgBinary
)

type HTTPUpgrader = ws.HTTPUpgrader

type wsOptions struct {
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

func (wso *wsOptions) ToOptions() *websocket.Options {
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

type Option func(*wsOptions)

// WithEngine overwrite default engine
func WithEngine(engine netty.Bootstrap) Option {
	return func(options *wsOptions) {
		options.engine = engine
	}
}

// WithServeMux overwrite default http.ServeMux
func WithServeMux(serveMux *http.ServeMux) Option {
	return func(options *wsOptions) {
		options.serveMux = serveMux
	}
}

// WithServeTLS serve port with TLS
func WithServeTLS(certFile, keyFile string) Option {
	return func(options *wsOptions) {
		options.certFile, options.keyFile = certFile, keyFile
	}
}

// WithBinary switch to binary message mode
func WithBinary() Option {
	return func(options *wsOptions) {
		options.messageType = MsgBinary
	}
}

// WithValidUTF8 enable UTF-8 checks for text frames payload
func WithValidUTF8() Option {
	return func(options *wsOptions) {
		options.checkUTF8 = true
	}
}

// WithMaxFrameSize set the maximum frame size
func WithMaxFrameSize(maxFrameSize int64) Option {
	return func(options *wsOptions) {
		options.maxFrameSize = maxFrameSize
	}
}

// WithBufferSize set the read/write buffer size
func WithBufferSize(readBufferSize, writeBufferSize int) Option {
	return func(options *wsOptions) {
		options.readBufferSize, options.writeBufferSize = readBufferSize, writeBufferSize
	}
}

// WithAsyncWrite enable async write
func WithAsyncWrite(writeQueueSize int, writeForever bool) Option {
	return func(options *wsOptions) {
		if options.engine != defaultEngine {
			panic("please use `netty.NewAsyncWriteChannel(...)` instead of `netty.NewChannel()` in your engine configure")
		}
		options.engine = netty.NewBootstrap(
			netty.WithTransport(websocket.New()),
			netty.WithChannel(netty.NewAsyncWriteChannel(writeQueueSize, writeForever)),
			netty.WithChannelHolder(nil),
			netty.WithClientInitializer(clientInitializer),
			netty.WithChildInitializer(childInitializer),
		)
	}
}

// WithUpgrader set the HTTPUpgrader
func WithUpgrader(upgrader HTTPUpgrader) Option {
	return func(options *wsOptions) {
		options.upgrader = upgrader
	}
}
