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

// Package forensicfs provides defaults for read-only file system items.
package forensicfs

import (
	"os"
	"syscall"
)

// FileDefaults implements default methods for files.
type FileDefaults struct{}

// Readdirnames returns an error for files.
func (*FileDefaults) Readdirnames(count int) ([]string, error) { return nil, syscall.EPERM }

// FileInfoDefaults implements default methods for file infos.
type FileInfoDefaults struct{}

// IsDir returns false for file infos.
func (*FileInfoDefaults) IsDir() bool { return false }

// DirectoryDefaults implements default methods for directory infos.
type DirectoryDefaults struct{}

// Read returns an error for directories.
func (*DirectoryDefaults) Read(p []byte) (n int, err error) { return 0, syscall.EPERM }

// ReadAt returns an error for directories.
func (*DirectoryDefaults) ReadAt(p []byte, off int64) (n int, err error) { return 0, syscall.EPERM }

// Seek returns an error for directories.
func (*DirectoryDefaults) Seek(offset int64, whence int) (int64, error) { return 0, syscall.EPERM }

// DirectoryInfoDefaults implements default methods for directory infos.
type DirectoryInfoDefaults struct{}

// IsDir returns true for directory infos.
func (*DirectoryInfoDefaults) IsDir() bool { return true }

// Mode returns os.ModeDir for directory infos.
func (*DirectoryInfoDefaults) Mode() os.FileMode { return os.ModeDir }

// Size returns 0 for directory infos.
func (*DirectoryInfoDefaults) Size() int64 { return 0 }
