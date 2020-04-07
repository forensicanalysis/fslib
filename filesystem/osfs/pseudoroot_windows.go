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

package osfs

import (
	"errors"
	"strings"
	"syscall"
	"unsafe"
)

// Readdirnames lists all partitions in the window pseudo root.
func (*Root) Readdirnames(n int) (partitions []string, err error) {
	kernel32, _ := syscall.LoadDLL("kernel32.dll")
	getLogicalDriveStringsHandle, _ := kernel32.FindProc("GetLogicalDriveStringsA")

	buffer := [1024]byte{}
	bufferSize := uint32(len(buffer))

	hr, _, _ := getLogicalDriveStringsHandle.Call(
		uintptr(unsafe.Pointer(&bufferSize)), //nolint:gosec
		uintptr(unsafe.Pointer(&buffer)),     //nolint:gosec
	)
	if hr == 0 {
		return nil, errors.New("partitions could not be listed")
	}
	for i := 0; i < int(hr); i += 4 {
		partitions = append(partitions, strings.ToUpper(string(buffer[i])))
	}

	return partitions, nil
}
