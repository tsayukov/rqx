// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"context"

	"github.com/tsayukov/optparams"
)

// WithContext sets the given [context.Context] for the current request.
func WithContext(ctx context.Context) optparams.Func[doParams] {
	return func(params *doParams) error {
		params.ctx = ctx
		return nil
	}
}