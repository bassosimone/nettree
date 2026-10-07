//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/tree/v3.20.1/internal/mocks
//

package netstub

import (
	"net"
	"time"

	"github.com/bassosimone/nettree/pkg/runtimex"
)

// FuncConn allows to mock any [net.Conn].
type FuncConn struct {
	ReadFunc             func([]byte) (int, error)
	WriteFunc            func([]byte) (int, error)
	CloseFunc            func() error
	LocalAddrFunc        func() net.Addr
	RemoteAddrFunc       func() net.Addr
	SetDeadlineFunc      func(time.Time) error
	SetReadDeadlineFunc  func(time.Time) error
	SetWriteDeadlineFunc func(time.Time) error
}

var _ net.Conn = &FuncConn{}

// Read implements [net.Conn].
func (fc *FuncConn) Read(b []byte) (int, error) {
	runtimex.Assert(fc.ReadFunc != nil)
	return fc.ReadFunc(b)
}

// Write implements [net.Conn].
func (fc *FuncConn) Write(b []byte) (int, error) {
	runtimex.Assert(fc.WriteFunc != nil)
	return fc.WriteFunc(b)
}

// Close implements [net.Conn].
func (fc *FuncConn) Close() error {
	runtimex.Assert(fc.CloseFunc != nil)
	return fc.CloseFunc()
}

// LocalAddr implements [net.Conn].
func (fc *FuncConn) LocalAddr() net.Addr {
	runtimex.Assert(fc.LocalAddrFunc != nil)
	return fc.LocalAddrFunc()
}

// RemoteAddr implements [net.Conn].
func (fc *FuncConn) RemoteAddr() net.Addr {
	runtimex.Assert(fc.RemoteAddrFunc != nil)
	return fc.RemoteAddrFunc()
}

// SetDeadline implements [net.Conn].
func (fc *FuncConn) SetDeadline(t time.Time) error {
	runtimex.Assert(fc.SetDeadlineFunc != nil)
	return fc.SetDeadlineFunc(t)
}

// SetReadDeadline implements [net.Conn].
func (fc *FuncConn) SetReadDeadline(t time.Time) error {
	runtimex.Assert(fc.SetReadDeadlineFunc != nil)
	return fc.SetReadDeadlineFunc(t)
}

// SetWriteDeadline implements [net.Conn].
func (fc *FuncConn) SetWriteDeadline(t time.Time) error {
	runtimex.Assert(fc.SetWriteDeadlineFunc != nil)
	return fc.SetWriteDeadlineFunc(t)
}
