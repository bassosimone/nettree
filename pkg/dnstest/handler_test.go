// SPDX-License-Identifier: GPL-3.0-or-later

package dnstest

import (
	"net/netip"
	"slices"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestHandlerConfigRemove(t *testing.T) {
	// add www.example.com
	config := NewHandlerConfig()
	config.AddNetipAddr("www.example.com", netip.MustParseAddr("1.1.1.1"))

	// ensure that we get the expected response
	rrs, found := config.Lookup("www.example.com", dns.TypeA)
	assert.True(t, found)
	assert.True(t, len(rrs) == 1)

	// remove www.example.com
	config.Remove("www.example.com")

	// ensure that we get the expected response
	rrs, found = config.Lookup("www.example.com", dns.TypeA)
	assert.True(t, !found)
	assert.True(t, len(rrs) == 0)
}

func TestHandlerConfigClone(t *testing.T) {
	// add www.example.com
	config1 := NewHandlerConfig()
	config1.AddNetipAddr("www.example.com", netip.MustParseAddr("1.1.1.1"))

	// clone and remove www.example.com
	config2 := config1.Clone()
	config2.Remove("www.example.com")

	// ensure that we get the expected response (config1)
	rrs, found := config1.Lookup("www.example.com", dns.TypeA)
	assert.True(t, found)
	assert.True(t, len(rrs) == 1)

	// ensure that we get the expected response (config2)
	rrs, found = config2.Lookup("www.example.com", dns.TypeA)
	assert.True(t, !found)
	assert.True(t, len(rrs) == 0)
}

// collectAddrs extracts all A and AAAA records from a DNS message's Answer section
// and returns them as a sorted slice of strings for stable comparison.
func collectAddrs(resp *dns.Msg) (output []string) {
	for _, rec := range resp.Answer {
		switch rec := rec.(type) {
		case *dns.A:
			output = append(output, rec.A.String())
		case *dns.AAAA:
			output = append(output, rec.AAAA.String())
		}
	}
	slices.Sort(output)
	return
}

// collectCNAMEs extracts all CNAME records from a DNS message's Answer section
// and returns them as a sorted slice of strings for stable comparison.
func collectCNAMEs(answer []dns.RR) (output []string) {
	for _, rec := range answer {
		if cname, ok := rec.(*dns.CNAME); ok {
			output = append(output, cname.Target)
		}
	}
	slices.Sort(output)
	return
}

func TestHandlerPrepareResponse(t *testing.T) {
	type testCase struct {
		name           string
		getConfig      func() *HandlerConfig
		getQuery       func() *dns.Msg
		expectedRcode  int
		validateAnswer func(t *testing.T, resp *dns.Msg)
	}

	testCases := []testCase{
		{
			name: "successful A record lookup",
			getConfig: func() *HandlerConfig {
				config := NewHandlerConfig()
				config.AddNetipAddr("www.example.com", netip.MustParseAddr("1.1.1.1"))
				return config
			},
			getQuery: func() *dns.Msg {
				query := new(dns.Msg)
				query.Question = append(query.Question, dns.Question{
					Name:   dns.CanonicalName("www.example.com"),
					Qtype:  dns.TypeA,
					Qclass: dns.ClassINET,
				})
				return query
			},
			expectedRcode: dns.RcodeSuccess,
			validateAnswer: func(t *testing.T, resp *dns.Msg) {
				addrs := collectAddrs(resp)
				assert.Equal(t, []string{"1.1.1.1"}, addrs)
			},
		},

		{
			name: "type not found",
			getConfig: func() *HandlerConfig {
				config := NewHandlerConfig()
				config.AddNetipAddr("www.example.com", netip.MustParseAddr("1.1.1.1"))
				return config
			},
			getQuery: func() *dns.Msg {
				query := new(dns.Msg)
				query.Question = append(query.Question, dns.Question{
					Name:   dns.CanonicalName("www.example.com"),
					Qtype:  dns.TypeAAAA,
					Qclass: dns.ClassINET,
				})
				return query
			},
			expectedRcode: dns.RcodeSuccess,
			validateAnswer: func(t *testing.T, resp *dns.Msg) {
				assert.Empty(t, resp.Answer)
			},
		},

		{
			name: "name not found",
			getConfig: func() *HandlerConfig {
				return NewHandlerConfig()
			},
			getQuery: func() *dns.Msg {
				query := new(dns.Msg)
				query.Question = append(query.Question, dns.Question{
					Name:   dns.CanonicalName("nonexistent.example.com"),
					Qtype:  dns.TypeA,
					Qclass: dns.ClassINET,
				})
				return query
			},
			expectedRcode: dns.RcodeNameError,
			validateAnswer: func(t *testing.T, resp *dns.Msg) {
				assert.Empty(t, resp.Answer)
			},
		},

		{
			name: "cname chase",
			getConfig: func() *HandlerConfig {
				config := NewHandlerConfig()
				config.AddCNAME("alias.example.com", "real.example.com")
				config.AddNetipAddr("real.example.com", netip.MustParseAddr("8.8.8.8"))
				return config
			},
			getQuery: func() *dns.Msg {
				query := new(dns.Msg)
				query.Question = append(query.Question, dns.Question{
					Name:   dns.CanonicalName("alias.example.com"),
					Qtype:  dns.TypeA,
					Qclass: dns.ClassINET,
				})
				return query
			},
			expectedRcode: dns.RcodeSuccess,
			validateAnswer: func(t *testing.T, resp *dns.Msg) {
				addrs := collectAddrs(resp)
				assert.Equal(t, []string{"8.8.8.8"}, addrs)
				cnames := collectCNAMEs(resp.Answer)
				assert.Equal(t, []string{dns.CanonicalName("real.example.com")}, cnames)
			},
		},

		{
			name: "invalid query (no question)",
			getConfig: func() *HandlerConfig {
				return NewHandlerConfig()
			},
			getQuery: func() *dns.Msg {
				return &dns.Msg{}
			},
			expectedRcode: dns.RcodeRefused,
			validateAnswer: func(t *testing.T, resp *dns.Msg) {
				assert.Empty(t, resp.Answer)
			},
		},

		{
			name: "cname loop",
			getConfig: func() *HandlerConfig {
				config := NewHandlerConfig()
				config.AddCNAME("a.example.com", "b.example.com")
				config.AddCNAME("b.example.com", "a.example.com")
				return config
			},
			getQuery: func() *dns.Msg {
				query := new(dns.Msg)
				query.Question = append(query.Question, dns.Question{
					Name:   dns.CanonicalName("a.example.com"),
					Qtype:  dns.TypeA,
					Qclass: dns.ClassINET,
				})
				return query
			},
			expectedRcode: dns.RcodeServerFailure,
			validateAnswer: func(t *testing.T, resp *dns.Msg) {
				assert.Empty(t, resp.Answer)
			},
		},

		{
			name: "invalid class (CHAOS)",
			getConfig: func() *HandlerConfig {
				return NewHandlerConfig()
			},
			getQuery: func() *dns.Msg {
				query := &dns.Msg{}
				query.Question = append(query.Question, dns.Question{
					Name:   dns.CanonicalName("www.example.com"),
					Qtype:  dns.TypeA,
					Qclass: dns.ClassCHAOS,
				})
				return query
			},
			expectedRcode: dns.RcodeRefused,
			validateAnswer: func(t *testing.T, resp *dns.Msg) {
				assert.Empty(t, resp.Answer)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewHandler(tc.getConfig())
			response := handler.PrepareResponse(tc.getQuery())

			assert.NotNil(t, response)
			assert.Equal(t, tc.expectedRcode, response.Rcode)
			if tc.validateAnswer != nil {
				tc.validateAnswer(t, response)
			}
		})
	}
}
