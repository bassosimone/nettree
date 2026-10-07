//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/blob/v3.20.1/internal/mocks/writer.go
//

package iotest

import "io"

// FuncWriter allows to mock any [io.Writer].
type FuncWriter = FuncWriteCloser

var _ io.Writer = &FuncWriter{}
