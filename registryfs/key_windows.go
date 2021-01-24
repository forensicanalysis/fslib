package registryfs

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io"
	"io/fs"
	"os"
)

// Key is an entry in the registry.
type Key struct {
	Key  *registry.Key
	name string
	path string
	fs   *FS
}

func (rk *Key) Read([]byte) (int, error) {
	return 0, nil
}

// Name returns the name of the file.
func (rk *Key) Name() string {
	return rk.name
}

// ReadDir returns up to n sub keys of a key.
func (rk *Key) ReadDir(n int) (entries []fs.DirEntry, err error) {
	var items []fs.DirEntry
	subKeyNames, err := rk.Key.ReadSubKeyNames(n)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("error ReadSubKeyNames: %w", err)
	}
	for _, subKeyName := range subKeyNames {
		subKey, err := registry.OpenKey(*rk.Key, subKeyName, registry.READ|registry.QUERY_VALUE|registry.ENUMERATE_SUB_KEYS)
		if err != nil {
			return nil, err
		}
		info, err := subKey.Stat()
		if err != nil {
			return nil, err
		}
		subKey.Close()

		items = append(items, &KeyInfo{KeyInfo: info, name: subKeyName})
	}
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
