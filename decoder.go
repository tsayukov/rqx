// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

// Decoder reads from [io.Reader] and stores its decoded content
// to the value pointed to by the given interface.
type Decoder func(from io.Reader, to any) error

func jsonDecoder(from io.Reader, to any) error {
	return json.NewDecoder(from).Decode(to)
}

func xmlDecoder(from io.Reader, to any) error {
	return xml.NewDecoder(from).Decode(to)
}
