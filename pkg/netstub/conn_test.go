//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/tree/v3.20.1/internal/mocks
//

package netstub

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFuncConn(t *testing.T) {
	wantErr := errors.New("mocked error")
	wantAddr := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)}
	conn := &FuncConn{
		ReadFunc: func([]byte) (int, error) {
			return 0, wantErr
		},
		WriteFunc: func([]byte) (int, error) {
			return 0, wantErr
		},
		CloseFunc: func() error {
			return wantErr
		},
		LocalAddrFunc: func() net.Addr {
			return wantAddr
		},
		RemoteAddrFunc: func() net.Addr {
			return wantAddr
		},
		SetDeadlineFunc: func(time.Time) error {
			return wantErr
		},
		SetReadDeadlineFunc: func(time.Time) error {
			return wantErr
		},
		SetWriteDeadlineFunc: func(time.Time) error {
			return wantErr
		},
	}

	buf := make([]byte, 8)
	_, err := conn.Read(buf)
	require.ErrorIs(t, err, wantErr)

	_, err = conn.Write(buf)
	require.ErrorIs(t, err, wantErr)

	err = conn.Close()
	require.ErrorIs(t, err, wantErr)

	assert.Equal(t, wantAddr, conn.LocalAddr())
	assert.Equal(t, wantAddr, conn.RemoteAddr())

	deadline := time.Now()
	err = conn.SetDeadline(deadline)
	require.ErrorIs(t, err, wantErr)

	err = conn.SetReadDeadline(deadline)
	require.ErrorIs(t, err, wantErr)

	err = conn.SetWriteDeadline(deadline)
	require.ErrorIs(t, err, wantErr)
}
