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

package forensicfs

import (
	"os"
	"testing"
)

func TestFileInfoDefaults_IsDir(t *testing.T) {
	tests := []struct {
		name string
		f    *FileInfoDefaults
		want bool
	}{
		{"no dir", &FileInfoDefaults{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FileInfoDefaults{}
			if got := f.IsDir(); got != tt.want {
				t.Errorf("FileInfoDefaults.IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectoryDefaults_Read(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		d       *DirectoryDefaults
		args    args
		wantN   int
		wantErr bool
	}{
		{"get error", &DirectoryDefaults{}, args{[]byte{}}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DirectoryDefaults{}
			gotN, err := d.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("DirectoryDefaults.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("DirectoryDefaults.Read() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestDirectoryDefaults_ReadAt(t *testing.T) {
	type args struct {
		p   []byte
		off int64
	}
	tests := []struct {
		name    string
		d       *DirectoryDefaults
		args    args
		wantN   int
		wantErr bool
	}{
		{"get error", &DirectoryDefaults{}, args{[]byte{}, 0}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DirectoryDefaults{}
			gotN, err := d.ReadAt(tt.args.p, tt.args.off)
			if (err != nil) != tt.wantErr {
				t.Errorf("DirectoryDefaults.ReadAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("DirectoryDefaults.ReadAt() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestDirectoryDefaults_Seek(t *testing.T) {
	type args struct {
		offset int64
		whence int
	}
	tests := []struct {
		name    string
		d       *DirectoryDefaults
		args    args
		want    int64
		wantErr bool
	}{
		{"get error", &DirectoryDefaults{}, args{0, 0}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DirectoryDefaults{}
			got, err := d.Seek(tt.args.offset, tt.args.whence)
			if (err != nil) != tt.wantErr {
				t.Errorf("DirectoryDefaults.Seek() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DirectoryDefaults.Seek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectoryInfoDefaults_IsDir(t *testing.T) {
	tests := []struct {
		name string
		d    *DirectoryInfoDefaults
		want bool
	}{
		{"get false", &DirectoryInfoDefaults{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DirectoryInfoDefaults{}
			if got := d.IsDir(); got != tt.want {
				t.Errorf("DirectoryInfoDefaults.IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectoryInfoDefaults_Mode(t *testing.T) {
	tests := []struct {
		name string
		d    *DirectoryInfoDefaults
		want os.FileMode
	}{
		{"ModeDir", &DirectoryInfoDefaults{}, os.ModeDir},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DirectoryInfoDefaults{}
			if got := d.Mode(); got != tt.want {
				t.Errorf("DirectoryInfoDefaults.Mode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDirectoryInfoDefaults_Size(t *testing.T) {
	tests := []struct {
		name string
		d    *DirectoryInfoDefaults
		want int64
	}{
		{"Size", &DirectoryInfoDefaults{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DirectoryInfoDefaults{}
			if got := d.Size(); got != tt.want {
				t.Errorf("DirectoryInfoDefaults.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
