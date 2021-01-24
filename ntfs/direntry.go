package ntfs

import (
	"io/fs"
	"time"

	"www.velocidex.com/golang/go-ntfs/parser"
)

type DirEntry struct {
	info *parser.FileInfo
}

func (d *DirEntry) Name() string {
	return d.info.Name
}

func (d *DirEntry) IsDir() bool {
	return d.info.IsDir
}

func (d *DirEntry) Size() int64 {
	return d.info.Size
}

func (d *DirEntry) Mode() fs.FileMode {
	if d.IsDir() {
		return fs.ModeDir
	}
	return 0
}

func (d *DirEntry) ModTime() time.Time {
	return d.info.Mtime
}

func (d *DirEntry) Sys() interface{} {
	return d.info
}

func (d *DirEntry) Type() fs.FileMode {
	if d.IsDir() {
		return fs.ModeDir
	}
	return 0
}

func (d *DirEntry) Info() (fs.FileInfo, error) {
	return d, nil
}
