// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"errors"
	"io"
	"net/http"
)

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
func Do(httpMethod HTTPMethod, url string, opts ...Option) (retErr error) {
	params, err := newDoParams(opts...)
	if err != nil {
		return err
	}

	url = params.urlBuilder.build(url)

	for {
		req, err := prepareRequest(httpMethod, url, params)
		if err != nil {
			return err
		}

		if err := params.handler.applyBefore(req); err != nil {
			return err
		}

		resp, err := params.client.Do(req)
		if err != nil {
			return err
		}

		// We need to close the body manually because we are operating
		// in the loop.
		closeBody := func(body io.ReadCloser) {
			if err := body.Close(); err != nil {
				retErr = errors.Join(retErr, err)
			}
		}

		if err := params.handler.applyAfter(resp); err != nil {
			closeBody(resp.Body)
			return err
		}

		match, err := params.handler.matchOK(resp)
		if match {
			closeBody(resp.Body)
			return err // nil or error
		}

		if err := params.handler.matchError(resp); err != nil {
			if errors.Is(err, errRateLimit) && params.handler.rateLimitResponse != nil {
				if err := params.handler.rateLimitResponse(params.ctx, resp); err != nil {
					return err
				}

				closeBody(resp.Body)
				continue
			}

			closeBody(resp.Body)
			return err
		}

		unhandledError := newUnhandledResponse(resp)
		closeBody(resp.Body)
		return unhandledError
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
