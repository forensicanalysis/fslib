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

// Package testfs provides a in memory forensicfs implementation for testing.
package testfs

import (
	"bytes"
	"io/fs"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/forensicanalysis/fslib/filesystem"
	"github.com/forensicanalysis/fslib/forensicfs"
)

// FS implements a read-only memory file system for testing.
type FS struct {
	items map[string]fs.File
}

// Name returns the name of the file system.
func (*FS) Name() string { return "FS" }

// Open opens a file for reading.
func (fsys *FS) Open(name string) (fs.File, error) {
	name, err := filesystem.Clean(name)
	if err != nil {
		return nil, err
	}

	name = strings.Trim(name, "/")
	if fsys.items == nil {
		fsys.items = map[string]fs.File{"": &Directory{fs: fsys, path: ""}}
	}
	if item, ok := fsys.items[name]; ok {
		return item, nil
	}
	return nil, os.ErrNotExist
}

// Stat returns an os.FileInfo object that describes a file.
func (fsys *FS) Stat(name string) (os.FileInfo, error) {
	name, err := filesystem.Clean(name)
	if err != nil {
		return nil, err
	}

	name = strings.Trim(name, "/")
	if fsys.items == nil {
		fsys.items = map[string]fs.File{"": &Directory{fs: fsys, path: ""}}
	}
	if item, ok := fsys.items[name]; ok {
		return item.Stat()
	}
	return nil, os.ErrNotExist
}

// CreateDir adds a directory and all required parent directories to the file
// system.
func (fsys *FS) CreateDir(name string) {
	name = strings.Trim(name, "/")
	if fsys.items == nil {
		fsys.items = map[string]fs.File{"": &Directory{fs: fsys, path: ""}}
	}
	parts := strings.Split(name, "/")
	for i := range parts {
		name = strings.Join(parts[:i+1], "/")
		fsys.items[name] = &Directory{fs: fsys, path: name}
	}
}

// CreateFile adds a file and all required parent directories to the file system.
func (fsys *FS) CreateFile(name string, data []byte) {
	name = strings.TrimLeft(name, "/")
	if fsys.items == nil {
		fsys.items = map[string]fs.File{"": &Directory{fs: fsys, path: ""}}
	}
	fsys.items[name] = &File{name: path.Base(name), data: bytes.NewReader(data)}
}

// File describes a single file in the test file system.
type File struct {
	forensicfs.FileDefaults
	forensicfs.FileInfoDefaults
	name string
	data *bytes.Reader
}

// Name returns the name of the file.
func (f *File) Name() (name string) { return f.name }

// Read reads bytes into the passed buffer.
func (f *File) Read(dst []byte) (int, error) { return f.data.Read(dst) }

// ReadAt reads bytes starting at off into passed buffer.
func (f *File) ReadAt(b []byte, off int64) (int, error) { return f.data.ReadAt(b, off) }

// Seek move the current offset to the given position.
func (f *File) Seek(offset int64, whence int) (pos int64, err error) {
	return f.data.Seek(offset, whence)
}

// Size returns the file size.
func (f *File) Size() (n int64) { return f.data.Size() }

// Close does nothing for test file systems.
func (*File) Close() error { return nil }

// Stat returns the file itself as os:FileInfo.
func (f *File) Stat() (os.FileInfo, error) { return f, nil }

// Mode returns the os.FileMode.
func (*File) Mode() os.FileMode { return 0 }

// ModTime returns the current time for the file.
func (*File) ModTime() time.Time { return time.Now() }

// Sys returns nil for test file systems.
func (*File) Sys() interface{} { return nil }

// Directory describes a single directory in the test file system.
type Directory struct {
	forensicfs.DirectoryDefaults
	forensicfs.DirectoryInfoDefaults
	fs   *FS
	path string
}

// Name returns the name of the directory.
func (d *Directory) Name() (name string) { return path.Base(name) }

// Readdirnames returns up to n child items of a directory.
func (d *Directory) Readdirnames(n int) (items []string, err error) {
	seen := map[string]bool{}
	for itemPath := range d.fs.items {
		folder := ""
		if d.path != "" {
			folder = d.path + "/"
		}
		if strings.HasPrefix(itemPath, folder) {
			parts := strings.Split(itemPath[len(folder):], "/")
			if _, ok := seen[parts[0]]; !ok && parts[0] != "" {
				seen[parts[0]] = true
				items = append(items, parts[0])
			}
		}
	}
	sort.Strings(items)
	if n > 0 && len(items) > n {
		return items[:n], err
	}
	return
}

// Close does nothing for test file systems.
func (*Directory) Close() error { return nil }

// Stat returns the directory itself as os:FileInfo.
func (d *Directory) Stat() (os.FileInfo, error) { return d, nil }

// ModTime returns the current time for the directory.
func (d *Directory) ModTime() time.Time { return time.Now() }

// Sys returns nil for test file systems.
func (d *Directory) Sys() interface{} { return nil }
