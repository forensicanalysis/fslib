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

// Package fat16 provides a forensicfs implementation of the FAT16 file systems.
package fat16

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/forensicanalysis/fslib/fsio"
)

// FS implements a read-only file system for the FAT16 file system.
type FS struct {
	vh      volumeHeader
	decoder fsio.ReadSeekerAt
}

// New creates a new fat16 FS.
func New(decoder fsio.ReadSeekerAt) (*FS, error) {
	// parser volume header
	vh := volumeHeader{}
	err := binary.Read(decoder, binary.LittleEndian, &vh)
	if err != nil && err != io.EOF {
		return nil, err
	}
	decoder.Seek(0, 0) // nolint: errcheck

	return &FS{vh: vh, decoder: decoder}, err
}

/*
func (m *FS) getVolumeName() (string, error) {
	rootDirStart := (int64(m.vh.SectorsPerFat)*int64(m.vh.FatCount) + 1) * 512

	_, err := decoder.Seek(rootDirStart, os.SEEK_SET)
	if err != nil {
		return "", err
	}

	for i := uint16(0); i < 5; i++ {
		firstByte, err := firstByte(decoder)
		if err != nil {
			return "", err
		}

		// test if entry exists
		if firstByte != 0x00 {
			de := directoryEntry{}

			err := binary.Read(decoder, binary.LittleEndian, &de)
			if err != nil {
				return "", err
			}

			// hide volume label
			if de.FileAttributes&0x08 != 0 {
				return formatFilename(&de), nil
			}
		}
	}
	return "", errors.New("Volumename not found")
}
*/

// Open opens a file for reading.
func (m *FS) Open(name string) (f fs.File, err error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}

	if name == "." {
		name = ""
	}

	name, de, err := m.getDirectoryEntry(2, m.vh.RootdirEntryCount, name)
	if err != nil {
		return nil, err
	}
	f = NewItem(name, m, &de.directoryEntry)

	return f, nil
}

// Item describes files and directories in the FAT16 file system.
type Item struct {
	*io.SectionReader
	name           string
	fs             *FS
	directoryEntry *directoryEntry
}

// NewItem creates a new fat16 Item.
func NewItem(name string, fs *FS, directoryEntry *directoryEntry) *Item {
	log.Println("NewItem directoryEntry.Startingcluster: ", directoryEntry.Startingcluster)
	cluster := int64(directoryEntry.Startingcluster)

	pos := getOffset(cluster, fs.vh)

	size := int64(directoryEntry.FileSize)
	if size == 0 {
		size = int64(fs.vh.SectorSize) * int64(fs.vh.SectorsPerCluster)
	}
	log.Println("directoryEntry.FileSize", directoryEntry.FileSize, "size ", size)

	return &Item{
		name:           name,
		fs:             fs,
		directoryEntry: directoryEntry,
		SectionReader:  io.NewSectionReader(fs.decoder, pos, size),
	}
}

// Name returns the name of the file.
func (i *Item) Name() string {
	return i.name // string(bytes.TrimRight(i.directoryEntry.Filename[:], "\x00"))
}

func (i *Item) ReadDir(n int) ([]fs.DirEntry, error) {
	if !i.IsDir() {
		return nil, errors.New("cannot call Readdirnames on a file")
	}
	de := i.directoryEntry

	size := int64(de.FileSize)
	if size == 0 {
		size = int64(i.fs.vh.SectorSize) * int64(i.fs.vh.SectorsPerCluster)
	}

	log.Printf("Readdirnames startingcluster: %d size: %d", de.Startingcluster, size)
	entries, err := i.fs.getDirectoryEntries(int64(de.Startingcluster), uint16(size/32))
	var infos []fs.DirEntry
	for name, entry := range entries {
		if name != "." && name != ".." {
			infos = append(infos, entry)
			n--
			if n == 0 {
				break
			}
		}
	}
	return infos, err
}

// Read reads bytes into the passed buffer.
func (i *Item) Read(p []byte) (n int, err error) {
	if i.IsDir() {
		return 0, syscall.EPERM
	}
	return i.SectionReader.Read(p)
}

// ReadAt reads bytes starting at off into passed buffer.
func (i *Item) ReadAt(p []byte, off int64) (n int, err error) {
	if i.IsDir() {
		return 0, syscall.EPERM
	}
	return i.SectionReader.ReadAt(p, off)
}

// Seek move the current offset to the given position.
func (i *Item) Seek(offset int64, whence int) (int64, error) {
	if i.IsDir() {
		return 0, syscall.EPERM
	}
	return i.SectionReader.Seek(offset, whence)
}

// Close closes the file freeing the resource. Usually additional IO operations
// fail after closing.
func (*Item) Close() error { return nil }

// Stat return an os.FileInfo object that describes a file.
func (i *Item) Stat() (os.FileInfo, error) { return i, nil }

// Mode returns the os.FileMode.
func (i *Item) Mode() os.FileMode {
	var mode os.FileMode
	if i.directoryEntry.FileAttributes&0x10 != 0 {
		mode |= os.ModeDir
	}
	return mode
}

// ModTime returns the modification time.
func (*Item) ModTime() time.Time { return time.Time{} } // TODO

// Sys returns underlying data source.
func (i *Item) Sys() interface{} { return i.directoryEntry }

// IsDir returns if the item is a file.
func (i *Item) IsDir() bool { return i.directoryEntry.FileAttributes&0x10 != 0 }

// Size returns the item's size.
func (i *Item) Size() int64 { return int64(i.directoryEntry.FileSize) }
