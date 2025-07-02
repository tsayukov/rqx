// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"errors"
	"net/http"
	"slices"

	"github.com/tsayukov/optparams"
)

// RateLimitStatuses are HTTP response status codes that are returned
// when the rate limit is reached.
type RateLimitStatuses responseStatuses

var errRateLimit = errors.New("rate limit exceeded")

// Cooldown adds the given [RateLimitHandler] to the response handlers.
// Note that when the request body is [io.Closer], [RateLimitHandler]
// is not allowed, because the body will be closed by [net/http.Client.Do]
// before the next attempt.
func (rc RateLimitStatuses) Cooldown(handler RateLimitHandler) optparams.Func[doParams] {
	return func(params *doParams) error {
		if handler == nil {
			return errors.New("rate limit handler is nil")
		}

		params.handler.rateLimitResponse = handler

		params.handler.errorResponses = append(params.handler.errorResponses,
			func(resp *http.Response) error {
				if !slices.Contains(rc, resp.StatusCode) {
					return nil
				}

				return errRateLimit
			})

		return nil
	}
}
