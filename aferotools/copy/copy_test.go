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

package copy

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/spf13/afero"
)

func generateTestFS(t *testing.T) (afero.Fs, string) {
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), dir)

	err = fs.MkdirAll("dir/subdir", 07777)
	if err != nil {
		t.Error("error creating file")
	}
	f, err := fs.Create("foo.txt")
	if err != nil {
		t.Error("error creating file foo.txt")
	}
	f, err = fs.Create("dir/subdir/subfoo.txt")
	if err != nil {
		t.Error("error creating file subfoo.txt")
	}
	f, err = fs.Create("dir/subdir/dir")
	if err != nil {
		t.Error("error creating file dir")
	}
	f.Write([]byte("test"))

	return fs, dir
}

func TestItem(t *testing.T) {
	fs, temp := generateTestFS(t)
	defer os.RemoveAll(temp)

	type args struct {
		srcfs  afero.Fs
		destfs afero.Fs
		src    string
		dest   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Copy file", args{fs, fs, "foo.txt", "bar2.txt"}, false},
		{"Copy dir", args{fs, fs, "dir/subdir", "bar"}, false},
		{"Copy not existing item", args{fs, fs, "bar.txt", "bar.txt"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Item(tt.args.srcfs, tt.args.destfs, tt.args.src, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("Item() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile(t *testing.T) {
	fs, temp := generateTestFS(t)
	defer os.RemoveAll(temp)

	type args struct {
		srcfs  afero.Fs
		destfs afero.Fs
		src    string
		dest   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Copy not existing file", args{fs, fs, "bar.txt", "bar2.txt"}, true},
		{"Copy to existing folder", args{fs, fs, "foo.txt", "dir"}, true},
		{"Copy to existing file", args{fs, fs, "foo.txt", "foo.txt/bar.txt"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := File(tt.args.srcfs, tt.args.destfs, tt.args.src, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("File() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirectory(t *testing.T) {
	fs, temp := generateTestFS(t)
	defer os.RemoveAll(temp)

	type args struct {
		srcfs  afero.Fs
		destfs afero.Fs
		src    string
		dest   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Copy not existing file", args{fs, fs, "dir3", "dir2"}, true},
		{"Copy to existing file", args{fs, fs, "dir", "foo.txt/bar.txt"}, true},
		{"Copy to existing file", args{fs, fs, "dir/subdir", ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Directory(tt.args.srcfs, tt.args.destfs, tt.args.src, tt.args.dest); (err != nil) != tt.wantErr {
				t.Errorf("Directory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
