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

package osfs

import (
	"errors"
	"io/fs"
	"strings"
	"syscall"
	"time"
	"unsafe"
	"sort"
)

type SimpleEntry struct {
	name string
}

func (e *SimpleEntry) Name() string               { return e.name }
func (e *SimpleEntry) IsDir() bool                { return true }
func (e *SimpleEntry) Type() fs.FileMode          { return fs.ModeDir }
func (e *SimpleEntry) Info() (fs.FileInfo, error) { return e, nil }
func (e *SimpleEntry) Size() int64                { return 0 }
func (e *SimpleEntry) Mode() fs.FileMode          { return fs.ModeDir }
func (e *SimpleEntry) ModTime() time.Time         { return time.Time{} }
func (e *SimpleEntry) Sys() interface{}           { return nil }

func SimpleEntries(names []string) (entries []fs.DirEntry) {
	for _, name := range names {
		entries = append(entries, &SimpleEntry{name: name})
	}
	return entries
}

// ReadDir lists all partitions in the window pseudo root.
func (r *Root) ReadDir(n int) (partitions []fs.DirEntry, err error) {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return nil, err
	}
	getLogicalDriveStringsProc, err := kernel32.FindProc("GetLogicalDriveStringsA")
	if err != nil {
		return nil, err
	}

	lpBuffer := [1024]byte{}
	nBufferLength := uint32(len(lpBuffer))

	returnLength, _, _ := getLogicalDriveStringsProc.Call(
		uintptr(unsafe.Pointer(&nBufferLength)), //nolint:gosec
		uintptr(unsafe.Pointer(&lpBuffer)),      //nolint:gosec
	)
	if returnLength == 0 {
		return nil, errors.New("partitions could not be listed")
	}
	for i := 0; i < int(returnLength); i += 4 {
		partitions = append(partitions, &SimpleEntry{strings.ToUpper(string(lpBuffer[i]))})
	}

	sort.Sort(fslib.ByName(partitions))

	return partitions, nil
}
