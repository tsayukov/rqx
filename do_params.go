// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/tsayukov/optparams"
)

// doParams holds required and optional arguments of [Do].
type doParams struct {
	ctx        context.Context
	client     *http.Client
	urlBuilder urlBuilder
	headers    http.Header
	body       io.Reader
	handler    handler
}

func newDoParams(opts ...Option) (*doParams, error) {
	params := &doParams{
		headers: make(http.Header),
	}

	opts = append(opts,
		optparams.Default[doParams](&params.ctx, context.Background()),
		optparams.Default[doParams](&params.client, http.DefaultClient),
	)

	if err := optparams.Apply(params, opts...); err != nil {
		return nil, err
	}

	if params.handler.rateLimitResponse != nil && params.body != nil {
		_, ok := params.body.(io.Closer)
		if ok { // if the body is io.Closer
			return nil, errors.New("rate limit handler cannot be set if body is io.Closer")
		}
	}

	return params, nil
}
