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

// Package fstests provides functions for testing implementations of the
// forensicfs.
package fstests

import (
	fsys "io/fs"
	"os"
	"testing"
	"time"

	"github.com/forensicanalysis/fslib/filesystem/testfs"

	"github.com/forensicanalysis/fslib/fsio"
)

func TestGetDefaultContainerTests(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"Defaults", 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDefaultContainerTests(); len(got) != tt.want {
				t.Errorf("GetDefaultContainerTests() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestRunTest(t *testing.T) {
	fs := &testfs.FS{}
	fs.CreateFile("test.bar1", []byte("test"))
	fs.CreateFile("test.bar2", []byte("test"))
	fs.CreateFile("test.bar3", []byte("test"))
	n := func(f fsio.ReadSeekerAt) (fsys.FS, error) { return fs, nil }
	type args struct {
		t     *testing.T
		name  string
		file  string
		new   func(fsio.ReadSeekerAt) (fsys.FS, error)
		tests map[string]*PathTest
	}
	tests := []struct {
		name string
		args args
	}{
		{"RunTest Folder", args{t, "FS", "filesystem/ntfs.dd", n, map[string]*PathTest{
			"mem": {
				TestName:         "",                                              //string
				Path:             "/",                                             //string
				FileName:         ".",                                             //string
				InfoSize:         0,                                               //int64
				InfoMode:         os.ModeDir,                                      //os.FileMode
				InfoModTime:      time.Time{},                                     //time.Time
				InfoIsDir:        true,                                            //bool
				InfoSys:          nil,                                             //interface{}
				FileReaddirnames: []string{"test.bar1", "test.bar2", "test.bar3"}, //[]string
				Head:             []byte(""),                                      //[]byte
			}},
		}},
		{"RunTest File", args{t, "FS", "filesystem/ntfs.dd", n, map[string]*PathTest{
			"mem": {
				TestName:         "",             //string
				Path:             "/test.bar1",   //string
				FileName:         "test.bar1",    //string
				InfoSize:         4,              //int64
				InfoMode:         0,              //os.FileMode
				InfoModTime:      time.Time{},    //time.Time
				InfoIsDir:        false,          //bool
				InfoSys:          nil,            //interface{}
				FileReaddirnames: []string{},     //[]string
				Head:             []byte("test"), //[]byte
			}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunTest(tt.args.t, tt.args.name, tt.args.file, tt.args.new, tt.args.tests)
		})
	}
}
