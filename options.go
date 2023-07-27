/*
 * Copyright 2023 the go-netty project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	noDelay           bool
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
		noDelay:         true,
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
		NoDelay:           wso.noDelay,
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

// WithNoDelay controls whether the operating system should delay
// packet transmission in hopes of sending fewer packets (Nagle's
// algorithm).  The default is true (no delay), meaning that data is
// sent as soon as possible after a Write.
func WithNoDelay(noDelay bool) Option {
	return func(o *options) {
		o.noDelay = noDelay
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
