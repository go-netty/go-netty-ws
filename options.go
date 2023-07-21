package nettyws

import (
	"crypto/tls"
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

type options struct {
	engine            netty.Bootstrap
	serveMux          *http.ServeMux
	tls               *tls.Config
	checkUTF8         bool
	maxFrameSize      int64
	readBufferSize    int
	writeBufferSize   int
	messageType       MessageType
	compressEnabled   bool
	compressLevel     int
	compressThreshold int64
}

func parseOptions(opt ...Option) *options {
	opts := &options{
		engine:          defaultEngine,
		serveMux:        http.NewServeMux(),
		messageType:     MsgText,
		readBufferSize:  0,
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
		TLS:               wso.tls,
		OpCode:            opCode,
		CheckUTF8:         wso.checkUTF8,
		MaxFrameSize:      wso.maxFrameSize,
		ReadBufferSize:    wso.readBufferSize,
		WriteBufferSize:   wso.writeBufferSize,
		Backlog:           256,
		NoDelay:           false,
		CompressEnabled:   wso.compressEnabled,
		CompressLevel:     wso.compressLevel,
		CompressThreshold: wso.compressThreshold,
		Dialer:            ws.DefaultDialer,
		Upgrader:          ws.DefaultHTTPUpgrader,
		ServeMux:          wso.serveMux,
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
func WithServeTLS(tls *tls.Config) Option {
	return func(options *options) {
		options.tls = tls
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

// WithCompress enable message compression with level, messages below the threshold will not be compressed.
func WithCompress(compressLevel int, compressThreshold int64) Option {
	return func(options *options) {
		options.compressEnabled = true
		options.compressLevel = compressLevel
		options.compressThreshold = compressThreshold
	}
}
