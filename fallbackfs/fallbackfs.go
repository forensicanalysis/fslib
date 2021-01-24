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

// Package fallbackfs implements a meta filesystem that wraps a sequence of file
// systems. If opening a file on the first file system fails, it is tried on the
// next file systems of the filesystem. This can be useful when the first file
// system has access restrictions that can be circumvented this way. A live Windows
// file system can be backed with a raw disk file system for example, to enable
// extraction of locked files.
package fallbackfs

import (
	"fmt"
	"io/fs"
)

// New creates a new fallback FS.
func New(filesystems ...fs.FS) *FS {
	return &FS{fallbackFilesystems: filesystems}
}

// FS implements a read-only meta file system where failing method calls to
// higher level file systems are passed to other file systems.
type FS struct {
	fallbackFilesystems []fs.FS
}

// Open opens a file for reading.
func (fsys *FS) Open(name string) (item fs.File, err error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}

	for _, fallbackFilesystem := range fsys.fallbackFilesystems {
		item, err = fallbackFilesystem.Open(name)
		if err == nil {
			return
		}
	}

	return
}

// Stat returns an fs.FileInfo object that describes a file.
func (fsys *FS) Stat(name string) (info fs.FileInfo, err error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}

	for _, fallbackFilesystem := range fsys.fallbackFilesystems {
		info, err = fs.Stat(fallbackFilesystem, name)
		if err == nil {
			return
		}
	}

	return
}