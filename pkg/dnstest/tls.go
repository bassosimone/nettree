// SPDX-License-Identifier: GPL-3.0-or-later

package dnstest

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/bassosimone/nettree/pkg/runtimex"
	"github.com/miekg/dns"
)

// TLSListenConfig is the [*net.ListenConfig] used by [MustNewTLSServer].
type TLSListenConfig interface {
	Listen(ctx context.Context, network, address string) (net.Listener, error)
}

// Ensure that [*net.ListenConfig] implements [TLSListenConfig].
var _ TLSListenConfig = &net.ListenConfig{}

// MustNewTLSServer returns a new [*TLSServer] ready to use.
//
// This method PANICS on failure.
func MustNewTLSServer(lc TLSListenConfig, address string, cert tls.Certificate, handler dns.Handler) *TLSServer {
	listener := runtimex.PanicOnError1(lc.Listen(context.Background(), "tcp", address))
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	tlsListener := tls.NewListener(listener, config)
	srv := &TLSServer{
		address: listener.Addr().String(),
		done:    make(chan struct{}),
		srv: &dns.Server{
			Listener:  tlsListener,
			Handler:   handler,
			TLSConfig: config,
		},
	}
	go func() {
		srv.srv.ActivateAndServe() // in background
		close(srv.done)
	}()
	return srv
}

// TLSServer is a server for testing DNS-over-TLS.
type TLSServer struct {
	// address is the address to use.
	address string

	// done is closed when done.
	done chan struct{}

	// srv is the server.
	srv *dns.Server
}

// Address returns the listening TLS address for this server.
func (srv *TLSServer) Address() string {
	return srv.address
}

// Close closes the socket used by this server.
func (srv *TLSServer) Close() {
	runtimex.PanicOnError0(srv.srv.Shutdown())
	<-srv.done
}
