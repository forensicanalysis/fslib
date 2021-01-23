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

package zip

import (
	"github.com/forensicanalysis/fslib"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// File describes files and directories in the zip file system.
type File struct {
	internal afero.File
}

// Close closes the file freeing the resource. Other IO operations fail after
// closing.
func (f *File) Close() (err error) {
	return f.internal.Close()
}

// Read reads bytes into the passed buffer.
func (f *File) Read(p []byte) (n int, err error) {
	return f.internal.Read(p)
}

// ReadAt reads bytes starting at off into passed buffer.
func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	return f.internal.ReadAt(p, off)
}

// Seek move the current offset to the given position.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.internal.Seek(offset, whence)
}

// Name returns the name of the file.
func (f *File) Name() string {
	return filepath.ToSlash(filepath.Base(f.internal.Name()))
}

// Readdirnames returns up to n child items of a directory.
func (f *File) ReadDir(count int) ([]fs.DirEntry, error) {
	if count == 0 {
		count = -1
	}

	infos, err := f.internal.Readdir(count)
	if err != nil {
		return nil, err
	}
	entries := fslib.InfosToEntries(infos)
	return entries, nil
}

// Readdirnames returns up to n child items of a directory.
func (f *File) Readdirnames(count int) (names []string, err error) {
	if count == 0 {
		count = -1
	}
	return f.internal.Readdirnames(count)
}

// Stat return an os.FileInfo object that describes a file.
func (f *File) Stat() (os.FileInfo, error) {
	/*if f.Name() == "/" {
		return &RootInfo{}, nil
	}*/
	return f.internal.Stat()
}

// Sys returns underlying data source.
func (f *File) Sys() interface{} {
	attr := map[string]interface{}{
		// "modified": f.zipfile.Modified.In(time.UTC),
	}

	mode, err := f.Stat()
	if err == nil {
		attr["mode"] = mode
	}

	return attr
}
