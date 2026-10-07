//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/blob/v3.20.1/internal/mocks/http_test.go
//

package httptestx

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFuncRoundTripper(t *testing.T) {
	wantErr := errors.New("mocked error")
	transport := &FuncRoundTripper{
		RoundTripFunc: func(*http.Request) (*http.Response, error) {
			return nil, wantErr
		},
	}
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	resp, err := transport.RoundTrip(req)

	require.Error(t, err)
	require.ErrorIs(t, err, wantErr)
	assert.Nil(t, resp)
}
