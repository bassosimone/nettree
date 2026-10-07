//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// Adapted from: https://github.com/ooni/probe-cli/tree/v3.20.1/internal/mocks
//

package netstub

import (
	"context"

	"github.com/bassosimone/nettree/pkg/runtimex"
)

// FuncResolver allows to mock any resolver with a LookupHost method.
type FuncResolver struct {
	LookupHostFunc func(ctx context.Context, name string) ([]string, error)
}

var _ interface {
	LookupHost(ctx context.Context, name string) ([]string, error)
} = &FuncResolver{}

// LookupHost implements a resolver with a LookupHost method.
func (fr *FuncResolver) LookupHost(ctx context.Context, name string) ([]string, error) {
	runtimex.Assert(fr.LookupHostFunc != nil)
	return fr.LookupHostFunc(ctx, name)
}
