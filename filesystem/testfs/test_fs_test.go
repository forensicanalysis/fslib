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

package testfs

import (
	"reflect"
	"testing"

	"github.com/forensicanalysis/fslib"
)

func getInFS() *FS {
	infs := &FS{}
	content := []byte("test")
	dirs := []string{"/dir/", "/dir/a/", "/dir/b/", "/dir/a/a/", "/dir/a/b/", "/dir/b/a/", "/dir/b/b/"}
	for _, dir := range dirs {
		infs.CreateDir(dir)
	}
	files := []string{"/foo.bin", "/dir/bar.bin", "/dir/baz.bin", "/dir/a/a/foo.bin", "/dir/a/b/foo.bin", "/dir/b/a/foo.bin", "/dir/b/b/foo.bin"}
	for _, file := range files {
		infs.CreateFile(file, content)
	}
	return infs
}

func TestTestFS_CreateDir(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		{"CreateDir", args{"/dir"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := getInFS()
			fs.CreateDir(tt.args.name)
		})
	}
}

func TestTestFS_CreateFile(t *testing.T) {
	type args struct {
		name string
		data []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{"Create File", args{"foo", []byte("test")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := getInFS()
			fs.CreateFile(tt.args.name, tt.args.data)
		})
	}
}

func TestTestFS_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"Name", "FS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := getInFS()
			if got := fs.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTestFS_Open(t *testing.T) {
	fs := getInFS()
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    fslib.Item
		wantErr bool
	}{
		{"Open", args{"/dir"}, &Directory{path: "dir", fs: fs}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := fs.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Open() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTestFolder_Readdirnames(t *testing.T) {
	type args struct {
		in0 int
	}
	tests := []struct {
		name      string
		dir       string
		args      args
		wantItems []string
		wantErr   bool
	}{
		{"Readdirnames", "/dir", args{0}, []string{"a", "b", "bar.bin", "baz.bin"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := getInFS()
			d, _ := fs.Open(tt.dir)
			gotItems, err := d.Readdirnames(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("Readdirnames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotItems, tt.wantItems) {
				t.Errorf("Readdirnames() gotItems = %v, want %v", gotItems, tt.wantItems)
			}
		})
	}
}
