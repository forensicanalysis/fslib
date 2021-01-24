package fat16

import (
	"io/fs"
	"os"
	"strings"
	"time"
)

type directoryEntry struct {
	Filename          [8]byte
	FilenameExtension [3]byte
	FileAttributes    byte
	_                 [10]byte // Reserved
	Timecreated       [2]byte
	Datecreated       [2]byte
	Startingcluster   uint16
	FileSize          uint32
}

func formatFilename(de *directoryEntry) string {
	filename := strings.TrimSpace(string(de.Filename[:]))
	if de.FilenameExtension[0] != 0x20 {
		filename = filename + "." + strings.TrimSpace(string(de.FilenameExtension[:]))
	}
	return filename
}

type namedEntry struct {
	directoryEntry
	name string
}

func (d *namedEntry) Name() string {
	return d.name
}

func (d *namedEntry) IsDir() bool {
	return d.FileAttributes&0x10 != 0
}

func (d *namedEntry) Size() int64 {
	return int64(d.FileSize)
}

func (d *namedEntry) Mode() fs.FileMode {
	if d.FileAttributes&0x10 != 0 {
		return os.ModeDir
	}
	return 0
}

func (d *namedEntry) ModTime() time.Time {
	return time.Time{} // TODO parse d.Timecreated
}

func (d *namedEntry) Type() fs.FileMode {
	return d.Mode()
}

func (d *namedEntry) Info() (fs.FileInfo, error) {
	return d, nil
}

func (d *namedEntry) Sys() interface{} {
	return d
}
