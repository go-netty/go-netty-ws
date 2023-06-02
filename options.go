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

type wsOptions struct {
	engine          netty.Bootstrap
	serveMux        *http.ServeMux
	certFile        string
	keyFile         string
	checkUTF8       bool
	maxFrameSize    int
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
		MaxFrameSize:    int64(wso.maxFrameSize),
		ReadBufferSize:  wso.readBufferSize,
		WriteBufferSize: wso.writeBufferSize,
		Dialer:          ws.DefaultDialer,
		Upgrader:        ws.DefaultHTTPUpgrader,
		ServeMux:        wso.serveMux,
	}
}

type Option func(*wsOptions)

func WithEngine(engine netty.Bootstrap) Option {
	return func(options *wsOptions) {
		options.engine = engine
	}
}

func WithServeMux(serveMux *http.ServeMux) Option {
	return func(options *wsOptions) {
		options.serveMux = serveMux
	}
}

func WithServeTLS(certFile, keyFile string) Option {
	return func(options *wsOptions) {
		options.certFile, options.keyFile = certFile, keyFile
	}
}

func WithBinary() Option {
	return func(options *wsOptions) {
		options.messageType = MsgBinary
	}
}

func WithValidUTF8() Option {
	return func(options *wsOptions) {
		options.checkUTF8 = true
	}
}

func WithMaxFrameSize(maxFrameSize int) Option {
	return func(options *wsOptions) {
		options.maxFrameSize = maxFrameSize
	}
}

func WithBufferSize(readBufferSize, writeBufferSize int) Option {
	return func(options *wsOptions) {
		options.readBufferSize, options.writeBufferSize = readBufferSize, writeBufferSize
	}
}
