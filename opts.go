// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

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

var ErrBodyAlreadyExists = errors.New("body already exists")

// WithBody adds the given data as the body content. If the body is already set,
// it causes the [ErrBodyAlreadyExists] error.
func WithBody(data io.Reader) optparams.Func[doParams] {
	return func(params *doParams) error {
		if params.body != nil {
			return ErrBodyAlreadyExists
		}

		params.body = data

		return nil
	}
}

// WithBytes adds the given bytes as the body content. If the body is already
// set, it causes the [ErrBodyAlreadyExists] error.
func WithBytes(data []byte) optparams.Func[doParams] {
	return func(params *doParams) error {
		if params.body != nil {
			return ErrBodyAlreadyExists
		}

		params.body = bytes.NewReader(data)

		return nil
	}
}

// WithTextPlain adds the given text as the body content and sets the content
// type as "text/plain". If the body is already set, it causes
// the [ErrBodyAlreadyExists] error.
func WithTextPlain(data string) optparams.Func[doParams] {
	return optparams.Join[doParams](
		func(params *doParams) error {
			if params.body != nil {
				return ErrBodyAlreadyExists
			}

			params.body = strings.NewReader(data)

			return nil
		},
		WithContentType(string(ContentTextPlain)),
	)
}

// WithJSON encodes the given data in JSON format as the body content and sets
// the content type as "application/json". If the body is already set, it causes
// the [ErrBodyAlreadyExists] error.
func WithJSON(data any) optparams.Func[doParams] {
	return optparams.Join[doParams](
		func(params *doParams) error {
			if params.body != nil {
				return ErrBodyAlreadyExists
			}

			var buffer bytes.Buffer
			if err := json.NewEncoder(&buffer).Encode(data); err != nil {
				return err
			}
			params.body = bytes.NewReader(buffer.Bytes())

			return nil
		},
		WithContentType(string(ContentJSON)),
	)
}

// WithXML encodes the given data in XML format as the body content and sets
// the content type as "application/xml". If the body is already set, it causes
// the [ErrBodyAlreadyExists] error.
func WithXML(data any) optparams.Func[doParams] {
	return optparams.Join[doParams](
		func(params *doParams) error {
			if params.body != nil {
				return ErrBodyAlreadyExists
			}

			var buffer bytes.Buffer
			if err := xml.NewEncoder(&buffer).Encode(data); err != nil {
				return err
			}
			params.body = bytes.NewReader(buffer.Bytes())

			return nil
		},
		WithContentType(string(ContentXML)),
	)
}

// WithMultipartForm returns [MultipartFormBuilder] to add multipart sections
// sequentially before calling the [MultipartFormBuilder.Body] method.
func WithMultipartForm() *MultipartFormBuilder {
	var b MultipartFormBuilder
	b.mw = multipart.NewWriter(&b.buf)
	return &b
}

// WithHandlerBeforeResponse adds the given handler to call it right before
// the sending HTTP request.
func WithHandlerBeforeResponse(handler BeforeResponseHandler) optparams.Func[doParams] {
	return func(params *doParams) error {
		params.handler.beforeResponse = append(params.handler.beforeResponse, handler)
		return nil
	}
}

// WithHandlerAfterResponse adds the given handler to call it immediately after
// receiving non-nil [net/http.Response].
func WithHandlerAfterResponse(handler AfterResponseHandler) optparams.Func[doParams] {
	return func(params *doParams) error {
		params.handler.afterResponse = append(params.handler.afterResponse, handler)
		return nil
	}
}
