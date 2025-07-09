// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_urlBuilder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		urlFunc  func() (string, error)
		want     string
		hasError bool
	}{
		{
			name: "Empty",
			urlFunc: func() (string, error) {
				u := &urlBuilder{}
				return u.build(""), nil
			},
			want: "",
		},
		{
			name: "Only base stripped URL",
			urlFunc: func() (string, error) {
				u := &urlBuilder{}
				return u.build("https://www.example.com"), nil
			},
			want: "https://www.example.com",
		},
		{
			name: "Only base URL with slash",
			urlFunc: func() (string, error) {
				u := &urlBuilder{}
				return u.build("https://www.example.com/"), nil
			},
			want: "https://www.example.com",
		},
		{
			name: "URL with no path",
			urlFunc: func() (string, error) {
				u := &urlBuilder{}
				if err := u.appendPaths(); err != nil {
					return "", err
				}
				return u.build("https://www.example.com"), nil
			},
			want: "https://www.example.com",
		},
		{
			name: "URL with stripped path",
			urlFunc: func() (string, error) {
				u := &urlBuilder{}
				if err := u.appendPaths("one"); err != nil {
					return "", err
				}
				return u.build("https://www.example.com"), nil
			},
			want: "https://www.example.com/one",
		},
		{
			name: "URL with path with slashes",
			urlFunc: func() (string, error) {
				u := &urlBuilder{}
				if err := u.appendPaths("/one/"); err != nil {
					return "", err
				}
				return u.build("https://www.example.com"), nil
			},
			want: "https://www.example.com/one",
		},
		{
			name: "URL with stripped paths",
			urlFunc: func() (string, error) {
				u := &urlBuilder{}
				if err := u.appendPaths("one", "two", "three/four"); err != nil {
					return "", err
				}
				return u.build("https://www.example.com"), nil
			},
			want: "https://www.example.com/one/two/three/four",
		},
		{
			name: "URL with paths with slashes",
			urlFunc: func() (string, error) {
				u := &urlBuilder{}
				if err := u.appendPaths("/one", "two/", "/three/four"); err != nil {
					return "", err
				}
				return u.build("https://www.example.com"), nil
			},
			want: "https://www.example.com/one/two/three/four",
		},
		{
			name: "URL with nil query",
			urlFunc: func() (string, error) {
				u := &urlBuilder{}
				if err := u.appendQuery(nil); err != nil {
					return "", err
				}
				return u.build("https://www.example.com"), nil
			},
			want: "https://www.example.com",
		},
		{
			name: "URL with query",
			urlFunc: func() (string, error) {
				data := struct {
					First  string   `url:"first"`
					Second []string `url:"second,brackets"`
				}{
					First:  "1",
					Second: []string{"2", "3", "4", "5"},
				}

				u := &urlBuilder{}
				if err := u.appendQuery(&data); err != nil {
					return "", err
				}
				return u.build("https://www.example.com"), nil
			},
			want: "https://www.example.com?first=1&second%5B%5D=2&second%5B%5D=3&second%5B%5D=4&second%5B%5D=5",
		},
		{
			name: "URL with error query",
			urlFunc: func() (string, error) {
				u := &urlBuilder{}
				if err := u.appendQuery(42); err != nil {
					return "", err
				}
				return u.build("https://www.example.com"), nil
			},
			hasError: true,
		},
		{
			name: "URL with paths and query",
			urlFunc: func() (string, error) {
				data := struct {
					First  string   `url:"first"`
					Second []string `url:"second,brackets"`
				}{
					First:  "1",
					Second: []string{"2"},
				}

				u := &urlBuilder{}
				if err := u.appendQuery(&data); err != nil {
					return "", err
				}
				if err := u.appendPaths("/one/two/three/four/"); err != nil {
					return "", err
				}

				return u.build("https://www.example.com"), nil
			},
			want: "https://www.example.com/one/two/three/four?first=1&second%5B%5D=2",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			url, err := tt.urlFunc()

			if tt.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, url)
			}
		})
	}
}

func Test_FromInt(t *testing.T) {
	assert.Equal(t, "42", FromInt(42))
	assert.Equal(t, "42", FromInt(int8(42)))
	assert.Equal(t, "42", FromInt(int16(42)))
	assert.Equal(t, "42", FromInt(int32(42)))
	assert.Equal(t, "42", FromInt(int64(42)))
}

func Test_FromUint(t *testing.T) {
	assert.Equal(t, "42", FromUint(uint(42)))
	assert.Equal(t, "42", FromUint(uint8(42)))
	assert.Equal(t, "42", FromUint(uint16(42)))
	assert.Equal(t, "42", FromUint(uint32(42)))
	assert.Equal(t, "42", FromUint(uint64(42)))
}
