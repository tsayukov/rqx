// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/tsayukov/optparams"
)

func optionalBool[B ~bool](value ...B) bool {
	return bool(len(value) > 0 && value[0])
}

// WithContext sets the given [context.Context] for the current request.
func WithContext(ctx context.Context) optparams.Func[doParams] {
	return func(params *doParams) error {
		params.ctx = ctx
		return nil
	}
}

// WithClient sets the given [net/http.Client] for the current request.
func WithClient(c *http.Client) optparams.Func[doParams] {
	return func(params *doParams) error {
		params.client = c
		return nil
	}
}

// WithURLPaths appends the given paths separated by '/' to the URL. Note that
// the resulting URL is not escaped.
func WithURLPaths(paths ...string) optparams.Func[doParams] {
	return func(params *doParams) error {
		return params.urlBuilder.appendPaths(paths...)
	}
}

// WithQuery adds a properly escaped query string encoded from the given data.
func WithQuery(data any) optparams.Func[doParams] {
	return func(params *doParams) error {
		return params.urlBuilder.appendQuery(data)
	}
}

func WithHeader(key HeaderKey, value string, appendMode ...HeaderAppendMode) optparams.Func[doParams] {
	return withHeader(key, value, withHeaderOptions{
		isKeyCanonicalized: false,
		doesAddValueToEnd:  optionalBool(appendMode...),
	})
}

// WithContentType sets the HTTP Content-Type representation header, overwriting
// the previous one, if any.
func WithContentType(value string, appendMode ...HeaderAppendMode) optparams.Func[doParams] {
	return withHeader(HeaderContentType, value, withHeaderOptions{
		isKeyCanonicalized: true,
		doesAddValueToEnd:  optionalBool(appendMode...),
	})
}

// WithAccept sets the HTTP Accept request header, overwriting the previous one,
// if any.
func WithAccept(value string, appendMode ...HeaderAppendMode) optparams.Func[doParams] {
	return withHeader(HeaderAccept, value, withHeaderOptions{
		isKeyCanonicalized: true,
		doesAddValueToEnd:  optionalBool(appendMode...),
	})
}

// WithAuth sets the HTTP Authorization request header with the given value.
func WithAuth(value string, appendMode ...HeaderAppendMode) optparams.Func[doParams] {
	return withHeader(HeaderAuthorization, value, withHeaderOptions{
		isKeyCanonicalized: true,
		doesAddValueToEnd:  optionalBool(appendMode...),
	})
}

// WithBasicAuth sets the HTTP Authorization header to use HTTP Basic Authentication
// with the provided username and password.
func WithBasicAuth(username, password string) optparams.Func[doParams] {
	enc := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	return WithAuth("Basic " + enc)
}