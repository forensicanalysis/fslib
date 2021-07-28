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

// Package systemfs provides an io/fs implementation that uses the osfs as
// default, while a ntfs for every partition as a fallback on Windows, on UNIX the
// behavior is the same as osfs.
package systemfs

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"runtime"

	"github.com/forensicanalysis/fslib/ntfs"
	"github.com/forensicanalysis/fslib/osfs"
)

// New creates a new system FS.
func New() (fs.FS, error) {
	return newFS()
}

func newFS() (fs.FS, error) {
	if runtime.GOOS != "windows" {
		return osfs.New(), nil
	}

	fsys := &FS{}
	root := osfs.Root{}
	partitions, err := root.ReadDir(0)
	if err != nil {
		return fsys, err
	}

	var ntfsPartitions []string
	for _, partition := range partitions {
		_, close, err := fsys.NTFSOpen(partition.Name() + "/$MFT")

		if err == nil {
			ntfsPartitions = append(ntfsPartitions, partition.Name())
			close()
		}
	}
	fsys.ntfsPartitions = ntfsPartitions

	return fsys, nil
}

// FS implements a read-only file system for all operating systems.
type FS struct {
	ntfsPartitions []string
	cacheSize      int
	pageSize       int
}

func (systemfs *FS) setPageAndCacheSize(pageSize, cacheSize int) {
	systemfs.cacheSize = cacheSize
	systemfs.pageSize = pageSize
}

// Open opens a file for reading.
func (systemfs *FS) Open(name string) (item fs.File, err error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}

	if name == "." {
		return &Root{fs: systemfs}, nil
	}

	fsys := osfs.New()

	item, err = fsys.Open(name)
	if err == nil {
		return item, nil
	}
	if os.IsNotExist(err) && path.Base(name)[0] != '$' {
		return nil, err
	}

	if !contains(systemfs.ntfsPartitions, string(name[0])) {
		return nil, fmt.Errorf("not an NTFS parition (%s %s): %w", string(name[0]), systemfs.ntfsPartitions, err)
	}

	item, _, err = systemfs.NTFSOpen(name)
	if err != nil {
		err = fmt.Errorf("ntfs open failed: %w", err)
	}
	return item, err
}

func (systemfs *FS) NTFSOpen(name string) (fs.File, func() error, error) {
	base, err := os.Open(fmt.Sprintf("\\\\.\\%c:", name[0]))
	if err != nil {
		return nil, nil, fmt.Errorf("ntfs base open failed: %w", err)
	}

	lowLevelFS, err := ntfs.New(base, 0, 0)
	if err != nil {
		base.Close() // nolint:errcheck
		return nil, nil, fmt.Errorf("ntfs creation failed: %w", err)
	}

	log.Printf("low level open %s", name[2:])

	item, err := lowLevelFS.Open(name[2:])
	if err != nil {
		return nil, nil, err
	}

	i := &Item{File: item, base: base}
	return i, i.Close, nil
}

// Stat returns an fs.FileInfo object that describes a file.
func (systemfs *FS) Stat(name string) (info fs.FileInfo, err error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}

	if name == "." {
		return &Root{fs: systemfs}, nil
	}

	fsys := osfs.New()

	info, err = fsys.Stat(name)
	if err == nil {
		return info, nil
	}
	if os.IsNotExist(err) && path.Base(name)[0] != '$' {
		return info, err
	}

	if !contains(systemfs.ntfsPartitions, string(name[0])) {
		return info, fmt.Errorf("not an NTFS parition: %w", err)
	}

	base, err := os.Open(fmt.Sprintf("\\\\.\\%c:", name[1]))
	if err != nil {
		return nil, fmt.Errorf("ntfs base open failed: %w", err)
	}

	lowLevelFS, err := ntfs.New(base, 0, 0)
	if err != nil {
		base.Close() // nolint:errcheck
		return info, fmt.Errorf("ntfs creation failed: %w", err)
	}

	log.Printf("low level open %s", name[2:])

	info, err = fs.Stat(lowLevelFS, name[2:])
	return info, fmt.Errorf("ntfs stat failed: %w", err)
}

// Item describes files and directories in the file system.
type Item struct {
	fs.File
	base *os.File
}

// Close closes the file freeing the resource. Usually additional IO operations
// fail after closing.
func (i *Item) Close() error {
	i.File.Close() // nolint:errcheck
	return i.base.Close()
}

func contains(l []string, s string) bool {
	for _, e := range l {
		if e == s {
			return true
		}
	}
	return false
}
