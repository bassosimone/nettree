//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/blob/v3.20.1/internal/mocks/http.go
//

package httptestx

import (
	"net/http"

	"github.com/bassosimone/nettree/pkg/runtimex"
)

// FuncRoundTripper allows to mock any [http.RoundTripper].
type FuncRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

var _ http.RoundTripper = &FuncRoundTripper{}

// RoundTrip implements [http.RoundTripper].
func (rt *FuncRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	runtimex.Assert(rt.RoundTripFunc != nil)
	return rt.RoundTripFunc(req)
}
