// This file is licensed under the terms of the MIT License (see LICENSE file)
// Copyright (c) 2025 Pavel Tsayukov p.tsayukov@gmail.com

package rqx

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"
)

// MultipartFormBuilder is a builder to constructs consecutive multipart
// sections.
type MultipartFormBuilder struct {
	mw   *multipart.Writer
	buf  bytes.Buffer
	errs []error
}

func (b *MultipartFormBuilder) joinErrors(errs ...error) *MultipartFormBuilder {
	b.errs = append(b.errs, errs...)
	return b
}

func (b *MultipartFormBuilder) writePart(w io.Writer, r io.Reader) *MultipartFormBuilder {
	if _, err := io.Copy(w, r); err != nil {
		return b.joinErrors(err)
	}

	return b
}

// AddString adds a new multipart section with a header using the given field
// name and writes the content to the section's body.
func (b *MultipartFormBuilder) AddString(fieldName, content string) *MultipartFormBuilder {
	w, err := b.mw.CreateFormField(fieldName)
	if err != nil {
		return b.joinErrors(err)
	}

	return b.writePart(w, strings.NewReader(content))
}

// AddFile adds a new multipart section with a header using the given field name
// and writes the file content to the section's body.
func (b *MultipartFormBuilder) AddFile(fieldName string, file *os.File) *MultipartFormBuilder {
	return b.AddAsFile(fieldName, file, file.Name())
}

// AddAsFile adds a new multipart section with a header using the given field
// name and writes the content to the section's body as if it was a file with
// the given file name.
func (b *MultipartFormBuilder) AddAsFile(
	fieldName string,
	content io.Reader,
	fileName string,
) *MultipartFormBuilder {
	if closer, ok := content.(io.Closer); ok {
		defer func() { _ = closer.Close() }()
	}

	w, err := b.mw.CreateFormFile(fieldName, fileName)
	if err != nil {
		return b.joinErrors(err)
	}

	return b.writePart(w, content)
}

var quoteEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// AddAsFileWithType adds a new multipart section with a header using the given
// field name and writes the content to the section's body as if it was a file
// with the given file name and content type.
func (b *MultipartFormBuilder) AddAsFileWithType(
	fieldName string,
	content io.Reader,
	fileName, contentType string,
) *MultipartFormBuilder {
	if closer, ok := content.(io.Closer); ok {
		defer func() { _ = closer.Close() }()
	}

	h := make(textproto.MIMEHeader)
	h.Set(string(HeaderContentDisposition), fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
		escapeQuotes(fieldName), escapeQuotes(fileName),
	))
	h.Set(string(HeaderContentType), contentType)

	w, err := b.mw.CreatePart(h)
	if err != nil {
		return b.joinErrors(err)
	}

	return b.writePart(w, content)
}

// Body creates a body with the multipart sections and the proper content type.
func (b *MultipartFormBuilder) Body() Option {
	return func(params *doParams) error {
		if len(b.errs) > 0 {
			return errors.Join(b.errs...)
		}

		if err := b.mw.Close(); err != nil {
			return err
		}

		params.body = bytes.NewReader(b.buf.Bytes())
		params.headers[string(HeaderContentType)] = []string{b.mw.FormDataContentType()}

		return nil
	}
}
