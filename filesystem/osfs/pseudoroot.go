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
	"os"
	"syscall"
	"time"
)

// Root is a pseudo root directory for windows partitions.
type Root struct{  }

func (r *Root) Read([]byte) (int, error) {
	return 0, syscall.EPERM
}
// Name always returns / for window pseudo roots.
func (*Root) Name() (name string) { return "." }

// Close does not do anything for window pseudo roots.
func (*Root) Close() error { return nil }

// Size returns 0 for window pseudo roots.
func (*Root) Size() int64 { return 0 }

// Mode returns os.ModeDir for window pseudo roots.
func (*Root) Mode() os.FileMode { return os.ModeDir }

// ModTime returns the zero time (0001-01-01 00:00) for window pseudo roots.
func (*Root) ModTime() time.Time { return time.Time{} }

// IsDir returns true for window pseudo roots.
func (*Root) IsDir() bool { return true }

// Sys returns nil for window pseudo roots.
func (*Root) Sys() interface{} { return nil }

// Stat returns the windows pseudo roots itself as os.FileMode.
func (r *Root) Stat() (os.FileInfo, error) {
	return r, nil
}
