// SPDX-License-Identifier: GPL-3.0-or-later

package dnstest

import (
	"context"
	"net"

	"github.com/bassosimone/nettree/pkg/runtimex"
	"github.com/miekg/dns"
)

// TCPListenConfig is the [*net.ListenConfig] used by [MustNewTCPServer].
type TCPListenConfig interface {
	Listen(ctx context.Context, network, address string) (net.Listener, error)
}

// Ensure that [*net.ListenConfig] implements [TCPListenConfig].
var _ TCPListenConfig = &net.ListenConfig{}

// MustNewTCPServer returns a new [*TCPServer] ready to use.
//
// This method PANICS on failure.
func MustNewTCPServer(lc TCPListenConfig, address string, handler dns.Handler) *TCPServer {
	listener := runtimex.PanicOnError1(lc.Listen(context.Background(), "tcp", address))
	srv := &TCPServer{
		address: listener.Addr().String(),
		done:    make(chan struct{}),
		srv: &dns.Server{
			Listener: listener,
			Handler:  handler,
		},
	}
	go func() {
		srv.srv.ActivateAndServe() // in background
		close(srv.done)
	}()
	return srv
}

// TCPServer is a server for testing DNS-over-TCP.
type TCPServer struct {
	// address is the address to use.
	address string

	// done is closed when done.
	done chan struct{}

	// srv is the server.
	srv *dns.Server
}

// Address returns the listening TCP address for this server.
func (srv *TCPServer) Address() string {
	return srv.address
}

// Close closes the socket used by this server.
func (srv *TCPServer) Close() {
	runtimex.PanicOnError0(srv.srv.Shutdown())
	<-srv.done
}
