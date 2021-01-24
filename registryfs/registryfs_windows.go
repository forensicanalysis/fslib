// +build go1.8

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

package registryfs

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
)

var registryRoots = map[string]registry.Key{
	"HKEY_CLASSES_ROOT":     registry.CLASSES_ROOT,
	"HKEY_CURRENT_USER":     registry.CURRENT_USER,
	"HKEY_LOCAL_MACHINE":    registry.LOCAL_MACHINE,
	"HKEY_USERS":            registry.USERS,
	"HKEY_CURRENT_CONFIG":   registry.CURRENT_CONFIG,
	"HKEY_PERFORMANCE_DATA": registry.PERFORMANCE_DATA,
}

// New creates a new registry FS.
func New() *FS { return &FS{} }

// FS implements a read-only file system for Windows Registries.
type FS struct{}

// Name returns the name of the file system.
func (*FS) Name() (name string) { return "Registry FS" }

// Open opens a file for reading.
func (fsys *FS) Open(name string) (item fs.File, err error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}

	if name == "." {
		return &Root{fs: fsys}, nil
	}

	parts := strings.Split(name, "/")
	root := registryRoots[parts[0]]
	for i := range parts {
		parts[i] = strings.Replace(parts[i], `\`, `/`, -1)
	}

	p := filepath.Join(parts[1:]...)
	k, err := registry.OpenKey(root, p, registry.READ|registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	return &Key{Key: &k, name: path.Base(name), path: name, fs: fsys}, err
}
