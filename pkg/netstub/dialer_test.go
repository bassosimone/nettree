//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/tree/v3.20.1/internal/mocks
//

package netstub

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFuncDialer(t *testing.T) {
	wantErr := errors.New("mocked error")
	dialer := &FuncDialer{
		DialContextFunc: func(context.Context, string, string) (net.Conn, error) {
			return nil, wantErr
		},
	}

	conn, err := dialer.DialContext(context.Background(), "tcp", "example.com:80")

	require.ErrorIs(t, err, wantErr)
	require.Nil(t, conn)
}
