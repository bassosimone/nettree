// SPDX-License-Identifier: GPL-3.0-or-later

package netstub

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFuncListener(t *testing.T) {
	t.Run("Accept", func(t *testing.T) {
		wantErr := errors.New("mocked error")
		lis := &FuncListener{
			AcceptFunc: func() (net.Conn, error) {
				return nil, wantErr
			},
		}

		conn, err := lis.Accept()

		require.ErrorIs(t, err, wantErr)
		require.Nil(t, conn)
	})

	t.Run("Close", func(t *testing.T) {
		wantErr := errors.New("mocked error")
		lis := &FuncListener{
			CloseFunc: func() error {
				return wantErr
			},
		}

		err := lis.Close()

		require.ErrorIs(t, err, wantErr)
	})

	t.Run("Addr", func(t *testing.T) {
		want := &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8080}
		lis := &FuncListener{
			AddrFunc: func() net.Addr {
				return want
			},
		}

		got := lis.Addr()

		require.Equal(t, want, got)
	})
}
