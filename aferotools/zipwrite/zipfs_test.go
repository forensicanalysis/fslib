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

package zipwrite

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/afero"
)

func TestNewWriteZipFs(t *testing.T) {
	tests := []struct {
		name       string
		want       *FS
		wantWriter string
	}{
		{"Create write FS", &FS{zipwriter: zip.NewWriter(&bytes.Buffer{})}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if got := New(writer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWriteZipFs() = %v, want %v", got, tt.want)
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("NewWriteZipFs() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestZipWriteFs_Open(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fs      *FS
		args    args
		want    afero.File
		wantErr bool
	}{
		{"Open", &FS{zipwriter: zip.NewWriter(&bytes.Buffer{})}, args{"foo"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fs.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FS.Open() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipWriteFs_OpenFile(t *testing.T) {
	type args struct {
		name string
		flag int
		perm os.FileMode
	}
	tests := []struct {
		name    string
		fs      *FS
		args    args
		want    afero.File
		wantErr bool
	}{
		{"Open ro", &FS{zipwriter: zip.NewWriter(&bytes.Buffer{})}, args{"foo", os.O_RDONLY, 0}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fs.OpenFile(tt.args.name, tt.args.flag, tt.args.perm)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.OpenFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FS.OpenFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipWriteFs_Remove(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fs      *FS
		args    args
		wantErr bool
	}{
		{"Remove", &FS{zipwriter: zip.NewWriter(&bytes.Buffer{})}, args{"foo"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fs.Remove(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("FS.Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZipWriteFs_RemoveAll(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fs      *FS
		args    args
		wantErr bool
	}{
		{"Remove All", &FS{zipwriter: zip.NewWriter(&bytes.Buffer{})}, args{"foo"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fs.RemoveAll(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("FS.RemoveAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZipWriteFs_Rename(t *testing.T) {
	type args struct {
		oldname string
		newname string
	}
	tests := []struct {
		name    string
		fs      *FS
		args    args
		wantErr bool
	}{
		{"Rename", &FS{zipwriter: zip.NewWriter(&bytes.Buffer{})}, args{"foo", "bar"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fs.Rename(tt.args.oldname, tt.args.newname); (err != nil) != tt.wantErr {
				t.Errorf("FS.Rename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZipWriteFs_Stat(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fs      *FS
		args    args
		want    os.FileInfo
		wantErr bool
	}{
		{"Stat", &FS{zipwriter: zip.NewWriter(&bytes.Buffer{})}, args{"foo"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fs.Stat(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.Stat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FS.Stat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipWriteFs_Name(t *testing.T) {
	tests := []struct {
		name string
		fs   *FS
		want string
	}{
		{"Name", &FS{zipwriter: zip.NewWriter(&bytes.Buffer{})}, "FS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fs.Name(); got != tt.want {
				t.Errorf("FS.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipWriteFs_Chmod(t *testing.T) {
	type args struct {
		name string
		mode os.FileMode
	}
	tests := []struct {
		name    string
		fs      *FS
		args    args
		wantErr bool
	}{
		{"Chmod", &FS{zipwriter: zip.NewWriter(&bytes.Buffer{})}, args{"foo", 0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fs.Chmod(tt.args.name, tt.args.mode); (err != nil) != tt.wantErr {
				t.Errorf("FS.Chmod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZipWriteFs_Chtimes(t *testing.T) {
	type args struct {
		name  string
		atime time.Time
		mtime time.Time
	}
	tests := []struct {
		name    string
		fs      *FS
		args    args
		wantErr bool
	}{
		{"Remove", &FS{zipwriter: zip.NewWriter(&bytes.Buffer{})}, args{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fs.Chtimes(tt.args.name, tt.args.atime, tt.args.mtime); (err != nil) != tt.wantErr {
				t.Errorf("FS.Chtimes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZipWriteItem_Name(t *testing.T) {
	tests := []struct {
		name string
		item Item
		want string
	}{
		{"Name", Item{"foo", nil}, "foo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.item.Name(); got != tt.want {
				t.Errorf("Item.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipWriteItem_Read(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		item    Item
		args    args
		wantN   int
		wantErr bool
	}{
		{"Read", Item{"foo", nil}, args{nil}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.item.Read(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Item.Read() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestZipWriteItem_ReadAt(t *testing.T) {
	type args struct {
		b   []byte
		off int64
	}
	tests := []struct {
		name    string
		item    Item
		args    args
		wantN   int
		wantErr bool
	}{
		{"ReadAt", Item{"foo", nil}, args{nil, 0}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.item.ReadAt(tt.args.b, tt.args.off)
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.ReadAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Item.ReadAt() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestZipWriteItem_Seek(t *testing.T) {
	type args struct {
		offset int64
		whence int
	}
	tests := []struct {
		name    string
		item    Item
		args    args
		wantRet int64
		wantErr bool
	}{
		{"Seek", Item{"foo", nil}, args{0, io.SeekCurrent}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet, err := tt.item.Seek(tt.args.offset, tt.args.whence)
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.Seek() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRet != tt.wantRet {
				t.Errorf("Item.Seek() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestZipWriteItem_Stat(t *testing.T) {
	tests := []struct {
		name    string
		item    Item
		want    os.FileInfo
		wantErr bool
	}{
		{"Stat", Item{"foo", nil}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.item.Stat()
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.Stat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Item.Stat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipWriteItem_Sync(t *testing.T) {
	tests := []struct {
		name    string
		item    Item
		wantErr bool
	}{
		{"Sync", Item{"foo", nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.Sync(); (err != nil) != tt.wantErr {
				t.Errorf("Item.Sync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZipWriteItem_Truncate(t *testing.T) {
	type args struct {
		size int64
	}
	tests := []struct {
		name    string
		item    Item
		args    args
		wantErr bool
	}{
		{"Truncate", Item{"foo", nil}, args{0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.Truncate(tt.args.size); (err != nil) != tt.wantErr {
				t.Errorf("Item.Truncate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZipWriteItem_Close(t *testing.T) {
	tests := []struct {
		name    string
		item    Item
		wantErr bool
	}{
		{"Close", Item{"foo", nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.item.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Item.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestZipWriteItem_WriteAt(t *testing.T) {
	type args struct {
		b   []byte
		off int64
	}
	tests := []struct {
		name    string
		item    Item
		args    args
		wantN   int
		wantErr bool
	}{
		{"WriteAt", Item{"foo", nil}, args{nil, 0}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.item.WriteAt(tt.args.b, tt.args.off)
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.WriteAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Item.WriteAt() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestZipWriteItem_WriteString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		item    Item
		args    args
		wantN   int
		wantErr bool
	}{
		{"WriteString", Item{"foo", nil}, args{""}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := tt.item.WriteString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.WriteString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("Item.WriteString() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestZipWriteItem_Readdir(t *testing.T) {
	type args struct {
		count int
	}
	tests := []struct {
		name    string
		item    Item
		args    args
		want    []os.FileInfo
		wantErr bool
	}{
		{"ReadDir", Item{"foo", nil}, args{0}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.item.Readdir(tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.Readdir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Item.Readdir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZipWriteItem_Readdirnames(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		item    Item
		args    args
		want    []string
		wantErr bool
	}{
		{"Readdirnames", Item{"foo", nil}, args{0}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.item.Readdirnames(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Item.Readdirnames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Item.Readdirnames() = %v, want %v", got, tt.want)
			}
		})
	}
}
