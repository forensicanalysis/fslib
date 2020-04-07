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

package content

import (
	"bytes"
	"io"
	"testing"

	"github.com/forensicanalysis/fslib/fsio"
)

func TestStringsReaderFail(t *testing.T) {
	type args struct {
		r io.Reader
		w io.Writer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"strings reader", args{bytes.NewBuffer([]byte{0x65, 0x6c, 0x6c, 0x6f, 0x00}), &fsio.ErrorWriter{}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StringsReader(tt.args.r, tt.args.w)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringsReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestStringsReader(t *testing.T) {
	type args struct {
		r io.Reader
		w *bytes.Buffer
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"strings reader", args{bytes.NewBuffer([]byte("hello")), &bytes.Buffer{}}, "hello\n", false},
		{"error reader", args{&fsio.ErrorReader{}, &bytes.Buffer{}}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := StringsReader(tt.args.r, tt.args.w)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringsReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := tt.args.w.String(); gotW != tt.wantW {
				t.Errorf("StringsReader() gotW = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_extractString(t *testing.T) {
	in := []byte{0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x68, 0x00}
	type args struct {
		data          []byte
		currentString *bytes.Buffer
		w             io.Writer
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"extract", args{in, bytes.NewBuffer([]byte{0x68}), &bytes.Buffer{}}, false},
		{"write error", args{in, bytes.NewBuffer([]byte{0x68}), &fsio.ErrorWriter{}}, true},
		{"write 2. error", args{[]byte{0x00}, bytes.NewBuffer([]byte{0x65, 0x6c, 0x6c, 0x6f}), &fsio.ErrorWriter{Skip: 1}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := extractString(tt.args.data, tt.args.currentString, tt.args.w); (err != nil) != tt.wantErr {
				t.Errorf("extractString() error = %v, wantErr %v", err, tt.wantErr)
			}
			/* if result.String() != tt.want {
				t.Errorf("extractString() = %v, want %v", result.String(), tt.want)
			} */
		})
	}
}