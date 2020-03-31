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
	"os"
	"path"
	"syscall"

	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem"
	"github.com/forensicanalysis/fslib/fsio"
)

// FS implements a read-only file system for zip files.
type FS struct {
	r     *zip.Reader
	files map[string]map[string]*zip.File
}

func splitpath(name string) (dir, file string) {
	if len(name) == 0 || name[0] != '/' {
		name = "/" + name
	}
	name = path.Clean(name)
	dir, file = path.Split(name)
	dir = path.Clean(dir)
	return
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

	fs := &FS{r: zr, files: make(map[string]map[string]*zip.File)}
	for _, file := range zr.File {
		d, f := splitpath(file.Name)
		if _, ok := fs.files[d]; !ok {
			fs.files[d] = make(map[string]*zip.File)
		}
		if _, ok := fs.files[d][f]; !ok {
			fs.files[d][f] = file
		}
		if file.FileInfo().IsDir() {
			dirname := path.Join(d, f)
			if _, ok := fs.files[dirname]; !ok {
				fs.files[dirname] = make(map[string]*zip.File)
			}
		}
	}
	return fs, nil
}

// Name returns the name of the file system.
func (fs *FS) Name() string {
	return "ZIP"
}

// Open opens a file for reading.
func (fs *FS) Open(name string) (fslib.Item, error) {
	name, err := filesystem.Clean(name)
	if err != nil {
		return nil, err
	}

	d, f := splitpath(name)
	if f == "" {
		return &File{fs: fs, isdir: true}, nil // &pseudoRoot{}, nil //
	}
	if _, ok := fs.files[d]; !ok {
		return nil, &os.PathError{Op: "stat", Path: name, Err: syscall.ENOENT}
	}
	file, ok := fs.files[d][f]
	if !ok {
		return nil, &os.PathError{Op: "stat", Path: name, Err: syscall.ENOENT}
	}
	return &File{fs: fs, zipfile: file, isdir: file.FileInfo().IsDir()}, nil
}

// Stat returns an os.FileInfo object that describes a file.
func (fs *FS) Stat(name string) (os.FileInfo, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	return f.Stat()
}
