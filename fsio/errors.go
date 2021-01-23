// Copyright (c) 2020 Siemens AG
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

package fsio

import (
	"errors"
	"os"
)

// ErrorReader mocks a Reader that fails after some Reads.
type ErrorReader struct {
	Skip    int
	current int
}

// Read implements io.Read but fails after n attempts.
func (e *ErrorReader) Read(b []byte) (n int, err error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken reader")
	}
	e.current++
	return len(b), nil
}

// ErrorReaderAt mocks a ReaderAt that fails after some ReadAts.
type ErrorReaderAt struct {
	Skip    int
	current int
}

// ReadAt implements io.ReadAt but fails after n attempts.
func (e *ErrorReaderAt) ReadAt(b []byte, _ int64) (n int, err error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken readerAt")
	}
	e.current++
	return len(b), nil
}

// ErrorSeeker mocks a Seeker that fails after some Seek.
type ErrorSeeker struct {
	Skip      int
	Size      int64
	seekCount int
	position  int64
}

// Seek implements io.Seek but fails after n attempts.
func (e *ErrorSeeker) Seek(off int64, whence int) (int64, error) {
	if e.seekCount >= e.Skip {
		return 0, errors.New("seek failed")
	}
	e.seekCount++
	switch whence {
	case os.SEEK_CUR:
		e.position += off
	case os.SEEK_SET:
		e.position = off
	case os.SEEK_END:
		e.position = e.Size + off
	}
	return e.position, nil
}

// ErrorReadSeeker mocks a ReadSeeker that fails after some Reads or Seeks.
type ErrorReadSeeker struct {
	Skip    int
	current int
}

// Read implements io.Read but fails after n attempts.
func (e *ErrorReadSeeker) Read(b []byte) (n int, err error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken reader")
	}
	e.current++
	return len(b), nil
}

// Seek implements io.Seek but fails after n attempts.
func (e *ErrorReadSeeker) Seek(int64, int) (int64, error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken seek")
	}
	e.current++
	return 0, nil
}

// ErrorReadSeekerAt mocks a ReadSeekerAt that fails after some ReadSeeker.
type ErrorReadSeekerAt struct {
	Skip    int
	current int
}

// Read implements io.Read but fails after n attempts.
func (e *ErrorReadSeekerAt) Read(b []byte) (n int, err error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken reader")
	}
	e.current++
	return len(b), nil
}

// Seek implements io.Seek but fails after n attempts.
func (e *ErrorReadSeekerAt) Seek(int64, int) (int64, error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken seek")
	}
	e.current++
	return 0, nil
}

// ReadAt implements io.ReadAt but fails after n attempts.
func (e *ErrorReadSeekerAt) ReadAt(b []byte, _ int64) (n int, err error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken readerAt")
	}
	e.current++
	return len(b), nil
}

// ErrorWriter mocks a Writer that fails after some Writes.
type ErrorWriter struct {
	Skip    int
	current int
}

// Write implements io.Write but fails after n attempts.
func (e *ErrorWriter) Write(b []byte) (int, error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken writer")
	}
	e.current++
	return len(b), nil
}
