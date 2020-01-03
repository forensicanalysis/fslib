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
// Paths
//
// Paths in fslib use the following conventions:
//     - all paths are separated by forward slashes '/' (yes, even the windows registry)
//     - all paths need to start with forward slashes '/' (exception: the OSFS accepts relative paths)
package fslib

import (
	"io"
	"os"
)

// FS is an interface for read-only file systems.
type FS interface {
	// Name returns the name of the file system.
	Name() string

	// Open opens a file for reading.
	Open(path string) (Item, error)

	// Stat return an os.FileInfo object that describes a file.
	Stat(path string) (os.FileInfo, error)
}

// Item is an interface for elements (e.g. files and directories) in read-only file
// systems.
type Item interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker

	// Name returns the name of the file.
	Name() string

	// Stat returns an os.FileInfo object that describes a file.
	Stat() (os.FileInfo, error)

	// Readdirnames returns up to n child items of a directory.
	Readdirnames(n int) ([]string, error)
}
