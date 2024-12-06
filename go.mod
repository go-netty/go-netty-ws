module github.com/go-netty/go-netty-ws

go 1.22

require (
	github.com/go-netty/go-netty v1.6.6
	github.com/go-netty/go-netty-transport v1.7.12
	github.com/gobwas/ws v1.4.0
)

require (
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	golang.org/x/sys v0.26.0 // indirect
)

//replace github.com/go-netty/go-netty => ../go-netty
//replace github.com/go-netty/go-netty-transport => ../go-netty-transport
