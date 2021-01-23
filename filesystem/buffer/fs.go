package buffer

import (
	"fmt"
	"io/fs"
	"path"
)

// New creates a new buffer FS.
func New(fsys fs.FS) *FS {
	return &FS{internal: fsys}
}

// FS implements a read-only meta file system where failing method calls to
// higher level file systems are passed to other file systems.
type FS struct {
	internal fs.FS
}

// Open opens a file for reading.
func (fsys *FS) Open(name string) (item fs.File, err error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}

	f, err := fsys.internal.Open(name)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	return &File{
		name:   path.Base(name),
		fs:     fsys,
		file:   f,
		size:   info.Size(),
		offset: 0,
		isdir:  info.IsDir(),
		closed: false,
		buf:    nil,
	}, nil
}
