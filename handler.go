// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"context"
	"net/http"
)

type (
	handler struct {
		beforeResponse []BeforeResponseHandler
		afterResponse  []AfterResponseHandler

		okResponse     okResponseHandler
		errorResponses []errorResponseHandler

		rateLimitResponse RateLimitHandler
	}

	// BeforeResponseHandler handles [net/http.Request] right before the sending
	// HTTP request.
	BeforeResponseHandler func(*http.Request) error

	// AfterResponseHandler handles [net/http.Response] immediately after
	// receiving non-nil [net/http.Response].
	AfterResponseHandler func(*http.Response) error

	responseStatuses []int

	// okResponseHandler handles [net/http.Response] whose HTTP status code
	// matches one of [OKStatuses].
	okResponseHandler func(*http.Response) (any, error)

	// errorResponseHandler handles [net/http.Response] whose HTTP status code
	// matches one of [ErrorStatuses].
	errorResponseHandler func(*http.Response) error

	// RateLimitHandler handles [net/http.Response] whose HTTP status code
	// matches one of [RateLimitStatuses].
	RateLimitHandler func(ctx context.Context, resp *http.Response) error
)

func (h *handler) applyBefore(req *http.Request) error {
	for _, fn := range h.beforeResponse {
		if err := fn(req); err != nil {
			return err
		}
	}

	return nil
}

func (h *handler) applyAfter(resp *http.Response) error {
	for _, fn := range h.afterResponse {
		if err := fn(resp); err != nil {
			return err
		}
	}

	return nil
}

func (h *handler) matchOK(resp *http.Response) (match bool, _ error) {
	if h.okResponse == nil {
		return false, nil
	}

	result, err := h.okResponse(resp)
	if result != nil || err != nil {
		return true, err
	}

	return false, nil
}

func (h *handler) matchError(resp *http.Response) error {
	for _, errorHandler := range h.errorResponses {
		err := errorHandler(resp)
		if err != nil {
			return err
		}
	}

	return nil
}
