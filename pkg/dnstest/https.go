// SPDX-License-Identifier: GPL-3.0-or-later

package dnstest

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/bassosimone/nettree/pkg/runtimex"
	"github.com/miekg/dns"
)

// HTTPSListenConfig is the [*net.ListenConfig] used by [MustNewHTTPSServer].
type HTTPSListenConfig interface {
	Listen(ctx context.Context, network, address string) (net.Listener, error)
}

// Ensure that [*net.ListenConfig] implements [HTTPSListenConfig].
var _ HTTPSListenConfig = &net.ListenConfig{}

// newUnstartedServer creates an httptest.Server with timeouts so that in-flight handlers
// cannot block Close() forever when the underlying gVisor stack tears down.
func newUnstartedServer(handler http.Handler) *httptest.Server {
	httpSrv := httptest.NewUnstartedServer(handler)
	httpSrv.Config.ReadHeaderTimeout = 5 * time.Second
	httpSrv.Config.ReadTimeout = 30 * time.Second
	httpSrv.Config.WriteTimeout = 30 * time.Second
	httpSrv.Config.IdleTimeout = 30 * time.Second
	return httpSrv
}

// MustNewHTTPSServer returns a new [*HTTPSServer] ready to use.
//
// This method PANICS on failure.
func MustNewHTTPSServer(
	lc HTTPSListenConfig, address string, cert tls.Certificate, handler *Handler) *HTTPSServer {
	listener := runtimex.PanicOnError1(lc.Listen(context.Background(), "tcp", address))
	hs := newUnstartedServer(HTTPSHandler{handler})
	hs.Listener = listener
	hs.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	hs.EnableHTTP2 = true
	hs.StartTLS()
	srv := &HTTPSServer{
		address: listener.Addr().String(),
		srv:     hs,
	}
	return srv
}

// HTTPSServer is a server for testing DNS-over-HTTPS.
type HTTPSServer struct {
	// address is the address to use.
	address string

	// srv is the HTTPS server.
	srv *httptest.Server
}

// URL returns the URL for this server.
func (srv *HTTPSServer) URL() string {
	return srv.srv.URL
}

// Close closes the socket used by this server.
func (srv *HTTPSServer) Close() {
	srv.srv.Close()
}

// HTTPSHandler handles DoH requests.
type HTTPSHandler struct {
	Handler *Handler
}

// Ensure that [HTTPSHandler] implements [http.Handler].
var _ http.Handler = HTTPSHandler{}

// ServeHTTP implements [http.Handler].
func (hh HTTPSHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
	}()
	runtimex.Assert(req.Method == "POST")
	runtimex.Assert(req.Header.Get("content-type") == "application/dns-message")
	rawQuery := runtimex.PanicOnError1(io.ReadAll(req.Body))
	query := &dns.Msg{}
	runtimex.PanicOnError0(query.Unpack(rawQuery))
	resp := hh.Handler.PrepareResponse(query)
	rawResp := runtimex.PanicOnError1(resp.Pack())
	w.Header().Set("content-type", "application/dns-message")
	w.Write(rawResp)
}
