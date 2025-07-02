// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"net/http"
	"slices"

	"github.com/tsayukov/optparams"
)

// ErrorStatuses are HTTP error response status codes.
type ErrorStatuses responseStatuses

// To sets a handler for [ErrorStatuses]. The handler uses [Decoder] to read
// and store decoded [net/http.Response.Body] to the value pointed to by the given
// resultError.
func (e ErrorStatuses) To(resultError error, decoder Decoder) optparams.Func[doParams] {
	return func(params *doParams) error {
		params.handler.errorResponses = append(params.handler.errorResponses,
			func(resp *http.Response) error {
				if !slices.Contains(e, resp.StatusCode) {
					return nil
				}

				if err := decoder(resp.Body, resultError); err != nil {
					return err
				}

				return resultError
			},
		)

		return nil
	}
}

// ToJSON sets a handler for [ErrorStatuses]. The handler reads and stores
// JSON-decoded [net/http.Response.Body] to the value pointed to by the given
// resultError.
func (e ErrorStatuses) ToJSON(resultError error) optparams.Func[doParams] {
	return e.To(resultError, jsonDecoder)
}

// ToXML sets a handler for [ErrorStatuses]. The handler reads and stores
// XML-decoded [net/http.Response.Body] to the value pointed to by the given
// resultError.
func (e ErrorStatuses) ToXML(resultError error) optparams.Func[doParams] {
	return e.To(resultError, xmlDecoder)
}
