// SPDX-License-Identifier: GPL-3.0-or-later

// Package iox contains small, context-aware extensions for the [io] package.
package iox

import (
	"bytes"
	"context"
	"errors"
	"io"
	"sync"
)

// ErrClosed is returned when writing on a closed [*LockedWriteCloser].
var ErrClosed = errors.New("locked writer is closed")

// LockedWriteCloser is a concurrency safe [io.WriteCloser] wrapper.
//
// It serializes writes, makes Close idempotent, and keeps track of the number of
// bytes successfully written.
//
// All methods are safe for concurrent use.
//
// Close is serialized with Write, so it may block until an in-flight Write returns.
//
// Construct using [NewLockedWriteCloser].
type LockedWriteCloser struct {
	err error
	mu  sync.RWMutex
	num int
	w   io.WriteCloser
}

// NewLockedWriteCloser wraps an [io.WriteCloser] and returns a concurrency-safe wrapper.
func NewLockedWriteCloser(w io.WriteCloser) *LockedWriteCloser {
	return &LockedWriteCloser{w: w}
}

// LockedWrite writes the given bytes to the underlying [io.WriteCloser].
//
// The returned error is nil, [ErrClosed] when closed, or the error ocurred
// when attempting to write into the underlying [io.WriteCloser].
func (w *LockedWriteCloser) LockedWrite(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.err; err != nil {
		return 0, err
	}
	count, err := w.w.Write(data)
	w.num += count
	return count, err
}

// Count returns the number of bytes successfully written so far.
func (w *LockedWriteCloser) Count() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.num
}

// Close ensures that subsequent writes would fail with [ErrClosed].
//
// Returns nil, [ErrClosed], or the error occurred when closing the [io.WriteCloser].
func (w *LockedWriteCloser) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.err; err != nil {
		return err
	}
	w.err = ErrClosed
	return w.w.Close()
}

// writeAdapter adapts [*LockedWriteCloser] to be an [io.Writer].
type writerAdapter struct {
	w *LockedWriteCloser
}

// Write implements [io.Writer].
func (w writerAdapter) Write(buf []byte) (int, error) {
	return w.w.LockedWrite(buf)
}

// CopyContext is a context-interruptible variant of [io.Copy].
//
// It copies from rc into lwc in a background goroutine. On return, it always
// closes lwc. It only closes rc when the context is canceled, to unblock any
// in-flight Read in the background goroutine.
//
// On success, rc is NOT closed. The caller MUST ensure rc is closed after
// CopyContext returns (e.g., via defer).
//
// The returned error is either caused by I/O or by the context.
func CopyContext(ctx context.Context, lwc *LockedWriteCloser, rc io.ReadCloser) (int, error) {
	// 1. prepare for receiving the background read result
	errch := make(chan error, 1)

	// 2. do in background so we can be interrupted
	go func() {
		_, err := io.Copy(writerAdapter{lwc}, rc)
		errch <- err
	}()

	// 3. wait and collect the error
	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
		// 4a. close the reader to unblock the goroutine's Read
		rc.Close()
	case err = <-errch:
		// 4b. success: do NOT close rc, the caller is responsible
	}

	// 5. always close the writer so the byte count is stable
	lwc.Close()

	// 6. access the number of bytes written once we have closed the
	// writer, so the number is stable ("happens after").
	count := lwc.Count()

	// 7. return to the caller
	return count, err
}

// ReadAllContext is a context-interruptible variant of [io.ReadAll].
//
// It reads from rc into an internal buffer using [CopyContext]. On success, rc
// is NOT closed (the caller MUST close it, e.g., via defer). On context cancellation,
// rc is closed to unblock any in-flight Read.
//
// The returned byte slice contains whatever was read before the error (if any),
// so partial results are available even when the context is canceled.
func ReadAllContext(ctx context.Context, rc io.ReadCloser) ([]byte, error) {
	buf := &bytes.Buffer{}
	_, err := CopyContext(ctx, NewLockedWriteCloser(NopWriteCloser(buf)), rc)
	return buf.Bytes(), err
}

// NopWriteCloser wraps an [io.Writer] and returns a no-op [io.WriteCloser].
//
// This is useful when [CopyContext] needs to stream into a writer that does not
// require closing, such as a [*bytes.Buffer].
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}

// nopWriteCloser adapts an [io.Writer] to become an [io.WriteCloser].
type nopWriteCloser struct {
	io.Writer
}

// Close implements [io.Closer].
func (w nopWriteCloser) Close() error {
	return nil
}

// LimitReadCloser wraps rc such that reads are limited to n bytes
// while Close forwards to the underlying rc.
func LimitReadCloser(rc io.ReadCloser, n int64) io.ReadCloser {
	return readCloser{io.LimitReader(rc, n), rc}
}

// readCloser adapts an [io.Reader] plus an [io.Closer] to an [io.ReadCloser].
type readCloser struct {
	io.Reader
	io.Closer
}
