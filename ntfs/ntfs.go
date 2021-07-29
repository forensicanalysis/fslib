// Copyright (c) 2019-2020 Siemens AG
// Copyright (c) 2019-2021 Jonas Plum
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

// Package ntfs provides an io/fs implementation of the New Technology File
// System (NTFS).
//
// Currently alternate data streams are not supported.
package ntfs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"

	"www.velocidex.com/golang/go-ntfs/parser"
)

const (
	defaultPageSize  = 1024 * 1024
	defaultCacheSize = 100 * 1024 * 1024
)

func checkPageSizeAndCacheSize(pageSize int64, cacheSize int) (int64, int) {
	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	if cacheSize <= 0 {
		cacheSize = defaultCacheSize
	}
	return pageSize, cacheSize
}

// New creates a new ntfs FS.
func New(r io.ReaderAt) (fs *FS, err error) {
	return NewWithSize(r, defaultPageSize, defaultCacheSize)
}

// NewWithSize creates a new ntfs FS with specific pageSize and cacheSize.
func NewWithSize(r io.ReaderAt, pageSize int64, cacheSize int) (fs *FS, err error) {
	pageSize, cacheSize = checkPageSizeAndCacheSize(pageSize, cacheSize)
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error parsing file system as NTFS")
		}
	}()
	reader, err := parser.NewPagedReader(r, pageSize, cacheSize)
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
	if !valid || strings.Contains(name, `\`) {
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
