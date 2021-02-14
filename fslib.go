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

// Package fslib project contains a collection of packages to parse file
// systems, archives and similar data. The included packages can be used to
// access disk images of with different partitioning and file systems.
// Additionally, file systems for live access to the currently mounted file system
// and registry (on Windows) are implemented.
package fslib

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"runtime"
	"sort"
)

const windows = "windows"

func ReadDir(file fs.File, n int) (items []fs.DirEntry, err error) {
	if directory, ok := file.(fs.ReadDirFile); ok {
		return directory.ReadDir(n)
	}
	return nil, fmt.Errorf("%v does not implement ReadDir", file)
}

// ToFSPath converts a normal path (e.g. 'C:\Windows') to a fs path
// ('C/Windows').
func ToFSPath(systemPath string) (name string, err error) {
	name, err = filepath.Abs(systemPath)
	if err != nil {
		return "", err
	}
	if runtime.GOOS == windows {
		name = filepath.ToSlash(name)
		name = name[:1] + name[2:]
		return name, nil
	}
	return name[1:], nil
}

func DirEntries(n int, items []fs.DirEntry, dirOffset int) ([]fs.DirEntry, int, error) {
	sort.Sort(ByName(items))

	// directory already exhausted
	if n <= 0 && dirOffset >= len(items) {
		return nil, 0, nil
	}

	var err error
	// read till end
	if n > 0 && dirOffset+n > len(items) {
		err = io.EOF
		if dirOffset > len(items) {
			return nil, 0, err
		}
	}

	offset := 0
	if n > 0 && dirOffset+n <= len(items) {
		items = items[dirOffset : dirOffset+n]
		offset += n
	} else {
		items = items[dirOffset:]
		offset += len(items)
	}

	return items, offset, err
}

type ByName []fs.DirEntry

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }
