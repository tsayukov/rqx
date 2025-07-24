// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"errors"
	"net/http"
)

// Do sends an HTTP request given [HTTPMethod], URL, and optional parameters.
//
// By default, [context.Background] is used. To set an appropriate context,
// use optional [WithContext].
//
// By default, [net/http.DefaultClient] is used. To set an appropriate
// [net/http.Client], use optional [WithClient].
//
// URL options:
//   - [WithURLPaths];
//   - [WithQuery].
//
// Headers options:
//   - [WithHeader];
//   - [WithContentType];
//   - [WithAccept].
//
// Authorization options:
//   - [WithAuth];
//   - [WithBasicAuth].
//
// Body options:
//   - [WithBody];
//   - [WithBytes];
//   - [WithTextPlain];
//   - [WithJSON];
//   - [WithXML];
//   - [WithMultipartForm].
//
// Handler options:
//   - [WithHandlerBeforeResponse];
//   - [WithHandlerAfterResponse];
//   - [WithOK];
//   - [WithError];
//   - [WithRateLimit].
//
// Error Wrapper options:
//   - [WithErrorPrefix];
//   - [WithErrorWrapper].
func Do(httpMethod HTTPMethod, url string, opts ...Option) error {
	params, err := newDoParams(opts...)
	if err != nil {
		return err
	}

	url = params.urlBuilder.build(url)

	for {
		tryAgain, err := do(httpMethod, url, params)
		if err != nil {
			return err
		}
		if tryAgain {
			continue
		}

		return nil
	}
}

// Get is a shortcut for [Do] for the [GET] HTTP method.
func Get(url string, opts ...Option) error {
	return Do(GET, url, opts...)
}

// Post is a shortcut for [Do] for the [POST] HTTP method.
func Post(url string, opts ...Option) error {
	return Do(POST, url, opts...)
}

// Put is a shortcut for [Do] for the [PUT] HTTP method.
func Put(url string, opts ...Option) error {
	return Do(PUT, url, opts...)
}

// Delete is a shortcut for [Do] for the [DELETE] HTTP method.
func Delete(url string, opts ...Option) error {
	return Do(DELETE, url, opts...)
}

// Options is a shortcut for [Do] for the [OPTIONS] HTTP method.
func Options(url string, opts ...Option) error {
	return Do(OPTIONS, url, opts...)
}

// Patch is a shortcut for [Do] for the [PATCH] HTTP method.
func Patch(url string, opts ...Option) error {
	return Do(PATCH, url, opts...)
}

func prepareRequest(httpMethod HTTPMethod, url string, params *doParams) (*http.Request, error) {
	req, err := http.NewRequestWithContext(params.ctx, string(httpMethod), url, params.body)
	if err != nil {
		return nil, err
	}

	for key, values := range params.headers {
		// No need to call Header.Add() for each value:
		// the key has been already canonicalized.
		req.Header[key] = append(req.Header[key], values...)
	}

	return req, nil
}

func do(httpMethod HTTPMethod, url string, params *doParams) (tryAgain bool, retErr error) {
	req, err := prepareRequest(httpMethod, url, params)
	if err != nil {
		return false, params.errorWrapper(err)
	}

	if err := params.handler.applyBefore(req); err != nil {
		return false, params.errorWrapper(err)
	}

	resp, err := params.client.Do(req)
	if err != nil {
		return false, params.errorWrapper(err)
	}

	defer func() { retErr = errors.Join(retErr, params.errorWrapper(resp.Body.Close())) }()

	if err := params.handler.applyAfter(resp); err != nil {
		return false, params.errorWrapper(err)
	}

	if match, err := params.handler.matchOK(resp); match { // if HTTP statuses are OK
		return false, params.errorWrapper(err) // nil or error
	}

	if err := params.handler.matchError(resp); err != nil {
		if errors.Is(err, errRateLimit) && params.handler.rateLimitResponse != nil {
			if err := params.handler.rateLimitResponse(params.ctx, resp); err != nil {
				return false, params.errorWrapper(err)
			}

			return true, nil
		}

		return false, params.errorWrapper(err)
	}

	return false, params.errorWrapper(newUnhandledResponse(resp))
}
