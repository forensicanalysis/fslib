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

// Package osfs provides a forensicfs implementation of the native OS file system.
// In windows paths are changed from "C:\Windows" to "/C/Windows" to comply with
// the path specifications of the fslib.
package osfs

import (
	"errors"
	"fmt"
	"github.com/forensicanalysis/fslib"
	"io/fs"
	"os"
	"runtime"
	"sort"
)

const windows = "windows"

// New wrapes the nativ file system.
func New() *FS {
	return &FS{}
}

// FS implements a read-only wrapper for the native file system.
type FS struct{}

// Name returns the name of the file system.
func (fs *FS) Name() string {
	return "OsFs"
}

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

// OpenSystemPath opens a normal path (e.g. 'C:\Windows') instead of a fslib path
// ('/C/Windows').
func (fs *FS) OpenSystemPath(syspath string) (item fs.File, err error) {
	syspath, err = fslib.ToForensicPath(syspath)
	if err != nil {
		return nil, err
	}
	return fs.Open(syspath)
}

// Open opens a file for reading.
func (fs *FS) Open(name string) (item fs.File, err error) {
	name, sysname, err := sysname(name)
	if err != nil {
		return nil, err
	}

	if name == "/" && runtime.GOOS == windows {
		return &Root{}, nil
	}

	file, err := os.Open(sysname) // #nosec
	if err != nil {
		return nil, err
	}

	return &Item{File: *file, syspath: sysname}, err
}

// Stat returns an os.FileInfo object that describes a file.
func (fs *FS) Stat(name string) (os.FileInfo, error) {
	name, sysname, err := sysname(name)
	if err != nil {
		return nil, err
	}

	if name == "/" && runtime.GOOS == windows {
		return &Root{}, nil
	}

	fi, err := os.Lstat(sysname)
	if err != nil {
		return nil, err
	}

	return &Info{fi, sysname}, nil
}

func sysname(name string) (string, string, error) {
	if runtime.GOOS == windows && len(name) > 1 && !isLetter(name[1]) {
		return "", "", errors.New("partition must be a letter")
	}
	if runtime.GOOS == windows && len(name) > 2 && name[2] != '/' {
		return "", "", errors.New("partition must be followed by a slash")
	}
	valid := fs.ValidPath(name)
	if !valid {
		return "", "", fmt.Errorf("path %s invalid", name)
	}
	if name == "/" {
		return "/", "/", nil
	}
	sysname := name
	if runtime.GOOS == windows {
		sysname = string(name[1]) + ":"
		if len(name) > 2 {
			sysname += name[2:]
		} else {
			sysname += "/"
		}
	}
	return name, sysname, nil
}

// Item describes files and directories in the native OS file system.
type Item struct {
	os.File
	syspath string
}

// Name returns the name of the item.
func (i *Item) Name() string {
	return i.File.Name()
}

// Readdirnames returns up to n child items of a directory.
func (i *Item) Readdirnames(n int) (items []string, err error) {
	items, err = i.File.Readdirnames(n)
	if items == nil {
		items = []string{}
	}
	sort.Strings(items)
	return items, err
}

// Close closes the file freeing the resource. Usually additional IO operations
// fail after closing.
func (i *Item) Close() error {
	return i.File.Close()
}

// Stat return an os.FileInfo object that describes a file.
func (i *Item) Stat() (os.FileInfo, error) {
	info, err := os.Lstat(i.syspath)
	return &Info{info, i.syspath}, err
}

// Info wraps os.FileInfo for native OS items.
type Info struct {
	os.FileInfo
	syspath string
}
