//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/blob/647b03f4270eb758106523fe6273e3ebdbcd599c/internal/runtimex/runtimex.go
//

// Package runtimex contains helpers for code paths that are not expected to fail and
// where failure indicates a programmer error or an unrecoverable condition.
//
// # When to use what
//
// Assert: For enforcing invariants in library code. Document the invariant and its
// justification in a comment above the assertion. Use these for conditions that should
// be impossible if the program is correct.
//
// PanicOnErrorN: For unwrapping `(value, ..., error)` returns where the error cannot
// occur in correct usage but ignoring the error feels sloppy.
//
// LogFatalOnErrorN: In main() functions when you want to log and exit.
//
// # History
//
// This package was originally inspired by [github.com/m-lab/go/rtx].
package runtimex

import (
	"errors"
	"log"
)

// Assert panics if the given value is false. The value passed to
// `panic()` is an error constructed using [errors.New].
//
// You typically use this function to assert runtime invariants in your codebase
// to make it more robust. Document the invariant and its justification in a
// comment above the assertion. For example:
//
//	// Invariant: msg is never nil after successful initialization. This is
//	// guaranteed by the constructor, which validates input.
//	runtimex.Assert(cfg.msg != nil)
//
// The correct approach is to assert for conditions that should be
// impossible if the program is correct.
func Assert(value bool) {
	if !value {
		panic(errors.New("assertion failed"))
	}
}

// PanicOnError0 panics if the given err is not nil. The value passed
// to `panic()` is the given err value.
//
// You typically use this function to assert runtime invariants
// in your codebase to make it more robust. For example:
//
//	runtimex.PanicOnError0(operationThatCannotFail())
//
// The correct approach is to assert for errors that cannot
// possibly happen (e.g., [json.Marshal] applied to a struct
// that can always be marshalled to a JSON string).
func PanicOnError0(err error) {
	if err != nil {
		panic(err)
	}
}

// PanicOnError1 panics if the given err is not nil. The value passed
// to `panic()` is the given err value. Otherwise, it returns the given
// value `v1`.
//
// This code is equivalent to:
//
//	v1, err := fx()
//	if err != nil {
//		panic(err)
//	}
//
// but is more compact and improves readability when chaining operations.
func PanicOnError1[T1 any](v1 T1, err error) T1 {
	if err != nil {
		panic(err)
	}
	return v1
}

// PanicOnError2 panics if the given err is not nil. The value passed
// to `panic()` is the given err value. Otherwise, it returns the given
// values `v1` and `v2`.
//
// This code is equivalent to:
//
//	v1, v2, err := fx()
//	if err != nil {
//		panic(err)
//	}
//
// but is more compact and improves readability when chaining operations.
func PanicOnError2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		panic(err)
	}
	return v1, v2
}

// PanicOnError3 panics if the given err is not nil. The value passed
// to `panic()` is the given err value. Otherwise, it returns the given
// values `v1`, v2, and `v3`.
//
// This code is equivalent to:
//
//	v1, v2, v3, err := fx()
//	if err != nil {
//		panic(err)
//	}
//
// but is more compact and improves readability when chaining operations.
func PanicOnError3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	if err != nil {
		panic(err)
	}
	return v1, v2, v3
}

// LogFatalFunc is the function that [LogFatalOnError0], [LogFatalOnError1],
// [LogFatalOnError2], and [LogFatalOnError3] invoke when they receive a
// non-nil error. The default is [log.Fatal], which logs and calls
// `os.Exit(1)`, bypassing every pending defer.
//
// Replace from `main` before starting goroutines to change the failure
// behavior — for example, to call `deferexit.Panic(1)` after logging so
// that deferred cleanup actually runs. Not safe to mutate concurrently.
var LogFatalFunc = log.Fatal

// LogFatalOnError0 exits with a fatal error if err is not nil.
//
// It is equivalent to:
//
//	if err != nil {
//		log.Fatal(err)
//	}
func LogFatalOnError0(err error) {
	if err != nil {
		LogFatalFunc(err)
	}
}

// LogFatalOnError1 exits with a fatal error if err is not nil. Otherwise,
// it returns the given value `v1`.
//
// It is equivalent to:
//
//	v1, err := fx()
//	if err != nil {
//		log.Fatal(err)
//	}
func LogFatalOnError1[T1 any](v1 T1, err error) T1 {
	if err != nil {
		LogFatalFunc(err)
	}
	return v1
}

// LogFatalOnError2 exits with a fatal error if err is not nil. Otherwise,
// it returns the given values `v1` and `v2`.
//
// It is equivalent to:
//
//	v1, v2, err := fx()
//	if err != nil {
//		log.Fatal(err)
//	}
func LogFatalOnError2[T1, T2 any](v1 T1, v2 T2, err error) (T1, T2) {
	if err != nil {
		LogFatalFunc(err)
	}
	return v1, v2
}

// LogFatalOnError3 exits with a fatal error if err is not nil. Otherwise,
// it returns the given values `v1`, `v2`, and `v3`.
//
// It is equivalent to:
//
//	v1, v2, v3, err := fx()
//	if err != nil {
//		log.Fatal(err)
//	}
func LogFatalOnError3[T1, T2, T3 any](v1 T1, v2 T2, v3 T3, err error) (T1, T2, T3) {
	if err != nil {
		LogFatalFunc(err)
	}
	return v1, v2, v3
}
