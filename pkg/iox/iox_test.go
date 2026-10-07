// SPDX-License-Identifier: GPL-3.0-or-later

package iox

import (
	"bytes"
	"context"
	"io"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/bassosimone/nettree/pkg/iotest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyContextWithCancelledContext(t *testing.T) {
	// Create a reader that blocks until Close is called.
	insideReader := make(chan struct{})
	unblockReader := make(chan struct{})
	closeCalled := &atomic.Bool{}
	rc := &iotest.FuncReadCloser{
		ReadFunc: func(b []byte) (int, error) {
			close(insideReader)
			<-unblockReader
			return 0, io.EOF
		},
		CloseFunc: func() error {
			closeCalled.Store(true)
			close(unblockReader)
			return nil
		},
	}

	buff := &bytes.Buffer{}
	lwc := NewLockedWriteCloser(NopWriteCloser(buff))

	// Create a cancellable context.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Arrange to cancel the context once the reader has entered.
	go func() {
		<-insideReader
		cancel()
	}()

	// Invoke CopyContext with a reader that blocks on Read.
	count, err := CopyContext(ctx, lwc, rc)

	// Because the context is canceled we should see zero and error.
	require.True(t, count == 0)
	require.ErrorIs(t, err, context.Canceled)

	// CopyContext MUST close the reader on cancellation to unblock the goroutine.
	assert.True(t, closeCalled.Load())
	assert.Equal(t, 0, buff.Len())

	// Make sure that subsequent calls to write are not propagated
	// and that we're getting the expected error.
	count, err = lwc.LockedWrite([]byte("0xabad1dea"))
	require.True(t, count == 0)
	require.ErrorIs(t, err, ErrClosed)

	// Make sure subsequent calls to close are not propagated
	// and that we're getting the expected error.
	err = lwc.Close()
	require.ErrorIs(t, err, ErrClosed)
}

func TestCopyContextSuccess(t *testing.T) {
	// Create a reader that returns the full payload and then EOF.
	const payload = "hello from iox"
	closeCalled := &atomic.Bool{}
	rc := &iotest.FuncReadCloser{
		ReadFunc: strings.NewReader(payload).Read,
		CloseFunc: func() error {
			closeCalled.Store(true)
			return nil
		},
	}

	buff := &bytes.Buffer{}
	lwc := NewLockedWriteCloser(NopWriteCloser(buff))

	count, err := CopyContext(context.Background(), lwc, rc)
	require.NoError(t, err)
	assert.Equal(t, len(payload), count)
	assert.Equal(t, payload, buff.String())

	// CopyContext MUST NOT close the reader on success.
	assert.False(t, closeCalled.Load())
}

func TestReadAllContextWithCancelledContext(t *testing.T) {
	// Create a reader that blocks until Close is called.
	insideReader := make(chan struct{})
	unblockReader := make(chan struct{})
	closeCalled := &atomic.Bool{}
	rc := &iotest.FuncReadCloser{
		ReadFunc: func(b []byte) (int, error) {
			close(insideReader)
			<-unblockReader
			return 0, io.EOF
		},
		CloseFunc: func() error {
			closeCalled.Store(true)
			close(unblockReader)
			return nil
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-insideReader
		cancel()
	}()

	data, err := ReadAllContext(ctx, rc)
	require.ErrorIs(t, err, context.Canceled)
	assert.Empty(t, data)
	assert.True(t, closeCalled.Load())
}

func TestReadAllContextSuccess(t *testing.T) {
	const payload = "hello from iox"
	closeCalled := &atomic.Bool{}
	rc := &iotest.FuncReadCloser{
		ReadFunc: strings.NewReader(payload).Read,
		CloseFunc: func() error {
			closeCalled.Store(true)
			return nil
		},
	}

	data, err := ReadAllContext(context.Background(), rc)
	require.NoError(t, err)
	assert.Equal(t, payload, string(data))

	// ReadAllContext MUST NOT close the reader on success.
	assert.False(t, closeCalled.Load())
}

func TestLimitReadCloser(t *testing.T) {
	// Limit reads while keeping the close behavior of the wrapped reader.
	payload := "iox-extra"
	closed := &atomic.Bool{}
	rc := &iotest.FuncReadCloser{
		ReadFunc: strings.NewReader(payload).Read,
		CloseFunc: func() error {
			closed.Store(true)
			return nil
		},
	}
	limited := LimitReadCloser(rc, int64(len("iox")))

	out, err := io.ReadAll(limited)
	require.NoError(t, err)
	assert.Equal(t, "iox", string(out))

	err = limited.Close()
	require.NoError(t, err)
	assert.True(t, closed.Load())
}
