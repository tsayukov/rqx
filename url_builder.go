// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"strconv"
	"strings"

	querypkg "github.com/google/go-querystring/query"
)

// FromInt returns the string representation of the given integer value.
// Use it to construct a URL with path parameters.
func FromInt[T interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}](value T) string {
	return strconv.FormatInt(int64(value), 10)
}

// FromUint returns the string representation of the given unsigned integer
// value. Use it to construct a URL with path parameters.
func FromUint[T interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}](value T) string {
	return strconv.FormatUint(uint64(value), 10)
}

type urlBuilder struct {
	length  int
	paths   []string
	queries []string
}

func (u *urlBuilder) appendPaths(paths ...string) error {
	for _, p := range paths {
		trimmedPath := strings.Trim(p, "/")
		u.length += 1 + len(trimmedPath)
		u.paths = append(u.paths, trimmedPath)
	}

	return nil
}

func (u *urlBuilder) appendQuery(data any) error {
	if data == nil {
		return nil
	}

	values, err := querypkg.Values(data)
	if err != nil {
		return err
	}

	query := values.Encode()
	u.length += 1 + len(query)
	u.queries = append(u.queries, query)

	return nil
}

func (u *urlBuilder) build(base string) string {
	var url strings.Builder

	base = strings.TrimRight(base, "/")

	url.Grow(len(base) + u.length)

	url.WriteString(base)

	for _, p := range u.paths {
		url.WriteRune('/')
		url.WriteString(p)
	}

	if len(u.queries) == 0 {
		return url.String()
	}

	url.WriteRune('?')
	url.WriteString(u.queries[0])

	for _, q := range u.queries[1:] {
		url.WriteRune('&')
		url.WriteString(q)
	}

	return url.String()
}
