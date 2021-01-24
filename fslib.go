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

// Package fslib contains a collection of tools and libraries to parse file
// systems, archives and other data types. The included libraries can be used to
// access disk images of with different partitioning and file systems.
// Additionally, file systems for live access to the currently mounted file system
// and registry (on Windows) are implemented.
//
// File systems supported
//
// - Native OS file system (directory listing for Windows root provides list of drives)
// - ZIP
// - NTFS
// - FAT16
// - MBR
// - GPT
// - Windows Registry (live not from files)
//
// Meta file systems
//
// - ⭐ **Recursive FS**: Access container files on file systems recursively, e.g. `"ntfs.dd/forensic.zip/Computer forensics - Wikipedia.pdf"`
// - Buffer FS: Buffer accessed files of an underlying file system
// - System FS: Similar to the native OS file system, but falls back to NTFS on failing access on Windows
//
// Paths
//
// Paths in fslib use [io/fs](https://tip.golang.org/pkg/io/fs/#ValidPath) conventions:
//
// > - Path names passed to open are unrooted, slash-separated sequences of path elements, like “x/y/z”.
// > - Path names must not contain a “.” or “..” or empty element, except for the special case that the root directory is named “.”.
// > - Paths are slash-separated on all systems, even Windows.
// > - Backslashes must not appear in path names.
package fslib

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"runtime"
	"strings"
)

const windows = "windows"

func ReadDir(file fs.File, n int) (items []fs.DirEntry, err error) {
	if directory, ok := file.(fs.ReadDirFile); ok {
		return directory.ReadDir(n)
	}
	return nil, fmt.Errorf("%v does not implement ReadDir", file)
}

// ToForensicPath converts a normal path (e.g. 'C:\Windows') to a fs path
// ('C/Windows').
func ToForensicPath(systemPath string) (name string, err error) {
	name, err = filepath.Abs(systemPath)
	if err != nil {
		return "", err
	}
	if runtime.GOOS == windows {
		name = strings.Replace(name, "\\", "/", -1)
		name = name[:1] + name[2:]
		return name, nil
	}
	return name[1:], nil
}
