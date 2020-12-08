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
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/windows/registry"

	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem"
	"github.com/forensicanalysis/fslib/forensicfs"
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
func (fs *FS) Open(name string) (item fslib.Item, err error) {
	name, err = filesystem.Clean(name)
	if err != nil {
		return nil, err
	}
	name = name[1:]

	if name == "" {
		return &Root{fs: fs}, nil
	}

	parts := strings.Split(name, "/")
	root := registryRoots[parts[0]]
	for i := range parts {
		parts[i] = strings.Replace(parts[i], `\`, `/`, -1)
	}

	k, err := registry.OpenKey(root, filepath.Join(parts[1:]...), registry.READ|registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
	return &Key{Key: &k, name: path.Base(name), path: name, fs: fs}, err
}

// Stat returns an os.FileInfo object that describes a file.
func (fs *FS) Stat(name string) (os.FileInfo, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	return f.Stat()
}

// Root is a pseudo root for the Windows registry.
type Root struct {
	forensicfs.DirectoryDefaults
	fs *FS
}

// Name always returns / for registry pseudo roots.
func (r Root) Name() string { return "/" }

// Readdirnames lists all registry roots in the registry.
func (r Root) Readdirnames(int) (items []string, err error) {
	for name := range registryRoots {
		_, err := r.fs.Open("/" + name)
		if err == nil {
			items = append(items, name)
		}
	}
	sort.Strings(items)
	return items, nil
}

// Size returns 0 for registry pseudo roots.
func (r *Root) Size() int64 { return 0 }

// Mode returns os.ModeDir for registry pseudo roots.
func (r *Root) Mode() os.FileMode { return os.ModeDir }

// ModTime returns the zero time (0001-01-01 00:00) for registry pseudo roots.
func (r *Root) ModTime() time.Time { return time.Time{} }

// IsDir returns true for registry pseudo roots.
func (r *Root) IsDir() bool { return true }

// Sys returns nil for registry pseudo roots.
func (r *Root) Sys() interface{} { return nil }

// Close does not do anything for registry pseudo roots.
func (r *Root) Close() error { return nil }

// Stat returns the registry pseudo roots itself as os.FileMode.
func (r *Root) Stat() (os.FileInfo, error) { return r, nil }

// Key is an entry in the registry.
type Key struct {
	forensicfs.DirectoryDefaults
	Key  *registry.Key
	name string
	path string
	fs   *FS
}

// Name returns the name of the file.
func (rk *Key) Name() string {
	return rk.name
}

// Readdirnames returns up to n sub keys of a key.
func (rk *Key) Readdirnames(n int) (items []string, err error) {
	items = []string{}
	subKeyNames, err := rk.Key.ReadSubKeyNames(n)
	if err != nil && err != io.EOF {
		return items, fmt.Errorf("error ReadSubKeyNames: %w", err)
	}
	for _, subKeyName := range subKeyNames {
		items = append(items, strings.Replace(subKeyName, `/`, `\`, -1))
	}
	sort.Strings(items)
	return items, nil
}

// Close closes the key freeing the resource. Usually additional IO operations fail
// after closing.
func (rk *Key) Close() error { return rk.Key.Close() }

// Stat return an os.FileInfo object that describes a key.
func (rk *Key) Stat() (os.FileInfo, error) {
	info, err := rk.Key.Stat()
	return &KeyInfo{info, rk.name}, err
}

// KeyInfo describes a key.
type KeyInfo struct {
	*registry.KeyInfo
	name string
}

// Name returns the name of the key.
func (rk *KeyInfo) Name() string { return rk.name }

// Size returns the file size.
func (rk *KeyInfo) Size() int64 { return 0 }

// IsDir returns if the key has subkeys.
func (rk *KeyInfo) IsDir() bool { return rk.SubKeyCount > 0 }

// ModTime returns the modification time.
func (rk *KeyInfo) ModTime() time.Time { return rk.KeyInfo.ModTime().In(time.UTC) }

// Sys returns underlying data source.
func (rk *KeyInfo) Sys() interface{} { return rk.KeyInfo }

// Mode returns the os.FileMode.
func (rk *KeyInfo) Mode() os.FileMode {
	if rk.IsDir() {
		return os.ModeDir
	}
	return 0
}
