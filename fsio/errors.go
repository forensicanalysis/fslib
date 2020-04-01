package fsio

import (
	"errors"
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

type ErrorSeeker struct {
	Skip    int
	current int
}

func (e *ErrorSeeker) Seek(int64, int) (int64, error) {
	if e.current >= e.Skip {
		return 0, errors.New("broken seek")
	}
	e.current += 1
	return 0, nil
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
