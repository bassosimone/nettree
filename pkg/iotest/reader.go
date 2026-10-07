//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/blob/v3.20.1/internal/mocks/reader.go
//

package iotest

import "io"

// FuncReader allows to mock any [io.Reader].
type FuncReader = FuncReadCloser

var _ io.Reader = &FuncReader{}
