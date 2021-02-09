// Copyright (c) 2019-2020 Siemens AG
// Copyright (c) 2019-2021 Jonas Plum
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

// Package fat16 provides an io/fs implementation of the FAT16 file systems.
package fat16

import (
	"encoding/binary"
	"fmt"
	"github.com/forensicanalysis/fslib/fsio"
	"io"
	"io/fs"
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

// Open opens a file for reading.
func (m *FS) Open(name string) (f fs.File, err error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}
	// log.Println(">>", name)

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
