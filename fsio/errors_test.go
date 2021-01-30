// Copyright (c) 2019-2020 Siemens AG
// Copyright (c) 2019-2021 Jonas Plum
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
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		wantErr bool
	}{
		{"ErrorReader", &ErrorReader{}, true},
		{"ErrorReadSeeker", &ErrorReadSeeker{}, true},
		{"ErrorReadSeekerAt", &ErrorReadSeekerAt{}, true},
		{"ErrorReader", &ErrorReader{Skip: 1}, false},
		{"ErrorReadSeeker", &ErrorReadSeeker{Skip: 1}, false},
		{"ErrorReadSeekerAt", &ErrorReadSeekerAt{Skip: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.Read(nil)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Read() error =  %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReaderAt(t *testing.T) {
	tests := []struct {
		name    string
		r       io.ReaderAt
		wantErr bool
	}{
		{"ErrorReaderAt", &ErrorReaderAt{}, true},
		{"ErrorReadSeekerAt", &ErrorReadSeekerAt{}, true},
		{"ErrorReaderAt", &ErrorReaderAt{Skip: 1}, false},
		{"ErrorReadSeekerAt", &ErrorReadSeekerAt{Skip: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.ReadAt(nil, 0)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadAt() error =  %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSeek(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Seeker
		wantErr bool
	}{
		{"ErrorSeeker", &ErrorSeeker{}, true},
		{"ErrorReadSeeker", &ErrorReadSeeker{}, true},
		{"ErrorReadSeekerAt", &ErrorReadSeekerAt{}, true},
		{"ErrorSeeker", &ErrorSeeker{Skip: 1}, false},
		{"ErrorReadSeeker", &ErrorReadSeeker{Skip: 1}, false},
		{"ErrorReadSeekerAt", &ErrorReadSeekerAt{Skip: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.Seek(0, 0)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Seek() error =  %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Writer
		wantErr bool
	}{
		{"ErrorWriter", &ErrorWriter{}, true},
		{"ErrorWriter", &ErrorWriter{Skip: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.Write(nil)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Write() error =  %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
