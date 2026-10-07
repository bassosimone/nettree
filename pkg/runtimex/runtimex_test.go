// SPDX-License-Identifier: GPL-3.0-or-later

package runtimex

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssert(t *testing.T) {
	t.Run("with true value does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Assert(true)
		})
	})

	t.Run("with false value panics", func(t *testing.T) {
		assert.PanicsWithError(t, "assertion failed", func() {
			Assert(false)
		})
	})
}

func TestPanicOnError0(t *testing.T) {
	t.Run("with nil error does not panic", func(t *testing.T) {
		assert.NotPanics(t, func() {
			PanicOnError0(nil)
		})
	})

	t.Run("with non-nil error panics", func(t *testing.T) {
		expectedErr := errors.New("test error")
		assert.PanicsWithValue(t, expectedErr, func() {
			PanicOnError0(expectedErr)
		})
	})
}

func TestPanicOnError1(t *testing.T) {
	t.Run("with nil error returns value", func(t *testing.T) {
		expectedValue := "test value"
		var result string
		assert.NotPanics(t, func() {
			result = PanicOnError1(expectedValue, nil)
		})
		assert.Equal(t, expectedValue, result)
	})

	t.Run("with non-nil error panics", func(t *testing.T) {
		expectedErr := errors.New("test error")
		assert.PanicsWithValue(t, expectedErr, func() {
			PanicOnError1("value", expectedErr)
		})
	})
}

func TestPanicOnError2(t *testing.T) {
	t.Run("with nil error returns values", func(t *testing.T) {
		expectedV1 := "first"
		expectedV2 := 42
		var v1 string
		var v2 int
		assert.NotPanics(t, func() {
			v1, v2 = PanicOnError2(expectedV1, expectedV2, nil)
		})
		assert.Equal(t, expectedV1, v1)
		assert.Equal(t, expectedV2, v2)
	})

	t.Run("with non-nil error panics", func(t *testing.T) {
		expectedErr := errors.New("test error")
		assert.PanicsWithValue(t, expectedErr, func() {
			PanicOnError2("value1", "value2", expectedErr)
		})
	})
}

func TestPanicOnError3(t *testing.T) {
	t.Run("with nil error returns values", func(t *testing.T) {
		expectedV1 := "first"
		expectedV2 := 42
		expectedV3 := true
		var v1 string
		var v2 int
		var v3 bool
		assert.NotPanics(t, func() {
			v1, v2, v3 = PanicOnError3(expectedV1, expectedV2, expectedV3, nil)
		})
		assert.Equal(t, expectedV1, v1)
		assert.Equal(t, expectedV2, v2)
		assert.Equal(t, expectedV3, v3)
	})

	t.Run("with non-nil error panics", func(t *testing.T) {
		expectedErr := errors.New("test error")
		assert.PanicsWithValue(t, expectedErr, func() {
			PanicOnError3("value1", "value2", "value3", expectedErr)
		})
	})
}

func TestLogFatalOnError(t *testing.T) {
	// Save original LogFatalFunc and restore after each test
	originalLogFatal := LogFatalFunc
	defer func() { LogFatalFunc = originalLogFatal }()

	var fatalCalled bool
	var fatalValue any
	LogFatalFunc = func(v ...any) {
		fatalCalled = true
		fatalValue = v[0]
	}

	// Reset mocks before each subtest
	resetMocks := func() {
		fatalCalled = false
		fatalValue = nil
	}

	t.Run("LogFatalOnError0", func(t *testing.T) {
		t.Run("with nil error", func(t *testing.T) {
			resetMocks()
			LogFatalOnError0(nil)
			assert.False(t, fatalCalled)
		})

		t.Run("with non-nil error", func(t *testing.T) {
			resetMocks()
			err := errors.New("logfatal0")
			LogFatalOnError0(err)
			assert.True(t, fatalCalled)
			assert.Equal(t, err, fatalValue)
		})
	})

	t.Run("LogFatalOnError1", func(t *testing.T) {
		t.Run("with nil error", func(t *testing.T) {
			resetMocks()
			val := LogFatalOnError1(17, nil)
			assert.False(t, fatalCalled)
			assert.Equal(t, 17, val)
		})

		t.Run("with non-nil error", func(t *testing.T) {
			resetMocks()
			err := errors.New("logfatal1")
			_ = LogFatalOnError1(17, err)
			assert.True(t, fatalCalled)
			assert.Equal(t, err, fatalValue)
		})
	})

	t.Run("LogFatalOnError2", func(t *testing.T) {
		t.Run("with nil error", func(t *testing.T) {
			resetMocks()
			v1, v2 := LogFatalOnError2(17, "hi", nil)
			assert.False(t, fatalCalled)
			assert.Equal(t, 17, v1)
			assert.Equal(t, "hi", v2)
		})

		t.Run("with non-nil error", func(t *testing.T) {
			resetMocks()
			err := errors.New("logfatal2")
			_, _ = LogFatalOnError2(17, "hi", err)
			assert.True(t, fatalCalled)
			assert.Equal(t, err, fatalValue)
		})
	})

	t.Run("LogFatalOnError3", func(t *testing.T) {
		t.Run("with nil error", func(t *testing.T) {
			resetMocks()
			v1, v2, v3 := LogFatalOnError3(17, "hi", true, nil)
			assert.False(t, fatalCalled)
			assert.Equal(t, 17, v1)
			assert.Equal(t, "hi", v2)
			assert.Equal(t, true, v3)
		})

		t.Run("with non-nil error", func(t *testing.T) {
			resetMocks()
			err := errors.New("logfatal3")
			_, _, _ = LogFatalOnError3(17, "hi", true, err)
			assert.True(t, fatalCalled)
			assert.Equal(t, err, fatalValue)
		})
	})
}
