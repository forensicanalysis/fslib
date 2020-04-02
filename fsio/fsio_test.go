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

package fsio

import (
	"bytes"
	"errors"
	"io"
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
		wantErr          error
	}{
		{"get zero size", args{bytes.NewReader([]byte{})}, 0, 0, true, nil},
		{"get size", args{bytes.NewReader([]byte{0})}, 0, 1, true, nil},
		{"keep position", args{bytes.NewReader([]byte{0, 1, 2, 3})}, 2, 4, true, nil},

		{"fail 1. seek", args{&ErrorSeeker{Skip: 0 + seeksInTestSetup, Size: 4}}, 0, 0, false, &ErrSizeNotGet{}},
		{"fail 1. seek", args{&ErrorSeeker{Skip: 0 + seeksInTestSetup, Size: 4}}, 0, 0, false, ErrSeek},
		{"fail 2. seek", args{&ErrorSeeker{Skip: 1 + seeksInTestSetup, Size: 4}}, 0, 0, false, &ErrSizeNotGet{}},
		{"fail 2. seek", args{&ErrorSeeker{Skip: 1 + seeksInTestSetup, Size: 4}}, 0, 0, false, ErrSeek},
		{"fail 3. seek but get size", args{&ErrorSeeker{Skip: 2 + seeksInTestSetup, Size: 4}}, 0, 4, false, ErrNotResetSeek},
		{"fail 3. seek but get size", args{&ErrorSeeker{Skip: 2 + seeksInTestSetup, Size: 4}}, 0, 4, false, ErrSeek},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			_, err := tt.args.seeker.Seek(tt.currentPosition, io.SeekStart)
			if err != nil {
				t.Error(err)
			}

			// When
			got, err := GetSize(tt.args.seeker)

			// Then
			if err == nil {
				if tt.wantErr != nil {
					t.Errorf("GetSize() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			} else {
				if tt.wantErr == nil {
					t.Errorf("GetSize() error = %v, wantErr %v", err, tt.wantErr)
					return
				} else {
					// fmt.Println(err, "--", tt.wantErr)
					if !errors.Is(err, tt.wantErr) {
						t.Errorf("GetSize() type error = %v, wantErr %v", err, tt.wantErr)
						return
					}
				}
			}

			if got != tt.wantSize {
				t.Errorf("GetSize() got = %v, wantSize %v", got, tt.wantSize)
			}
			if tt.wantKeepPosition {
				positionAfterGetSize, err := tt.args.seeker.Seek(0, io.SeekCurrent)
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
