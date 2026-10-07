//
// SPDX-License-Identifier: BSD-3-Clause
//
// Adapted from: https://github.com/ooni/netem/blob/608dcbcd82b8eabcb675d482e2ca83cf3a41c27d/dnsserver.go
//

package dnstest

import (
	"context"
	"net"

	"github.com/bassosimone/nettree/pkg/runtimex"
	"github.com/miekg/dns"
)

// UDPListenConfig is the [*net.ListenConfig] used by [MustNewUDPServer].
type UDPListenConfig interface {
	ListenPacket(ctx context.Context, network, address string) (net.PacketConn, error)
}

// Ensure that [*net.ListenConfig] implements [UDPListenConfig].
var _ UDPListenConfig = &net.ListenConfig{}

// MustNewUDPServer returns a new [*UDPServer] ready to use.
//
// This method PANICS on failure.
func MustNewUDPServer(lc UDPListenConfig, address string, handler dns.Handler) *UDPServer {
	pconn := runtimex.PanicOnError1(lc.ListenPacket(context.Background(), "udp", address))
	srv := &UDPServer{
		address: pconn.LocalAddr().String(),
		done:    make(chan struct{}),
		srv: &dns.Server{
			PacketConn: pconn,
			Handler:    handler,
		},
	}
	go func() {
		srv.srv.ActivateAndServe() // in background
		close(srv.done)
	}()
	return srv
}

// UDPServer is a server for testing DNS-over-UDP.
type UDPServer struct {
	// address is the address to use.
	address string

	// done is closed when done.
	done chan struct{}

	// srv is the server.
	srv *dns.Server
}

// Address returns the listening UDP address for this server.
func (srv *UDPServer) Address() string {
	return srv.address
}

// Close closes the socket used by this server.
func (srv *UDPServer) Close() {
	runtimex.PanicOnError0(srv.srv.Shutdown())
	<-srv.done
}
