package osfs

import (
	"github.com/forensicanalysis/fslib/forensicfs"
	"os"
	"time"
)

// Root is a pseudo root directory for windows partitions.
type Root struct{ forensicfs.DirectoryDefaults }

// Name always returns / for window pseudo roots.
func (*Root) Name() (name string) { return "/" }

// Close does not do anything for window pseudo roots.
func (*Root) Close() error { return nil }

// Size returns 0 for window pseudo roots.
func (*Root) Size() int64 { return 0 }

// Mode returns os.ModeDir for window pseudo roots.
func (*Root) Mode() os.FileMode { return os.ModeDir }

// ModTime returns the zero time (0001-01-01 00:00) for window pseudo roots.
func (*Root) ModTime() time.Time { return time.Time{} }

// IsDir returns true for window pseudo roots.
func (*Root) IsDir() bool { return true }

// Sys returns nil for window pseudo roots.
func (*Root) Sys() interface{} { return nil }

// Stat returns the windows pseudo roots itself as os.FileMode.
func (root *Root) Stat() (os.FileInfo, error) {
	return root, nil
}
