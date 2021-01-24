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

// Package osfs provides an io/fs implementation of the native OS file system.
// In windows paths are changed from "C:\Windows" to "/C/Windows" to comply with
// the path specifications of the fslib.
package osfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"runtime"

	"github.com/forensicanalysis/fslib"
)

const windows = "windows"

// New wraps the native file system.
func New() *FS {
	return &FS{}
}

// FS implements a read-only wrapper for the native file system.
type FS struct{}

func isLetter(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

// OpenSystemPath opens a normal path (e.g. 'C:\Windows') instead of a fslib path
// ('/C/Windows').
func (fsys *FS) OpenSystemPath(syspath string) (item fs.File, err error) {
	syspath, err = fslib.ToForensicPath(syspath)
	if err != nil {
		return nil, err
	}

	return fsys.Open(syspath)
}

// Open opens a file for reading.
func (fsys *FS) Open(name string) (item fs.File, err error) {
	name, sysname, err := sysname(name)
	if err != nil {
		return nil, err
	}

	if name == "." && runtime.GOOS == windows {
		return &Root{}, nil
	}

	file, err := os.Open(sysname) // #nosec
	if err != nil {
		return nil, err
	}

	return &Item{File: *file, syspath: sysname}, err
}

// Stat returns an os.FileInfo object that describes a file.
func (fsys *FS) Stat(name string) (os.FileInfo, error) {
	name, sysname, err := sysname(name)
	if err != nil {
		return nil, err
	}

	if name == "." && runtime.GOOS == windows {
		return &Root{}, nil
	}

	fi, err := os.Lstat(sysname)
	if err != nil {
		return nil, err
	}

	return &Info{fi, sysname}, nil
}

func sysname(name string) (string, string, error) {
	valid := fs.ValidPath(name)
	if !valid {
		return "", "", fmt.Errorf("path %s invalid", name)
	}
	if name == "." {
		return ".", ".", nil
	}
	if runtime.GOOS == windows && len(name) > 0 && !isLetter(name[0]) {
		return "", "", fmt.Errorf("partition must be a letter is %s", name)
	}
	if runtime.GOOS == windows && len(name) > 1 && name[1] != '/' {
		return "", "", errors.New("partition must be followed by a slash")
	}

	sysname := "/" + name
	if runtime.GOOS == windows {
		sysname = string(name[0]) + ":"
		if len(name) > 1 {
			sysname += name[1:]
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
