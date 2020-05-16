// Copyright (c) 2019-2020 Siemens AG
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

// Package systemfs provides a forensicfs implementation that uses the osfs as
// default, while a ntfs for every partition as a fallback on Windows, on UNIX the
// behavior is the same as osfs.
package systemfs

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem"
	"github.com/forensicanalysis/fslib/filesystem/ntfs"
	"github.com/forensicanalysis/fslib/filesystem/osfs"
)

// New creates a new system FS.
func New() (fslib.FS, error) {
	if runtime.GOOS != "windows" {
		return osfs.New(), nil
	}

	fs := &FS{}
	partitions, err := listPartitions()
	if err != nil {
		return fs, err
	}

	fs.ntfsPartitions = strings.Join(partitions, "")
	ntfsPartitions := ""

	for _, partition := range partitions {
		f, err := fs.Open("/" + partition + "/$MFT")

		if err == nil {
			ntfsPartitions += partition
		}

		if closer, ok := f.(io.Closer); ok {
			closer.Close() // nolint:errcheck
		}
	}
	fs.ntfsPartitions = ntfsPartitions
	return fs, nil
}

// FS implements a read-only file system for all operating systems.
type FS struct {
	ntfsPartitions string
}

// Name returns the name of the file system.
func (*FS) Name() (name string) { return "System FS" }

// Open opens a file for reading.
func (systemfs *FS) Open(name string) (item fslib.Item, err error) {
	name, err = filesystem.Clean(name)
	if err != nil {
		return nil, err
	}

	fs := osfs.New()

	item, err = fs.Open(name)
	if err == nil {
		return item, nil
	}
	if os.IsNotExist(err) && path.Base(name)[0] != '$' {
		return nil, err
	}

	if !strings.ContainsRune(systemfs.ntfsPartitions, rune(name[1])) {
		return nil, err
	}

	item, _, err = systemfs.NTFSOpen(name)
	return item, err
}

func (systemfs *FS) NTFSOpen(name string) (fslib.Item, func() error, error) {
	base, err := os.Open(fmt.Sprintf("\\\\.\\%c:", name[1]))
	if err != nil {
		return nil, nil, fmt.Errorf("ntfs base open failed: %w", err)
	}

	lowLevelFS, err := ntfs.New(base)
	if err != nil {
		base.Close() // nolint:errcheck
		return nil, nil, fmt.Errorf("ntfs creation failed: %w", err)
	}

	log.Printf("low level open %s", name[2:])

	item, err := lowLevelFS.Open(name[2:])
	if err != nil {
		return nil, nil, err
	}
	i := &Item{Item: item, base: base}
	return i, i.Close, nil
}

// Stat returns an os.FileInfo object that describes a file.
func (systemfs *FS) Stat(name string) (info os.FileInfo, err error) {
	name, err = filesystem.Clean(name)
	if err != nil {
		return nil, err
	}

	fs := osfs.New()

	info, err = fs.Stat(name)
	if err == nil {
		return info, nil
	}
	if os.IsNotExist(err) && path.Base(name)[0] != '$' {
		return info, err
	}

	if !strings.ContainsRune(systemfs.ntfsPartitions, rune(name[1])) {
		return info, err
	}

	base, err := os.Open(fmt.Sprintf("\\\\.\\%c:", name[1]))
	if err != nil {
		return nil, fmt.Errorf("ntfs base open failed: %w", err)
	}

	lowLevelFS, err := ntfs.New(base)
	if err != nil {
		base.Close() // nolint:errcheck
		return info, fmt.Errorf("ntfs creation failed: %w", err)
	}

	log.Printf("low level open %s", name[2:])

	info, err = lowLevelFS.Stat(name[2:])
	return info, err
}

// Item describes files and directories in the file system.
type Item struct {
	fslib.Item
	base *os.File
}

// Close closes the file freeing the resource. Usually additional IO operations
// fail after closing.
func (i *Item) Close() error {
	i.Item.Close() // nolint:errcheck
	return i.base.Close()
}
