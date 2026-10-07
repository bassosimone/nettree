//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/tree/v3.20.1/internal/mocks
//

package netstub

import (
	"context"
	"net"

	"github.com/bassosimone/nettree/pkg/runtimex"
)

// FuncDialer allows to mock any dialer with a DialContext method.
type FuncDialer struct {
	DialContextFunc func(ctx context.Context, network, address string) (net.Conn, error)
}

var _ interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
} = &FuncDialer{}

// DialContext implements a dialer with a DialContext method.
func (fd *FuncDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	runtimex.Assert(fd.DialContextFunc != nil)
	return fd.DialContextFunc(ctx, network, address)
}
