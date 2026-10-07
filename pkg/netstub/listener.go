// SPDX-License-Identifier: GPL-3.0-or-later

package netstub

import (
	"net"

	"github.com/bassosimone/nettree/pkg/runtimex"
)

// FuncListener allows to mock any [net.Listener].
type FuncListener struct {
	AcceptFunc func() (net.Conn, error)
	CloseFunc  func() error
	AddrFunc   func() net.Addr
}

var _ net.Listener = &FuncListener{}

// Accept implements [net.Listener].
func (fl *FuncListener) Accept() (net.Conn, error) {
	runtimex.Assert(fl.AcceptFunc != nil)
	return fl.AcceptFunc()
}

// Close implements [net.Listener].
func (fl *FuncListener) Close() error {
	runtimex.Assert(fl.CloseFunc != nil)
	return fl.CloseFunc()
}

// Addr implements [net.Listener].
func (fl *FuncListener) Addr() net.Addr {
	runtimex.Assert(fl.AddrFunc != nil)
	return fl.AddrFunc()
}
