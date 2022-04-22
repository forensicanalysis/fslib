//go:build go1.8
// +build go1.8

package registryfs

import (
	"fmt"
	"io"
	"io/fs"

	"golang.org/x/sys/windows/registry"

	"github.com/forensicanalysis/fslib"
)

// Key is an entry in the registry.
type Key struct {
	Key       *registry.Key
	name      string
	path      string
	fs        *FS
	dirOffset int
}

func (rk *Key) Read([]byte) (int, error) {
	return 0, io.EOF
}

// Name returns the name of the file.
func (rk *Key) Name() string {
	return rk.name
}

// ReadDir returns up to n sub keys of a key.
func (rk *Key) ReadDir(n int) (entries []fs.DirEntry, err error) {
	var items []fs.DirEntry
	subKeyNames, err := rk.Key.ReadSubKeyNames(-1)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error ReadSubKeyNames: %w", err)
	}
	for _, subKeyName := range subKeyNames {
		info, err := subKeyInfo(rk, subKeyName)
		if err != nil {
			items = append(items, &KeyInfo{name: subKeyName, KeyInfo: &registry.KeyInfo{}})
			continue
		}

		items = append(items, &KeyInfo{name: subKeyName, KeyInfo: info})
	}

	items, offset, err := fslib.DirEntries(n, items, rk.dirOffset)
	rk.dirOffset += offset
	return items, err
}

func subKeyInfo(rk *Key, subKeyName string) (*registry.KeyInfo, error) {
	subKey, err := registry.OpenKey(*rk.Key, subKeyName, registry.READ)
	if err != nil {
		return nil, err
	}
	defer subKey.Close()
	info, err := subKey.Stat()
	if err != nil {
		return nil, err
	}
	return info, nil
}

// Close closes the key freeing the resource. Usually additional IO operations fail
// after closing.
func (rk *Key) Close() error { return rk.Key.Close() }

// Stat return a fs.FileInfo object that describes a key.
func (rk *Key) Stat() (fs.FileInfo, error) {
	info, err := rk.Key.Stat()
	return &KeyInfo{info, rk.name}, err
}
