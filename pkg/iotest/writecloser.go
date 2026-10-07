//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/blob/v3.20.1/internal/mocks/writer.go
//

package iotest

import (
	"io"

	"github.com/bassosimone/nettree/pkg/runtimex"
)

// FuncWriteCloser allows to mock any [io.WriteCloser].
type FuncWriteCloser struct {
	WriteFunc func(b []byte) (int, error)
	CloseFunc func() error
}

var _ io.WriteCloser = &FuncWriteCloser{}

// Write implements [io.Writer].
func (w *FuncWriteCloser) Write(b []byte) (int, error) {
	runtimex.Assert(w.WriteFunc != nil)
	return w.WriteFunc(b)
}

// Close implements [io.Closer].
func (w *FuncWriteCloser) Close() error {
	runtimex.Assert(w.CloseFunc != nil)
	return w.CloseFunc()
}
