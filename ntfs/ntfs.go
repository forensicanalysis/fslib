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

// Package ntfs provides a forensicfs implementation of the New Technology File
// System (NTFS).
package ntfs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"

	"www.velocidex.com/golang/go-ntfs/parser"
)

// New creates a new ntfs FS.
func New(r io.ReaderAt) (fs *FS, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error parsing file system as NTFS")
		}
	}()
	reader, err := parser.NewPagedReader(r, 1024*1024, 100*1024*1024)
	if err != nil {
		return nil, err
	}
	ntfsCtx, err := parser.GetNTFSContext(reader, 0)
	return &FS{ntfsCtx: ntfsCtx}, err
}

// FS implements a read-only file system for the NTFS.
type FS struct {
	ntfsCtx *parser.NTFSContext
}

// Open opens a file for reading.
func (fsys *FS) Open(name string) (item fs.File, err error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}
	name = "/" + name

	dir, err := fsys.ntfsCtx.GetMFT(5)
	if err != nil {
		return nil, err
	}
	entry, err := dir.Open(fsys.ntfsCtx, name)

	return &Item{entry: entry, name: path.Base(name), path: name, ntfsCtx: fsys.ntfsCtx}, err
}
