// SPDX-License-Identifier: GPL-3.0-or-later

package dnstest

import (
	"crypto/tls"
	"net"
	"net/netip"
	"testing"

	"github.com/bassosimone/nettree/pkg/pkitest"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestTLSWorks(t *testing.T) {
	// create config
	config := NewHandlerConfig()
	config.AddNetipAddr("www.example.com", netip.MustParseAddr("2606:4700::6812:1a78"))
	config.AddNetipAddr("www.example.com", netip.MustParseAddr("2606:4700::6812:1b78"))

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
	srv := MustNewTLSServer(&net.ListenConfig{}, "127.0.0.1:0", cert, handler)
	defer srv.Close()

	// create query
	query := &dns.Msg{}
	query.Question = append(query.Question, dns.Question{
		Name:   dns.CanonicalName("www.example.com"),
		Qtype:  dns.TypeAAAA,
		Qclass: dns.ClassINET,
	})

	// dial
	tlsCfg := &tls.Config{RootCAs: pki.CertPool(), ServerName: "dns.example.com"}
	conn, err := tls.Dial("tcp", srv.Address(), tlsCfg)
	assert.NoError(t, err)
	dconn := &dns.Conn{Conn: conn}

	// exchange
	err = dconn.WriteMsg(query)
	assert.NoError(t, err)
	resp, err := dconn.ReadMsg()
	assert.NoError(t, err)

	// get results
	addrs := collectAddrs(resp)
	expect := []string{"2606:4700::6812:1a78", "2606:4700::6812:1b78"}
	assert.Equal(t, expect, addrs)
}
