// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"net/textproto"

	"github.com/tsayukov/optparams"
)

type HeaderAppendMode bool

// HeaderAppendModeON makes [net/http.Header] use [net/http.Header.Add]
// instead of [net/http.Header.Set].
const HeaderAppendModeON HeaderAppendMode = true

type withHeaderOptions struct {
	isKeyCanonicalized bool
	doesAddValueToEnd  bool
}

func withHeader(key HeaderKey, value string, options withHeaderOptions) optparams.Func[doParams] {
	canonicalKey := string(key)
	if !options.isKeyCanonicalized {
		canonicalKey = textproto.CanonicalMIMEHeaderKey(canonicalKey)
	}

	return func(params *doParams) error {
		if options.doesAddValueToEnd {
			params.headers[canonicalKey] = append(params.headers[canonicalKey], value)
		} else {
			params.headers[canonicalKey] = []string{value}
		}

		return nil
	}
}
