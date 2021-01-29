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
	"bytes"
	"io"
	"os"
	"testing"
)

func TestDecoderAtWrapper_ReadAt(t *testing.T) {
	type fields struct {
		ReadSeeker io.ReadSeeker
	}
	type args struct {
		p   []byte
		off int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
	}{
		{"readat", fields{bytes.NewReader([]byte{1, 2})}, args{[]byte{1}, 0}, 1, false},
		{"readat eof", fields{bytes.NewReader([]byte{})}, args{nil, 0}, 0, true},
		{"fail 1. seek", fields{ReadSeeker: &ErrorReadSeeker{}}, args{nil, 0}, 0, true},
		{"fail 2. seek", fields{ReadSeeker: &ErrorReadSeeker{Skip: 1}}, args{nil, 0}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			da := &DecoderAtWrapper{
				ReadSeeker: tt.fields.ReadSeeker,
			}
			gotN, err := da.ReadAt(tt.args.p, tt.args.off)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("ReadAt() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestGetSize(t *testing.T) {
	type args struct {
		seeker io.Seeker
	}
	seeksInTestSetup := 1
	tests := []struct {
		name             string
		args             args
		currentPosition  int64
		wantSize         int64
		wantKeepPosition bool
		wantErr          bool
	}{
		{"get zero size", args{bytes.NewReader([]byte{})}, 0, 0, true, false},
		{"get size", args{bytes.NewReader([]byte{0})}, 0, 1, true, false},
		{"keep position", args{bytes.NewReader([]byte{0, 1, 2, 3})}, 2, 4, true, false},

		{"fail 1. seek", args{&ErrorSeeker{Skip: 0 + seeksInTestSetup, Size: 4}}, 0, 0, false, true},
		{"fail 2. seek", args{&ErrorSeeker{Skip: 1 + seeksInTestSetup, Size: 4}}, 0, 0, false, true},
		{"fail 3. seek but get size", args{&ErrorSeeker{Skip: 2 + seeksInTestSetup, Size: 4}}, 0, 4, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			_, err := tt.args.seeker.Seek(tt.currentPosition, os.SEEK_SET)
			if err != nil {
				t.Error(err)
			}

			// Run
			got, err := GetSize(tt.args.seeker)

			// Asserts
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantSize {
				t.Errorf("GetSize() got = %v, wantSize %v", got, tt.wantSize)
			}
			if tt.wantKeepPosition {
				positionAfterGetSize, err := tt.args.seeker.Seek(0, os.SEEK_CUR)
				if err != nil {
					t.Error(err)
				}
				if positionAfterGetSize != tt.currentPosition {
					t.Errorf("GetSize() positionAfterGetSize = %v, currentPosition %v", positionAfterGetSize, tt.currentPosition)
				}
			}
		})
	}
}
