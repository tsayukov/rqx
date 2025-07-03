// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"bufio"
	"net/textproto"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CanonicalHeaderKeys(t *testing.T) {
	for _, key := range allHeaderKeysFromFile(t, "const.go") {
		assert.Equal(t, textproto.CanonicalMIMEHeaderKey(key), key)
	}
}

func allHeaderKeysFromFile(t *testing.T, filename string) []string {
	t.Helper()

	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	pat := regexp.MustCompile("^\\s*Header\\w+\\s+HeaderKey\\s*=\\s*([\"`])([\\w_-]+)([\"`])")

	var keys []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := pat.FindStringSubmatch(line)
		if len(matches) != 4 {
			continue
		}

		if matches[1] != matches[3] {
			t.Fatalf("invalid string token: %v", matches)
		}

		keys = append(keys, matches[2])
	}

	return keys
}
