package fsio

import (
	"errors"
	"io"
)

type ErrorReader struct {
	Skip    int
	current int
}

func (e *ErrorReader) Read(b []byte) (n int, err error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken reader")
	}
	e.current += 1
	return len(b), nil
}

type ErrorReaderAt struct {
	Skip    int
	current int
}

func (e *ErrorReaderAt) ReadAt(b []byte, _ int64) (n int, err error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken readerAt")
	}
	e.current += 1
	return len(b), nil
}

var ErrSeek = errors.New("seek failed")

type ErrorSeeker struct {
	Skip      int
	Size      int64
	seekCount int
	position  int64
}

func (e *ErrorSeeker) Seek(off int64, whence int) (int64, error) {
	if e.seekCount >= e.Skip {
		return 0, ErrSeek
	}
	e.seekCount += 1
	switch whence {
	case io.SeekCurrent:
		e.position += off
	case io.SeekStart:
		e.position = off
	case io.SeekEnd:
		e.position = e.Size + off
	}
	return e.position, nil
}

type ErrorReadSeeker struct {
	Skip    int
	current int
}

func (e *ErrorReadSeeker) Read(b []byte) (n int, err error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken reader")
	}
	e.current += 1
	return len(b), nil
}

func (e *ErrorReadSeeker) Seek(int64, int) (int64, error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken seek")
	}
	e.current += 1
	return 0, nil
}

type ErrorReadSeekerAt struct {
	Skip    int
	current int
}

func (e *ErrorReadSeekerAt) Read(b []byte) (n int, err error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken reader")
	}
	e.current += 1
	return len(b), nil
}

func (e *ErrorReadSeekerAt) Seek(int64, int) (int64, error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken seek")
	}
	e.current += 1
	return 0, nil
}
func (e *ErrorReadSeekerAt) ReadAt(b []byte, _ int64) (n int, err error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken readerAt")
	}
	e.current += 1
	return len(b), nil
}

type ErrorWriter struct {
	Skip    int
	current int
}

func (e *ErrorWriter) Write(b []byte) (int, error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken writer")
	}
	e.current += 1
	return len(b), nil
}
