// Copyright (c) 2019 Siemens AG
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author(s): Jonas Plum

// Package ntfs provides a forensicfs implementation of the New Technology File
// System (NTFS).
package ntfs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"sort"
	"time"

	"www.velocidex.com/golang/go-ntfs/parser"
)

// New creates a new ntfs FS.
func New(r io.ReaderAt) (fs *FS, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("error parsing file system as NTFS")
		}
	}()
	reader, err := parser.NewPagedReader(r, 1024*1024, 100*1024*1024)
	if err != nil {
		return nil, err
	}
	ntfsCtx, err := parser.GetNTFSContext(reader, 0)
	return &FS{ntfsCtx: ntfsCtx}, err
}

// FS implements a read-only file system for the NTFS.
type FS struct {
	ntfsCtx *parser.NTFSContext
}

// Name returns the name of the file system.
func (*FS) Name() (name string) { return "NTFS" }

// Open opens a file for reading.
func (fsys *FS) Open(name string) (item fs.File, err error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}

	dir, err := fsys.ntfsCtx.GetMFT(5)
	if err != nil {
		return nil, err
	}
	entry, err := dir.Open(fsys.ntfsCtx, name)

	return &Item{entry: entry, name: path.Base(name), path: name, ntfsCtx: fsys.ntfsCtx}, err
}

// Stat returns an os.FileInfo object that describes a file.
func (fsys *FS) Stat(name string) (os.FileInfo, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}
	return f.Stat()
}

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

type DirEntry struct {
	info *parser.FileInfo
}

func (d DirEntry) Name() string {
	return d.info.Name
}

func (d DirEntry) IsDir() bool {
	return d.info.IsDir
}

func (d DirEntry) Size() int64 {
	return d.info.Size
}

func (d DirEntry) Mode() fs.FileMode {
	if d.IsDir() {
		return fs.ModeDir
	}
	return 0
}

func (d DirEntry) ModTime() time.Time {
	return d.info.Mtime
}

func (d DirEntry) Sys() interface{} {
	return d.info
}

func (d DirEntry) Type() fs.FileMode {
	if d.IsDir() {
		return fs.ModeDir
	}
	return 0
}

func (d DirEntry) Info() (fs.FileInfo, error) {
	return d, nil
}

func (i *Item) ReadDir(n int) (entries []fs.DirEntry, err error) {
	infos := parser.ListDir(i.ntfsCtx, i.entry)

	for _, info := range infos {
		if n != 0 && len(entries) == n {
			break
		}
		if info.Name == "" || info.Name == "." {
			continue
		}
		entries = append(entries, DirEntry{info})
		// TODO: some path like $BadClus:$Bad are not listed
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	return
}

// Readdirnames returns up to n child items of a directory.
func (i *Item) Readdirnames(n int) (items []string, err error) {
	infos := parser.ListDir(i.ntfsCtx, i.entry)

	for _, info := range infos {
		if n != 0 && len(items) == n {
			break
		}
		if info.Name == "" || info.Name == "." {
			continue
		}
		items = append(items, info.Name)
		// TODO: some path like $BadClus:$Bad are not listed
	}
	sort.Strings(items)
	return
}

// Close does not do anything for NTFS items.
func (i *Item) Close() error { return nil }

// Stat returns the MBR pseudo roots itself as os.FileMode.
func (i *Item) Stat() (os.FileInfo, error) { return i, nil }

// IsDir returns if the item is a file.
func (i *Item) IsDir() bool { return i.entry.IsDir(i.ntfsCtx) }

// ModTime returns the zero time (0001-01-01 00:00).
func (i *Item) ModTime() time.Time { return time.Time{} }

// Mode returns the os.FileMode.
func (i *Item) Mode() os.FileMode {
	if i.IsDir() {
		return os.ModeDir
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
