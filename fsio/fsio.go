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
	"errors"
	"io"
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
	pos, err := da.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}
	_, err = da.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}
	n, err = da.Read(p)
	if err != nil {
		return 0, err
	}
	_, err = da.Seek(pos, io.SeekStart)
	return n, err
}

var ErrNotResetSeek = errors.New("could not reset position")

type ErrSizeNotGet struct{ underlying error }

func (e *ErrSizeNotGet) Error() string { return "could not get size" }
func (e *ErrSizeNotGet) Unwrap() error { return e.underlying }
func (e *ErrSizeNotGet) Is(target error) bool {
	var x *ErrSizeNotGet
	return errors.As(target, &x)
}

type WrapperError struct {
	error
	Underlying error
}

func (e *WrapperError) Error() string   { return e.Error() + ": " + e.Underlying.Error() }
func (e *WrapperError) Unwrap() error   { return e.Underlying }
func (e *WrapperError) Is(t error) bool { return errors.Is(e.error, t) || errors.Is(e.Underlying, t) }

// GetSize return the size of an io.Seeker without changing the current offset
func GetSize(seeker io.Seeker) (int64, error) {
	pos, err := seeker.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, &ErrSizeNotGet{err}
	}
	end, err := seeker.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, &ErrSizeNotGet{err}
	}
	_, err = seeker.Seek(pos, io.SeekStart)
	if err != nil {
		return end, &WrapperError{ErrNotResetSeek, err}
	}
	return end, nil
}
