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

// +build !windows !go1.8

package registryfs

import (
	"errors"
	"io/fs"
	"os"
)

// New creates a new dummy registry FS.
func New() *FS { return &FS{} }

// FS implements a dummy file system for Windows Registries.
type FS struct{}

// Name returns the name of the file system.
func (*FS) Name() (name string) { return "Registry FS" }

// Open fails for non Windows operating systems.
func (m *FS) Open(name string) (item fs.File, err error) {
	return nil, errors.New("registry only supported on Windows")
}

// Stat fails for non Windows operating systems.
func (m *FS) Stat(name string) (os.FileInfo, error) {
	return nil, errors.New("registry only supported on Windows")
}
