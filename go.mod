module github.com/go-netty/go-netty-ws

go 1.24.0

require (
	github.com/go-netty/go-netty v1.6.7
	github.com/go-netty/go-netty-transport v1.7.14
	github.com/gobwas/ws v1.4.0
)

require (
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	golang.org/x/sys v0.38.0 // indirect
)

//replace github.com/go-netty/go-netty => ../go-netty
//replace github.com/go-netty/go-netty-transport => ../go-netty-transport
