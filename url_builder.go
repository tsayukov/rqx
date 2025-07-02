// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"strings"

	querypkg "github.com/google/go-querystring/query"
)

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
