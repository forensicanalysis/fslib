package registryfs

import (
	"io/fs"
	"os"
	"time"
)

// Root is a pseudo root for the Windows registry.
type Root struct {
	fs *FS
}

// Name always returns . for registry pseudo roots.
func (r Root) Name() string { return "." }

// ReadDir lists all registry roots in the registry.
func (r Root) ReadDir(int) (entries []fs.DirEntry, err error) {
	for name := range registryRoots {
		info, err := fs.Stat(r.fs, name)
		if err == nil {
			entries = append(entries, info.(*KeyInfo))
		}
	}
	return entries, nil
}

func (r *Root) Read([]byte) (int, error) { return 0, nil }

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
func (r *Root) Stat() (fs.FileInfo, error) { return r, nil }
