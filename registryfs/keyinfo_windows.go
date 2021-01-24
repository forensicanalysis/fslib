package registryfs

import (
	"io/fs"
	"time"

	"golang.org/x/sys/windows/registry"
)

// KeyInfo describes a key.
type KeyInfo struct {
	*registry.KeyInfo
	name string
}

func (rk *KeyInfo) Type() fs.FileMode { return rk.Mode() }

func (rk *KeyInfo) Info() (fs.FileInfo, error) { return rk, nil }

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

// Mode returns the fs.FileMode.
func (rk *KeyInfo) Mode() fs.FileMode {
	if rk.IsDir() {
		return fs.ModeDir
	}
	return 0
}
