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

// Package fallbackfs implements a meta filesystem that wraps a sequence of file
// systems. If opening a file on the first file system fails, it is tried on the
// next file systems of the filesystem. This can be useful when the first file
// system has access restrictions that can be circumvented this way. A live Windows
// file system can be backed with a raw disk file system for example, to enable
// extraction of locked files.
package fallbackfs

import (
	"fmt"
	"io"
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
	if !fs.ValidPath(name) {
		return nil, fmt.Errorf("path %s invalid", name)
	}

	for _, fallbackFilesystem := range fsys.fallbackFilesystems {
		item, err = fallbackFilesystem.Open(name)
		if err == nil {
			return &Item{name, item, fsys.fallbackFilesystems[1:]}, nil
		}
	}

	return
}

// Stat returns an fs.FileInfo object that describes a file.
func (fsys *FS) Stat(name string) (info fs.FileInfo, err error) {
	if !fs.ValidPath(name) {
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

type Item struct {
	path                string
	first               fs.File
	fallbackFilesystems []fs.FS
}

func (i *Item) Stat() (fs.FileInfo, error) {
	info, err := i.first.Stat()
	if err != nil {
		for _, fsys := range i.fallbackFilesystems {
			info, err = fs.Stat(fsys, i.path)
			if err == nil {
				break
			}
		}
	}
	return info, nil
}

func (i *Item) Read(bytes []byte) (int, error) {
	buf := make([]byte, len(bytes))
	n, err := i.first.Read(buf)
	if err != nil && err != io.EOF {
		for _, fsys := range i.fallbackFilesystems {
			file, err := fsys.Open(i.path)
			if err != nil {
				continue
			}
			buf = make([]byte, len(bytes))
			n, err = file.Read(buf)
			if err == nil || err == io.EOF {
				break
			}
		}
	}
	return copy(bytes, buf[:n]), err
}

func (i *Item) Close() error {
	return i.first.Close()
}

func (i *Item) ReadDir(n int) ([]fs.DirEntry, error) {
	if directory, ok := i.first.(fs.ReadDirFile); ok {
		return directory.ReadDir(n)
	}
	return nil, fmt.Errorf("%v does not implement ReadDir", i)
}
