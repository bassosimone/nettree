//
// SPDX-License-Identifier: GPL-3.0-or-later
//

package httptestx

import (
	"net/http"

	"github.com/bassosimone/nettree/pkg/runtimex"
)

// FuncClient allows to mock any Client with a Do method.
type FuncClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do implements the Do method.
func (c *FuncClient) Do(req *http.Request) (*http.Response, error) {
	runtimex.Assert(c.DoFunc != nil)
	return c.DoFunc(req)
}
