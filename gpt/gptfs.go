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

// Package gpt provides an io/fs implementation of the GUID partition table
// (GPT).
package gpt

import (
	"fmt"
	"io"
	"io/fs"
	"strconv"
	"strings"
)

// FS implements a read-only file system for Master Boot Records (MBR).
type FS struct {
	gpt *GptPartitionTable
}

// New creates a new gpt FS.
func New(decoder io.ReadSeeker) (*FS, error) {
	gpt := GptPartitionTable{}
	err := gpt.Decode(decoder)
	return &FS{gpt: &gpt}, err
}

// Open returns a File for the given location.
func (fsys *FS) Open(name string) (fs.File, error) {
	valid := fs.ValidPath(name)
	if !valid {
		return nil, fmt.Errorf("path %s invalid", name)
	}

	if name == "." {
		return &Root{gpt: fsys.gpt}, nil
	}
	if !strings.HasPrefix(name, "p") {
		return nil, fmt.Errorf("needs to start with 'p' is %s", name)
	}
	name = name[1:]
	index, err := strconv.Atoi(name)
	if err != nil {
		return nil, err
	}
	partitionEntry := fsys.gpt.Primary().Entries()[index]
	f := NewPartition(index, &partitionEntry)
	return f, nil
}
