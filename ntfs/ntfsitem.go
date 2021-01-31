package ntfs

import (
	"io"
	"io/fs"
	"os"
	"strings"

	"www.velocidex.com/golang/go-ntfs/parser"
)

// Item describes files and directories in the NTFS.
type Item struct {
	entry     *parser.MFT_ENTRY
	size      *int64
	attribute *parser.NTFS_ATTRIBUTE
	name      string
	offset    int64
	dirOffset int
	path      string
	ntfsCtx   *parser.NTFSContext
}

// Read reads bytes into the passed buffer.
func (i *Item) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	c, err := i.ReadAt(p, i.offset)
	i.offset += int64(c)
	return c, err
}

// ReadAt reads bytes starting at off into passed buffer.
func (i *Item) ReadAt(p []byte, off int64) (n int, err error) {
	if i.attribute == nil {
		attribute, err := i.entry.GetAttribute(i.ntfsCtx, 128, -1)
		if err != nil {
			return 0, err
		}
		i.attribute = attribute
	}

	n, err = i.attribute.Data(i.ntfsCtx).ReadAt(p, off)
	if int64(len(p)) > i.Size() {
		err = io.EOF
	}
	return
}

// Seek move the current offset to the given position.
func (i *Item) Seek(pos int64, whence int) (offset int64, err error) {
	if pos != 0 {
		switch whence {
		case os.SEEK_SET:
			i.offset = pos
		case os.SEEK_CUR:
			i.offset += pos
		case os.SEEK_END:
			i.offset = i.Size() + pos
		}
	}

	return i.offset, nil
}

// Size returns the item's size.
func (i *Item) Size() int64 {
	if i.size == nil {
		infos, err := parser.ModelMFTEntry(i.ntfsCtx, i.entry)
		if err != nil {
			return 0
		}
		i.size = &infos.Size
	}
	return *i.size
}

// ReadDir returns up to n child items of a directory.
func (i *Item) ReadDir(n int) (entries []fs.DirEntry, err error) {
	infos := parser.ListDir(i.ntfsCtx, i.entry)

	for _, info := range infos {
		if info.Name == "" || info.Name == "." || strings.Contains(info.Name, ":") {
			continue
		}
		entries = append(entries, &DirEntry{info})
	}

	// directory already exhausted
	if n <= 0 && i.dirOffset >= len(entries) {
		return nil, nil
	}

	// read till end
	if n > 0 && i.dirOffset+n > len(entries) {
		err = io.EOF
	}

	if n > 0 && i.dirOffset+n <= len(entries) {
		entries = entries[i.dirOffset : i.dirOffset+n]
		i.dirOffset += n
	} else {
		entries = entries[i.dirOffset:]
		i.dirOffset += len(entries)
	}

	return entries, err
}

// Close does not do anything for NTFS items.
func (i *Item) Close() error { return nil }

// Stat returns the MBR pseudo roots itself as fs.FileMode.
func (i *Item) Stat() (fs.FileInfo, error) {
	infos := parser.Stat(i.ntfsCtx, i.entry)

	return &DirEntry{infos[0]}, nil
}

/*
// IsDir returns if the item is a file.
func (i *Item) IsDir() bool { return i.entry.IsDir(i.ntfsCtx) }

// ModTime returns the zero time (0001-01-01 00:00).
func (i *Item) ModTime() time.Time {
	return time.Time{}
}

// Mode returns the fs.FileMode.
func (i *Item) Mode() fs.FileMode {
	if i.IsDir() {
		return fs.ModeDir
	}
	return 0
}

// Sys returns a map of NTFS item attributes.
func (i *Item) Sys() interface{} {
	infos, err := parser.ModelMFTEntry(i.ntfsCtx, i.entry)
	if err != nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"created":     infos.SI_Times.CreateTime.UTC().Format(time.RFC3339Nano),
		"modified":    infos.SI_Times.FileModifiedTime.UTC().Format(time.RFC3339Nano),
		"mftModified": infos.SI_Times.MFTModifiedTime.UTC().Format(time.RFC3339Nano),
		"accessed":    infos.SI_Times.AccessedTime.UTC().Format(time.RFC3339Nano),
	}
}
*/
