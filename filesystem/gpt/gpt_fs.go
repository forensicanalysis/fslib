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

// Package gpt provides a forensicfs implementation of the GUID partition table
// (GPT).
package gpt

import (
	"fmt"
	"io"
	fsys "io/fs"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/forensicanalysis/fslib/filesystem"
	"github.com/forensicanalysis/fslib/forensicfs"
	"github.com/forensicanalysis/fslib/fsio"
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

// Name returns the filesystem name.
func (fs *FS) Name() string {
	return "GPT"
}

// Open returns a File for the given location.
func (fs *FS) Open(name string) (fsys.File, error) {
	name, err := filesystem.Clean(name)
	if err != nil {
		return nil, err
	}

	if name == "/" {
		return &Root{gpt: fs.gpt}, nil
	}
	if !strings.HasPrefix(name, "/p") {
		return nil, fmt.Errorf("needs to start with '/p' is %s", name)
	}
	name = name[2:]
	index, err := strconv.Atoi(name)
	if err != nil {
		return nil, err
	}
	partitionEntry := fs.gpt.Primary().Entries()[index]
	f := NewPartition(index, &partitionEntry)
	return f, nil
}

// Stat returns an os.FileInfo object that describes a partition.
func (fs *FS) Stat(name string) (os.FileInfo, error) {
	f, err := fs.Open(name)
	if err != nil {
		return nil, err
	}
	return f.Stat()
}

// Partition implements fs.File
type Partition struct {
	forensicfs.FileDefaults
	*io.SectionReader
	name      int
	partition *PartitionEntry
}

// NewPartition creates a new Partition object for parsing GPT partitions.
func NewPartition(name int, partition *PartitionEntry) *Partition {
	return &Partition{
		name:      name,
		partition: partition,
		SectionReader: io.NewSectionReader(
			&fsio.DecoderAtWrapper{ReadSeeker: partition.decoder},
			int64(partition.FirstLba()*512),
			int64(partition.LastLba()-partition.FirstLba()*512),
		),
	}
}

// Name returns the name of a partition that consists of 'pX' where X is the
// number of the partition.
func (p *Partition) Name() string { return "p" + strconv.Itoa(p.name) }

// IsDir returns false for partition.
func (*Partition) IsDir() bool { return false }

// Size returns the partition size.
func (p *Partition) Size() int64 {
	return int64((p.partition.LastLba() - p.partition.FirstLba() + 1) * 512)
}

// Close does not do anything for GPT partitions.
func (p *Partition) Close() error { return nil }

// Stat return an os.FileInfo object that describes a file.
func (p *Partition) Stat() (os.FileInfo, error) { return p, nil }

// Mode returns 0 for partitions.
func (p *Partition) Mode() os.FileMode { return 0 }

// ModTime returns the zero time (0001-01-01 00:00) for partitions.
func (p *Partition) ModTime() time.Time { return time.Time{} }

// Sys returns the PartitionEntry.
func (p *Partition) Sys() interface{} { return p.partition }

func (p *Partition) Type() fsys.FileMode { return p.Mode() }

func (p *Partition) Info() (fsys.FileInfo, error) { return p, nil }


// Root is a pseudo root directory containing the partitions.
type Root struct {
	forensicfs.DirectoryDefaults
	gpt *GptPartitionTable
}

// Name always returns '/' for GPT roots.
func (r *Root) Name() string { return "/" }

func (r *Root) ReadDir(count int) ([]fsys.DirEntry, error) {
	var partitionInfos []fsys.DirEntry
	for index, partition := range r.gpt.Primary().Entries() {
		if count != 0 && index == count {
			return partitionInfos, nil
		}
		if partition.FirstLba() != 0 || partition.LastLba() != 0 {
			p := NewPartition(index, &partition)
			partitionInfos = append(partitionInfos, p)
		}
	}
	return partitionInfos, nil
}

// Readdirnames lists all partitions in the GPT.
func (r *Root) Readdirnames(count int) ([]string, error) {
	var partitionInfos []string
	for index, partition := range r.gpt.Primary().Entries() {
		if count != 0 && index == count {
			return partitionInfos, nil
		}
		if partition.FirstLba() != 0 || partition.LastLba() != 0 {
			partitionInfos = append(partitionInfos, "p"+strconv.Itoa(index))
		}
	}
	return partitionInfos, nil
}

// Size returns 0 for GPT pseudo roots.
func (r *Root) Size() int64 { return 0 }

// Mode returns os.ModeDir for GPT pseudo roots.
func (r *Root) Mode() os.FileMode { return os.ModeDir }

// ModTime returns the zero time (0001-01-01 00:00) for GPT pseudo roots.
func (r *Root) ModTime() time.Time { return time.Time{} }

// IsDir returns true for GPT pseudo roots.
func (r *Root) IsDir() bool { return true }

// Sys returns nil for GPT pseudo roots.
func (r *Root) Sys() interface{} { return nil }

// Close does not do anything for GPT pseudo roots.
func (r *Root) Close() error { return nil }

// Stat returns the GPT pseudo roots itself as os.FileMode.
func (r *Root) Stat() (os.FileInfo, error) { return r, nil }
