module github.com/go-netty/go-netty-ws

go 1.21

toolchain go1.22.0

require (
	github.com/go-netty/go-netty v1.6.4
	github.com/go-netty/go-netty-transport v1.7.9
	github.com/gobwas/ws v1.3.2
)

require (
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	golang.org/x/sys v0.15.0 // indirect
)

//replace github.com/go-netty/go-netty => ../go-netty
//replace github.com/go-netty/go-netty-transport => ../go-netty-transport
