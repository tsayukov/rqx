// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"net/http"
	"slices"

	"github.com/tsayukov/optparams"
)

// OKStatuses are HTTP response status codes that are successful.
type OKStatuses responseStatuses

// To sets a handler for [OKStatuses]. The handler uses [Decoder] to read
// and store decoded [net/http.Response.Body] to the value
// pointed to by the given result.
func (o OKStatuses) To(result any, decoder Decoder) optparams.Func[doParams] {
	return func(params *doParams) error {
		params.handler.okResponse = func(resp *http.Response) (any, error) {
			if !slices.Contains(o, resp.StatusCode) {
				return nil, nil
			}

			if err := decoder(resp.Body, result); err != nil {
				return nil, err
			}

			return result, nil
		}

		return nil
	}
}

// ToJSON sets a handler for [OKStatuses]. The handler reads and stores
// JSON-decoded [net/http.Response.Body] to the value pointed to by the given
// result.
func (o OKStatuses) ToJSON(result any) optparams.Func[doParams] {
	return o.To(result, jsonDecoder)
}

// ToXML sets a handler for [OKStatuses]. The handler reads and stores
// XML-decoded [net/http.Response.Body] to the value pointed to by the given
// result.
func (o OKStatuses) ToXML(result any) optparams.Func[doParams] {
	return o.To(result, xmlDecoder)
}
