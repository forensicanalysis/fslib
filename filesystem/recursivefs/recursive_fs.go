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

// Package recursivefs provides a forensicfs implementation that can open paths in
// nested forensicfs recursively. The forensicfs are identified using the filetype
// library. This way e.g. a file in a zip inside a disk image can be accessed.
package recursivefs

import (
	"errors"
	"fmt"
	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem/osfs"
	"github.com/forensicanalysis/fslib/fsio"
	"io/fs"
	"os"
	"sort"
)

type element struct {
	Parser string
	Key    string
}

// FS implements a read-only meta file system that can access nested file system
// structures.
type FS struct{}

// New creates a new recursive FS.
func New() *FS { return &FS{} }

// Name returns the filesystem name.
func (fsys *FS) Name() string { return "RecFS" }

// Open returns a File for the given location.
func (fsys *FS) Open(name string) (f fs.File, err error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}

	elems, err := parseRealPath(name)
	if err != nil {
		return
	}

	var childFS fs.FS = osfs.New()
	var childName = ""
	for _, elem := range elems {
		if f != nil {
			childFS, err = fsFromName(elem.Parser, f)
			if err != nil {
				return nil, err
			}

			fi, err := f.Stat()
			if err == nil && fi.IsDir() {
				f.Close() // nolint: errcheck
			}
		}

		f, err = childFS.Open(elem.Key)
		if err != nil {
			return nil, err
		}

		childName = elem.Key
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return &Item{File: f, path: name, recursiveFS: fsys, isFS: false}, nil
	}

	isFS, ifs, err := detectFsFromFile(name, f)
	if err != nil {
		return nil, err
	}

	f, err = reopen(f, childFS, childName)
	if err != nil {
		return nil, err
	}

	return &Item{File: f, path: name, innerFSName: ifs, recursiveFS: fsys, isFS: isFS}, nil
}

// Stat returns an os.FileInfo object that describes a file.
func (fsys *FS) Stat(name string) (os.FileInfo, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	return f.Stat()
}

// Item describes files and directories in the file system.
type Item struct {
	fs.File
	path        string
	innerFSName string
	recursiveFS *FS
	isFS        bool
}

// ReadDir returns up to n child items of a directory.
func (i *Item) ReadDir(n int) (items []fs.DirEntry, err error) {
	if !i.isFS {
		items, err = fslib.ReadDir(i.File, n)
	} else {
		if readSeekerAtItem, ok := i.File.(fsio.ReadSeekerAt); ok {
			innerFS, _ := fsFromName(i.innerFSName, readSeekerAtItem)
			root, _ := innerFS.Open(".")
			items, err = fslib.ReadDir(root, n)
		} else {
			return nil, errors.New("not a file does not implement Seek and ReadAt")
		}
	}
	if err != nil {
		return nil, fmt.Errorf("could not Readdirnames %#v: %w", i, err)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name() < items[j].Name()
	})

	return items, nil
}

// Stat return an os.FileInfo object that describes a file.
func (i *Item) Stat() (os.FileInfo, error) {
	info, err := i.File.Stat()
	return &Info{info, i.isFS}, err
}

// Info wraps the os.FileInfo.
type Info struct {
	os.FileInfo
	isFS bool
}

// IsDir returns if the item is a directory. Returns true for files that are file
// systems (e.g. zip archives).
func (m *Info) IsDir() bool {
	if m.isFS {
		return true
	}
	return m.FileInfo.IsDir()
}
