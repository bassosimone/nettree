//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/tree/v3.20.1/internal/mocks
//

package netstub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncAddr(t *testing.T) {
	wantNetwork := "tcp"
	wantString := "127.0.0.1:8080"
	addr := &FuncAddr{
		NetworkFunc: func() string {
			return wantNetwork
		},
		StringFunc: func() string {
			return wantString
		},
	}

	assert.Equal(t, wantNetwork, addr.Network())
	assert.Equal(t, wantString, addr.String())
}
