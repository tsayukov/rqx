// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slices"
)

// ErrorStatuses are HTTP error response status codes.
type ErrorStatuses[E error] responseStatuses

// To sets a handler for [ErrorStatuses]. The handler uses [Decoder] to read
// and store decoded [net/http.Response.Body] to the value pointed to by the error
// returned by the handler.
func (e ErrorStatuses[E]) To(decoder Decoder) Option {
	return func(params *doParams) error {
		params.handler.errorResponses = append(params.handler.errorResponses,
			func(resp *http.Response) error {
				if !slices.Contains(e, resp.StatusCode) {
					return nil
				}

				var resultError E
				if err := decoder(resp.Body, &resultError); err != nil {
					return err
				}

				return resultError
			},
		)

		return nil
	}
}

// ToJSON sets a handler for [ErrorStatuses]. The handler reads and stores
// JSON-decoded [net/http.Response.Body] to the value pointed to by the error
// returned by the handler.
func (e ErrorStatuses[E]) ToJSON() Option {
	return e.To(jsonDecoder)
}

// ToXML sets a handler for [ErrorStatuses]. The handler reads and stores
// XML-decoded [net/http.Response.Body] to the value pointed to by the error
// returned by the handler.
func (e ErrorStatuses[E]) ToXML() Option {
	return e.To(xmlDecoder)
}

// UnhandledResponseError is an error for the response that did not match
// any handlers.
type UnhandledResponseError struct {
	status  int
	headers http.Header
	body    *bytes.Buffer
}

func newUnhandledResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return &UnhandledResponseError{
		status:  resp.StatusCode,
		headers: resp.Header.Clone(),
		body:    bytes.NewBuffer(body),
	}
}

func (u *UnhandledResponseError) Error() string {
	return fmt.Sprintf(
		"unhandled response with status %d:\n\theader: %#v\n\tbody: %s",
		u.status, u.headers, u.body.String(),
	)
}

var _ error = (*UnhandledResponseError)(nil)
