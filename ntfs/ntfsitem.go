package ntfs

import (
	"io/fs"
	"os"
	"time"

	"www.velocidex.com/golang/go-ntfs/parser"
)

// Item describes files and directories in the NTFS.
type Item struct {
	entry   *parser.MFT_ENTRY
	name    string
	offset  int64
	path    string
	ntfsCtx *parser.NTFSContext
}

// Name returns the name of the file.
func (i *Item) Name() (name string) { return i.name }

// Read reads bytes into the passed buffer.
func (i *Item) Read(p []byte) (n int, err error) {
	c, err := i.ReadAt(p, i.offset)
	i.offset += int64(c)
	return c, err
}

// ReadAt reads bytes starting at off into passed buffer.
func (i *Item) ReadAt(p []byte, off int64) (n int, err error) {
	attribute, err := i.entry.GetAttribute(i.ntfsCtx, 128, -1)
	if err != nil {
		return 0, err
	}
	return attribute.Data(i.ntfsCtx).ReadAt(p, off)
}

// Seek move the current offset to the given position.
func (i *Item) Seek(pos int64, whence int) (offset int64, err error) {
	switch whence {
	case os.SEEK_SET:
		i.offset = pos
	case os.SEEK_CUR:
		i.offset += pos
	case os.SEEK_END:
		i.offset = i.Size() - pos
	}

	return i.offset, nil
}

// Size returns the item's size.
func (i *Item) Size() int64 {
	infos, err := parser.ModelMFTEntry(i.ntfsCtx, i.entry)
	if err != nil {
		return 0
	}
	return infos.Size
}

// ReadDir returns up to n child items of a directory.
func (i *Item) ReadDir(n int) (entries []fs.DirEntry, err error) {
	infos := parser.ListDir(i.ntfsCtx, i.entry)

	for _, info := range infos {
		if n != 0 && len(entries) == n {
			break
		}
		if info.Name == "" || info.Name == "." {
			continue
		}
		entries = append(entries, &DirEntry{info})
		// TODO: some paths like $BadClus:$Bad are not listed
	}
	return
}

// Close does not do anything for NTFS items.
func (i *Item) Close() error { return nil }

// Stat returns the MBR pseudo roots itself as fs.FileMode.
func (i *Item) Stat() (fs.FileInfo, error) { return i, nil }

// IsDir returns if the item is a file.
func (i *Item) IsDir() bool { return i.entry.IsDir(i.ntfsCtx) }

// ModTime returns the zero time (0001-01-01 00:00).
func (i *Item) ModTime() time.Time { return time.Time{} }

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
