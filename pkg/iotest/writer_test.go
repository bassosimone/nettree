//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/blob/v3.20.1/internal/mocks/writer_test.go
//

package iotest

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFuncWriter(t *testing.T) {
	w := &FuncWriter{
		WriteFunc: func(b []byte) (int, error) {
			return 0, errors.New("mocked error")
		},
	}
	b := make([]byte, 128)
	count, err := w.Write(b)
	require.Error(t, err)
	assert.Equal(t, count, 0)
}
