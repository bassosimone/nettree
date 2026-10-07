//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/tree/v3.20.1/internal/mocks
//

package netstub

import (
	"net"

	"github.com/bassosimone/nettree/pkg/runtimex"
)

// FuncAddr allows to mock any [net.Addr].
type FuncAddr struct {
	NetworkFunc func() string
	StringFunc  func() string
}

var _ net.Addr = &FuncAddr{}

// Network implements [net.Addr].
func (fa *FuncAddr) Network() string {
	runtimex.Assert(fa.NetworkFunc != nil)
	return fa.NetworkFunc()
}

// String implements [net.Addr].
func (fa *FuncAddr) String() string {
	runtimex.Assert(fa.StringFunc != nil)
	return fa.StringFunc()
}
