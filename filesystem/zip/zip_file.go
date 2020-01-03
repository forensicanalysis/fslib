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

package zip

import (
	"archive/zip"
	"io"
	"os"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/afero"
)

// File describes files and directories in the zip file system.
type File struct {
	fs            *FS
	zipfile       *zip.File
	reader        io.ReadCloser
	offset        int64
	isdir, closed bool
	buf           []byte
}

func (f *File) fillBuffer(offset int64) (err error) {
	if f.reader == nil {
		if f.reader, err = f.zipfile.Open(); err != nil {
			return
		}
	}
	if offset > int64(f.zipfile.UncompressedSize64) {
		offset = int64(f.zipfile.UncompressedSize64)
		err = io.EOF
	}
	if len(f.buf) >= int(offset) {
		return
	}
	buf := make([]byte, int(offset)-len(f.buf))
	n, _ := io.ReadFull(f.reader, buf)
	if n > 0 {
		f.buf = append(f.buf, buf[:n]...)
	}
	return
}

// Close closes the file freeing the resource. Other IO operations fail after
// closing.
func (f *File) Close() (err error) {
	f.zipfile = nil
	f.closed = true
	f.buf = nil
	if f.reader != nil {
		err = f.reader.Close()
		f.reader = nil
	}
	return
}

// Read reads bytes into the passed buffer.
func (f *File) Read(p []byte) (n int, err error) {
	if f.isdir {
		return 0, syscall.EISDIR
	}
	if f.closed {
		return 0, afero.ErrFileClosed
	}
	err = f.fillBuffer(f.offset + int64(len(p)))
	n = copy(p, f.buf[f.offset:])
	f.offset += int64(len(p))
	return
}

// ReadAt reads bytes starting at off into passed buffer.
func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	if f.isdir {
		return 0, syscall.EISDIR
	}
	if f.closed {
		return 0, afero.ErrFileClosed
	}
	err = f.fillBuffer(off + int64(len(p)))
	n = copy(p, f.buf[int(off):])
	return
}

// Seek move the current offset to the given position.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.isdir {
		return 0, syscall.EISDIR
	}
	if f.closed {
		return 0, afero.ErrFileClosed
	}
	switch whence {
	case io.SeekStart:
	case io.SeekCurrent:
		offset += f.offset
	case io.SeekEnd:
		offset += int64(f.zipfile.UncompressedSize64)
	default:
		return 0, syscall.EINVAL
	}
	if offset < 0 || offset > int64(f.zipfile.UncompressedSize64) {
		return 0, afero.ErrOutOfRange
	}
	f.offset = offset
	return offset, nil
}

// Name returns the name of the file.
func (f *File) Name() string {
	if f.zipfile == nil {
		return "/"
	}
	return strings.TrimSuffix(path.Base(f.zipfile.Name), "/")
}

func (f *File) getDirEntries() (map[string]*zip.File, error) {
	if !f.isdir {
		return nil, syscall.ENOTDIR
	}
	name := "/"
	if f.zipfile != nil {
		name = path.Join(splitpath(f.zipfile.Name))
	}
	entries, ok := f.fs.files[name]
	if !ok {
		return nil, &os.PathError{Op: "readdir", Path: name, Err: syscall.ENOENT}
	}
	return entries, nil
}

// Readdirnames returns up to n child items of a directory.
func (f *File) Readdirnames(count int) (fi []string, err error) {
	zipfiles, err := f.getDirEntries()
	if err != nil {
		return nil, err
	}
	for _, zipfile := range zipfiles {
		fi = append(fi, strings.TrimSuffix(path.Base(zipfile.Name), "/"))
		if count > 0 && len(fi) >= count {
			break
		}
	}
	return
}

// Stat return an os.FileInfo object that describes a file.
func (f *File) Stat() (os.FileInfo, error) {
	if f.zipfile == nil {
		return &RootInfo{}, nil
	}
	return f.zipfile.FileInfo(), nil
}

// Sys returns underlying data source.
func (f *File) Sys() interface{} {
	return map[string]interface{}{
		"mode":     f.zipfile.FileInfo().Mode(),
		"modified": f.zipfile.Modified.In(time.UTC),
	}
}

// RootInfo is a pseudo root os.FileInfo.
type RootInfo struct{}

// Name always returns / for zip pseudo roots.
func (i *RootInfo) Name() string { return "/" }

// Size returns 0 for zip pseudo roots.
func (i *RootInfo) Size() int64 { return 0 }

// Mode returns os.ModeDir for zip pseudo roots.
func (i *RootInfo) Mode() os.FileMode { return os.ModeDir }

// ModTime returns the zero time (0001-01-01 00:00) for zip pseudo roots.
func (i *RootInfo) ModTime() time.Time { return time.Time{} }

// IsDir returns true for zip pseudo roots.
func (i *RootInfo) IsDir() bool { return true }

// Sys returns nil for zip pseudo roots.
func (i *RootInfo) Sys() interface{} { return nil }
