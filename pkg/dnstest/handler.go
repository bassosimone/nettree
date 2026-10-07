//
// SPDX-License-Identifier: BSD-3-Clause
//
// Adapted from: https://github.com/ooni/netem/blob/608dcbcd82b8eabcb675d482e2ca83cf3a41c27d/dnsserver.go
//

package dnstest

import (
	"net/netip"
	"sync"

	"github.com/miekg/dns"
)

// HandlerConfig maps domain names to resource records.
//
// Construct using [NewHandlerConfig].
type HandlerConfig struct {
	mu  sync.Mutex
	rrs map[string][]dns.RR
}

// NewHandlerConfig constructs a [*HandlerConfig] instance.
func NewHandlerConfig() *HandlerConfig {
	return &HandlerConfig{
		mu:  sync.Mutex{},
		rrs: map[string][]dns.RR{},
	}
}

// Clone clones a [*HandlerConfig] instance.
func (c *HandlerConfig) Clone() *HandlerConfig {
	c.mu.Lock()
	out := NewHandlerConfig()
	for key, value := range c.rrs {
		v := make([]dns.RR, 0, len(value))
		v = append(v, value...)
		out.rrs[key] = v
	}
	c.mu.Unlock()
	return out
}

// handlerDefaultTTL is the defaultTTL used by [*HandlerConfig].
const handlerDefaultTTL = 3600

// AddNetipAddr adds a given [netip.Addr] to the [*HandlerConfig].
func (c *HandlerConfig) AddNetipAddr(name string, addr netip.Addr) {
	name = dns.CanonicalName(name)

	var record dns.RR
	switch addr.Is6() {
	case true:
		record = &dns.AAAA{
			Hdr: dns.RR_Header{
				Name:   name,
				Rrtype: dns.TypeAAAA,
				Class:  dns.ClassINET,
				Ttl:    handlerDefaultTTL,
			},
			AAAA: addr.AsSlice(),
		}

	default:
		record = &dns.A{
			Hdr: dns.RR_Header{
				Name:   name,
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    handlerDefaultTTL,
			},
			A: addr.AsSlice(),
		}
	}

	c.mu.Lock()
	c.rrs[name] = append(c.rrs[name], record)
	c.mu.Unlock()
}

// AddCNAME adds a CNAME alias record for the given name.
func (c *HandlerConfig) AddCNAME(name, cname string) {
	name, cname = dns.CanonicalName(name), dns.CanonicalName(cname)

	record := &dns.CNAME{
		Hdr: dns.RR_Header{
			Name:   name,
			Rrtype: dns.TypeCNAME,
			Class:  dns.ClassINET,
			Ttl:    handlerDefaultTTL,
		},
		Target: cname,
	}

	c.mu.Lock()
	c.rrs[name] = append(c.rrs[name], record)
	c.mu.Unlock()
}

// Remove removes records from the [*HandlerConfig].
func (c *HandlerConfig) Remove(name string) {
	c.mu.Lock()
	delete(c.rrs, dns.CanonicalName(name))
	c.mu.Unlock()
}

// Lookup searches for a name inside the [*HandlerConfig].
//
// A false return value indicates that the record does not exist
// while a true return value without records indicates that we don't
// have records for the given type.
func (c *HandlerConfig) Lookup(name string, qtype uint16) ([]dns.RR, bool) {
	var filtered []dns.RR
	c.mu.Lock()

	records, found := c.rrs[dns.CanonicalName(name)]
	for _, rr := range records {
		if qtype == rr.Header().Rrtype {
			filtered = append(filtered, rr)
		}
	}

	c.mu.Unlock()
	return filtered, found
}

// Handler is a [dns.Handler] using [*HandlerConfig] to serve responses.
//
// Construct using [NewHandler].
type Handler struct {
	cfg *HandlerConfig
}

// NewHandler returns a new [*Handler] instance.
func NewHandler(config *HandlerConfig) *Handler {
	return &Handler{config}
}

// Ensure that [handler] implements [dns.Handler].
var _ dns.Handler = &Handler{}

// ServeDNS implements [dns.Handler].
func (h *Handler) ServeDNS(rw dns.ResponseWriter, query *dns.Msg) {
	rw.WriteMsg(h.PrepareResponse(query))
}

// PrepareResponse returns a [*dns.Msg] response for the given [*dns.Msg] query.
func (h *Handler) PrepareResponse(query *dns.Msg) *dns.Msg {
	// 1. reject blatantly wrong queries
	if query.Response || len(query.Question) != 1 {
		resp := &dns.Msg{}
		resp.SetRcode(query, dns.RcodeRefused)
		return resp
	}

	// 2. find the corresponding record
	q0 := query.Question[0]
	if q0.Qclass != dns.ClassINET {
		resp := &dns.Msg{}
		resp.SetRcode(query, dns.RcodeRefused)
		return resp
	}

	// 3. lookup with the config, following CNAME chains
	var cnames []dns.RR
	qName, qType := q0.Name, q0.Qtype
	const maxCNAMEChain = 10
	for range maxCNAMEChain {
		// 3.1. execute the query requested by the user
		records, found := h.cfg.Lookup(qName, qType)

		switch {
		// 3.2. the query returned records
		case found && len(records) > 0:
			resp := &dns.Msg{}
			resp.SetReply(query)
			resp.Answer = append(cnames, records...)
			return resp

		// 3.3. no records but the name exists
		case found && len(records) <= 0:
			// 3.3.1. see whether a CNAME lookup could actually help
			records, found := h.cfg.Lookup(qName, dns.TypeCNAME)

			switch {
			// 3.3.2. we have at least a CNAME entry
			case found && len(records) >= 1:
				cnames = append(cnames, records...)
				// Type assertion is safe: we specifically queried for TypeCNAME,
				// so Config.Lookup only returns CNAME records.
				qName = records[0].(*dns.CNAME).Target

			// 3.3.3. otherwise, NOERROR (name exists but type not found)
			default:
				resp := &dns.Msg{}
				resp.SetReply(query) // NOERROR, empty answer
				return resp
			}

		// 3.4. otherwise, NXDOMAIN
		default:
			resp := &dns.Msg{}
			resp.SetRcode(query, dns.RcodeNameError)
			return resp
		}
	}

	// 3.5. CNAME chain too long: avoid possible loop
	resp := &dns.Msg{}
	resp.SetRcode(query, dns.RcodeServerFailure)
	return resp
}
