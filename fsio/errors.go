package fsio

import (
	"errors"
	"fmt"
)

type ErrorReader struct{}

func (b *ErrorReader) Read([]byte) (n int, err error) {
	return 0, errors.New("broken reader")
}

type ErrorReaderAt struct{}

func (b *ErrorReaderAt) ReadAt([]byte, int64) (n int, err error) {
	return 0, errors.New("broken readerAt")
}

type ErrorSeeker struct {
	Whence int
}

func (b *ErrorSeeker) Seek(_ int64, whence int) (int64, error) {
	if b.Whence == -1 || whence == b.Whence {
		return 0, errors.New("broken seeker")
	}
	return 0, nil
}

type ErrorReadSeeker struct {
	ErrorSeeker
	ErrorReader
}

type ErrorWriter struct {
	Skip    int
	current int
}

func (e *ErrorWriter) Write(b []byte) (int, error) {
	fmt.Println(string(b))
	if e.current >= e.Skip {
		return 0, errors.New("broken writer")
	}
	e.current += 1
	return len(b), nil
}
