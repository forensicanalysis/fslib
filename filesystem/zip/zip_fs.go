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

// Package zip provides a forensicfs implementation to access zip files.
package zip

import (
	"archive/zip"
	"fmt"
	"github.com/forensicanalysis/fslib/fsio"
	"io/fs"
)

// FS implements a read-only file system for zip files.
type FS struct {
	internal *zip.Reader
}

// New creates a new zip FS.
func New(base fsio.ReadSeekerAt) (*FS, error) {
	size, err := fsio.GetSize(base)
	if err != nil {
		return nil, err
	}

	zr, err := zip.NewReader(base, size)
	if err != nil {
		return nil, err
	}

	return &FS{zr}, nil
}

// Open opens a file for reading.
func (fsys *FS) Open(name string) (fs.File, error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}

	aferoItem, err := fsys.internal.Open(name)
	return &File{name, aferoItem}, err
}
