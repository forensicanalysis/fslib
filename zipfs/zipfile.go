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

package zipfs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"

	"github.com/forensicanalysis/fslib"
)

// File describes files and directories in the zip file system.
type File struct {
	path     string
	internal fs.File
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
	if readerAt, ok := f.internal.(io.ReaderAt); ok {
		return readerAt.ReadAt(p, off)
	}
	return 0, errors.New("does not implement ReadAt")
}

// Seek move the current offset to the given position.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	if readerAt, ok := f.internal.(io.Seeker); ok {
		return readerAt.Seek(offset, whence)
	}
	return 0, errors.New("does not implement Seek")
}

// Name returns the name of the file.
func (f *File) Name() string {
	return filepath.ToSlash(filepath.Base(f.path))
}

// Readdirnames returns up to n child items of a directory.
func (f *File) ReadDir(count int) ([]fs.DirEntry, error) {
	entries, err := fslib.ReadDir(f.internal, count)
	return uniqueEntries(entries), err
}

// Stat return an fs.FileInfo object that describes a file.
func (f *File) Stat() (fs.FileInfo, error) {
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

func uniqueEntries(entries []fs.DirEntry) []fs.DirEntry {
	keys := make(map[string]bool)
	var list []fs.DirEntry
	for _, entry := range entries {
		if fmt.Sprintf("%T", entry) != "zip.headerFileInfo" {
			continue
		}
		if _, value := keys[entry.Name()]; !value {
			keys[entry.Name()] = true
			list = append(list, entry)
		}
	}
	return list
}
