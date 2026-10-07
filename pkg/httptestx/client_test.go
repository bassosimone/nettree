//
// SPDX-License-Identifier: GPL-3.0-or-later
//

package httptestx

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFuncClient(t *testing.T) {
	wantErr := errors.New("mocked error")
	client := &FuncClient{
		DoFunc: func(*http.Request) (*http.Response, error) {
			return nil, wantErr
		},
	}
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)

	require.Error(t, err)
	require.ErrorIs(t, err, wantErr)
	assert.Nil(t, resp)
}
