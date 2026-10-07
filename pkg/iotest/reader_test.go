//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/blob/v3.20.1/internal/mocks/reader_test.go
//

package iotest

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFuncReader(t *testing.T) {
	r := &FuncReader{
		ReadFunc: func(b []byte) (int, error) {
			return 0, errors.New("mocked error")
		},
	}
	b := make([]byte, 128)
	count, err := r.Read(b)
	require.Error(t, err)
	assert.Equal(t, 0, count)
}
