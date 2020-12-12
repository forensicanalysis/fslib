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

// Package fsio provides IO interfaces and functions similar for file system
// operations.
package fsio

import (
	"io"
	"os"
)

// ReadSeekerAt combines the io.Reader, io.Seeker and io.ReaderAt interface.
type ReadSeekerAt interface {
	io.Reader
	io.Seeker
	io.ReaderAt
}

// DecoderAtWrapper wraps an io.ReadSeeker to provide a ReadAt method and implement
// the ReadSeekerAt interface.
type DecoderAtWrapper struct {
	io.ReadSeeker
}

// ReadAt reads len(b) bytes from the File starting at byte offset off.
func (da *DecoderAtWrapper) ReadAt(p []byte, off int64) (n int, err error) {
	pos, err := da.Seek(0, os.SEEK_CUR)
	if err != nil {
		return 0, err
	}
	_, err = da.Seek(off, os.SEEK_SET)
	if err != nil {
		return 0, err
	}
	n, err = da.Read(p)
	if err != nil {
		return 0, err
	}
	_, err = da.Seek(pos, os.SEEK_SET)
	return n, err
}

// GetSize return the size of an io.Seeker without changing the current offset.
func GetSize(seeker io.Seeker) (int64, error) {
	pos, err := seeker.Seek(0, os.SEEK_CUR)
	if err != nil {
		return 0, err
	}
	end, err := seeker.Seek(0, os.SEEK_END)
	if err != nil {
		return 0, err
	}
	_, err = seeker.Seek(pos, os.SEEK_SET)
	return end, err
}
