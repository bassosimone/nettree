//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/tree/v3.20.1/internal/mocks
//

package netstub

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFuncResolver(t *testing.T) {
	wantErr := errors.New("mocked error")
	resolver := &FuncResolver{
		LookupHostFunc: func(context.Context, string) ([]string, error) {
			return nil, wantErr
		},
	}

	addrs, err := resolver.LookupHost(context.Background(), "example.com")

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, addrs)
}
