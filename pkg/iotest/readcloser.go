//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/blob/v3.20.1/internal/mocks/reader.go
//

package iotest

import (
	"io"

	"github.com/bassosimone/nettree/pkg/runtimex"
)

// FuncReadCloser allows to mock any [io.ReadCloser].
type FuncReadCloser struct {
	ReadFunc  func(b []byte) (int, error)
	CloseFunc func() error
}

var _ io.ReadCloser = &FuncReadCloser{}

// Read implements [io.Reader].
func (r *FuncReadCloser) Read(b []byte) (int, error) {
	runtimex.Assert(r.ReadFunc != nil)
	return r.ReadFunc(b)
}

// Close implements [io.Closer].
func (r *FuncReadCloser) Close() error {
	runtimex.Assert(r.CloseFunc != nil)
	return r.CloseFunc()
}
