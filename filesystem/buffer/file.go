package buffer

import (
	"errors"
	"github.com/forensicanalysis/fslib"
	"io"
	"io/fs"
	"os"
	"syscall"
)

var (
	ErrOutOfRange = errors.New("out of range")
)

type File struct {
	name          string
	fs            *FS
	file          fs.File
	size          int64
	offset        int64
	isdir, closed bool
	buf           []byte
}

func (f *File) fillBuffer(offset int64) (err error) {
	if offset > f.size {
		offset = f.size
		err = io.EOF
	}
	if len(f.buf) >= int(offset) {
		return
	}
	buf := make([]byte, int(offset)-len(f.buf))
	if n, readErr := io.ReadFull(f.file, buf); n > 0 {
		f.buf = append(f.buf, buf[:n]...)
	} else if readErr != nil {
		err = readErr
	}
	return
}

func (f *File) Close() (err error) {
	f.file = nil
	f.closed = true
	f.buf = nil
	return
}

func (f *File) Read(p []byte) (n int, err error) {
	if f.isdir {
		return 0, syscall.EISDIR
	}
	if f.closed {
		return 0, fs.ErrClosed
	}
	err = f.fillBuffer(f.offset + int64(len(p)))
	n = copy(p, f.buf[f.offset:])
	f.offset += int64(n)
	return
}

func (f *File) ReadAt(p []byte, off int64) (n int, err error) {
	if f.isdir {
		return 0, syscall.EISDIR
	}
	if f.closed {
		return 0, fs.ErrClosed
	}
	err = f.fillBuffer(off + int64(len(p)))
	n = copy(p, f.buf[int(off):])
	return
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.isdir {
		return 0, syscall.EISDIR
	}
	if f.closed {
		return 0, fs.ErrClosed
	}
	switch whence {
	case os.SEEK_SET:
	case os.SEEK_CUR:
		offset += f.offset
	case os.SEEK_END:
		offset += f.size
	default:
		return 0, syscall.EINVAL
	}
	if offset < 0 || offset > f.size {
		return 0, ErrOutOfRange
	}
	f.offset = offset
	return offset, nil
}

func (f *File) Name() string {
	return f.name
}

func (f *File) ReadDir(count int) (fi []fs.DirEntry, err error) {
	return fslib.ReadDir(f.file, count)
}

func (f *File) Stat() (os.FileInfo, error) {
	return f.file.Stat()
}
