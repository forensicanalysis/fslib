// Copyright (c) 2019 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

// Package filetype provides functions to identify file types of binary data.
package filetype

import (
	"bufio"
	"io"
	"strings"

	"github.com/h2non/filetype/types"
)

// ID stores a unique identifier for filetypes.
type ID string

// Filetype represents a single file format.
type Filetype struct {
	ID         ID
	Mimetype   types.MIME
	Extensions []string
	Matcher    func(buf []byte) bool
	layer      int
}

var (
	filetypes = map[ID]*Filetype{}
	exttype   = map[string][]*Filetype{}
)

func init() { //nolint:gochecknoinits
	for _, filetypelist := range [][]*Filetype{filesystemTypes, importedTypes, baseTypes} {
		for _, filetype := range filetypelist {
			filetypes[filetype.ID] = filetype
			for _, ext := range filetype.Extensions {
				if _, ok := exttype[ext]; !ok {
					exttype[ext] = []*Filetype{}
				}
				exttype[ext] = append(exttype[ext], filetype)
			}
		}
	}
}

func newFiletype(id ID, t types.Type, matcher func([]byte) bool, layer int) *Filetype {
	return &Filetype{id, t.MIME, []string{t.Extension}, matcher, layer}
}

// DetectReader identifies a Filetype for an io.Reader.
func DetectReader(r *bufio.Reader) (*Filetype, error) {
	head, err := r.Peek(8192)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return Detect(head), nil
}

// DetectReaderByExtension identifies a Filetype for an io.Reader, the guess is
// used for performance.
func DetectReaderByExtension(r io.Reader, guess string) (*Filetype, error) {
	head := make([]byte, 8192)
	if _, err := r.Read(head); err != nil {
		return nil, err
	}

	guess = strings.TrimLeft(guess, ".")
	if _, ok := exttype[guess]; ok {
		return DetectByExtension(head, guess), nil
	}
	return Detect(head), nil
}

// Detect identifies a Filetype a byte slice.
func Detect(buf []byte) *Filetype {
	if len(buf) == 0 {
		return Empty
	}

	bestMatch := Binary
	for _, ft := range filetypes {
		if ft.Matcher(buf) {
			if ft.layer == 0 {
				return ft
			}
			if ft.layer < bestMatch.layer {
				bestMatch = ft
			}
		}
	}
	return bestMatch
}

// DetectByExtension identifies a Filetype a byte slice, the guess is used for
// performance.
func DetectByExtension(buf []byte, guess string) *Filetype {
	guess = strings.TrimLeft(guess, ".")

	if len(buf) == 0 {
		return Empty
	}

	bestMatch := Binary
	for _, ft := range exttype[guess] {
		if ft.Matcher(buf) {
			if ft.layer == 0 {
				return ft
			}
			if ft.layer < bestMatch.layer {
				bestMatch = ft
			}
		}
	}

	if bestMatch == Binary {
		return Detect(buf)
	}
	return bestMatch
}
