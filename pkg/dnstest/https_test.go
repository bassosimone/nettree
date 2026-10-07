// SPDX-License-Identifier: GPL-3.0-or-later

package dnstest

import (
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/netip"
	"testing"

	"github.com/bassosimone/nettree/pkg/pkitest"
	"github.com/bassosimone/nettree/pkg/runtimex"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestHTTPSInvalidRequest(t *testing.T) {
	// create the config and the handler
	config := NewHandlerConfig()
	handler := NewHandler(config)

	// create pki
	pki := pkitest.MustNewPKI("testdata")
	cert := pki.MustNewCert(&pkitest.SelfSignedCertConfig{
		CommonName:   "dns.example.com",
		DNSNames:     []string{"dns.example.com"},
		Organization: []string{"Example"},
	})

	// create server
	srv := MustNewHTTPSServer(&net.ListenConfig{}, "127.0.0.1:0", cert, handler)
	defer srv.Close()

	// create HTTP request with GET method
	httpReq := runtimex.PanicOnError1(http.NewRequest("GET", srv.URL(), nil))
	httpReq.Header.Set("content-type", "application/dns-message")

	// setup HTTPS client
	tlsCfg := &tls.Config{RootCAs: pki.CertPool(), ServerName: "dns.example.com"}
	tdialer := &tls.Dialer{NetDialer: &net.Dialer{}, Config: tlsCfg}
	client := &http.Client{Transport: &http.Transport{DialTLSContext: tdialer.DialContext}}

	// get response body
	httpResp, err := client.Do(httpReq)
	assert.NoError(t, err)
	defer httpResp.Body.Close()
	assert.True(t, httpResp.StatusCode == http.StatusBadRequest)
}

func TestHTTPSMissingContentType(t *testing.T) {
	// create the config and the handler
	config := NewHandlerConfig()
	handler := NewHandler(config)

	// create pki
	pki := pkitest.MustNewPKI("testdata")
	cert := pki.MustNewCert(&pkitest.SelfSignedCertConfig{
		CommonName:   "dns.example.com",
		DNSNames:     []string{"dns.example.com"},
		Organization: []string{"Example"},
	})

	// create server
	srv := MustNewHTTPSServer(&net.ListenConfig{}, "127.0.0.1:0", cert, handler)
	defer srv.Close()

	// create HTTP request containing query
	query := &dns.Msg{}
	query.Question = append(query.Question, dns.Question{
		Name:   dns.CanonicalName("www.example.com"),
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	})
	rawQuery := runtimex.PanicOnError1(query.Pack())
	httpReq := runtimex.PanicOnError1(http.NewRequest("POST", srv.URL(), bytes.NewReader(rawQuery)))
	// No content-type header

	// setup HTTPS client
	tlsCfg := &tls.Config{RootCAs: pki.CertPool(), ServerName: "dns.example.com"}
	tdialer := &tls.Dialer{NetDialer: &net.Dialer{}, Config: tlsCfg}
	client := &http.Client{Transport: &http.Transport{DialTLSContext: tdialer.DialContext}}

	// get response body
	httpResp, err := client.Do(httpReq)
	assert.NoError(t, err)
	defer httpResp.Body.Close()
	assert.True(t, httpResp.StatusCode == http.StatusBadRequest)
}

func TestHTTPSWorks(t *testing.T) {
	// create config
	config := NewHandlerConfig()
	config.AddNetipAddr("www.example.com", netip.MustParseAddr("104.20.34.220"))
	config.AddNetipAddr("www.example.com", netip.MustParseAddr("172.66.144.113"))

	// create handler
	handler := NewHandler(config)

	// create pki
	pki := pkitest.MustNewPKI("testdata")
	cert := pki.MustNewCert(&pkitest.SelfSignedCertConfig{
		CommonName:   "dns.example.com",
		DNSNames:     []string{"dns.example.com"},
		Organization: []string{"Example"},
	})

	// create server
	srv := MustNewHTTPSServer(&net.ListenConfig{}, "127.0.0.1:0", cert, handler)
	defer srv.Close()

	// create HTTP request containing query
	query := &dns.Msg{}
	query.Question = append(query.Question, dns.Question{
		Name:   dns.CanonicalName("www.example.com"),
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	})
	rawQuery := runtimex.PanicOnError1(query.Pack())
	httpReq := runtimex.PanicOnError1(http.NewRequest("POST", srv.URL(), bytes.NewReader(rawQuery)))
	httpReq.Header.Set("content-type", "application/dns-message")

	// setup HTTPS client
	tlsCfg := &tls.Config{RootCAs: pki.CertPool(), ServerName: "dns.example.com"}
	tdialer := &tls.Dialer{NetDialer: &net.Dialer{}, Config: tlsCfg}
	client := &http.Client{Transport: &http.Transport{DialTLSContext: tdialer.DialContext}}

	// get response body
	httpResp, err := client.Do(httpReq)
	assert.NoError(t, err)
	defer httpResp.Body.Close()
	assert.True(t, httpResp.StatusCode == http.StatusOK)
	rawResp, err := io.ReadAll(httpResp.Body)
	assert.NoError(t, err)

	// parse response body
	resp := &dns.Msg{}
	err = resp.Unpack(rawResp)
	assert.NoError(t, err)

	// get results
	addrs := collectAddrs(resp)
	expect := []string{"104.20.34.220", "172.66.144.113"}
	assert.Equal(t, expect, addrs)
}
