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

package ntfs

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"reflect"
	"sort"
	"testing"
	"testing/fstest"
	"testing/iotest"
	"time"

	"github.com/forensicanalysis/fslib/fsio"
	fslibtest "github.com/forensicanalysis/fslib/fstest"
)

func TestFS(t *testing.T) {
	b, err := os.ReadFile("../testdata/filesystem/ntfs.dd")
	if err != nil {
		t.Error(err)
	}
	r := bytes.NewReader(b)

	fsys, err := New(r)
	if err != nil {
		t.Error(err)
	}

	if err := fstest.TestFS(fsys, "image/alps.jpg"); err != nil {
		t.Fatal(err)
	}

	f, err := fsys.Open("README.md")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	content := []byte("# :mag: evidence\nSample data for forensics processing\n\nForensics software need to be able to parse and process many different file formats. This repository contains samples of different file formats that can be used to test forensics software. Each file is accompanied by an entry in the [evidence.json](evidence.json) file with some metadata. \n\nExample entry for this README.md:\n```\n[\n  …\n  {\n    \"name\": \"README.md\", \n    \"mime\": \"text/plain\", \n    \"generator\": \"github.com\"\n  },\n  …\n]\n```\n")
	err = iotest.TestReader(f, content)
	if err != nil {
		t.Error(err)
	}
}

func Test_NTFSImage(t *testing.T) {
	tests := fslibtest.GetDefaultContainerTests()

	extra := []string{
		"$AttrDef", "$BadClus", "$Bitmap", "$Boot", "$Extend", "$LogFile", "$MFT", "$MFTMirr",
		"$Secure", "$UpCase", "$Volume",
	}

	tests["rootTest"].InfoModTime = time.Date(2019, time.August, 21, 17, 40, 04, 0, time.UTC)
	tests["rootTest"].InfoMode = fs.ModeDir
	tests["rootTest"].FileReaddirnames = append(tests["rootTest"].FileReaddirnames, extra...)
	sort.Strings(tests["rootTest"].FileReaddirnames)

	tests["dir1Test"].InfoModTime = time.Date(2019, time.August, 21, 17, 40, 04, 0, time.UTC)
	tests["dir1Test"].InfoMode = fs.ModeDir

	tests["dir2Test"].InfoModTime = time.Date(2019, time.August, 21, 17, 40, 04, 0, time.UTC)
	tests["dir2Test"].InfoMode = fs.ModeDir

	tests["file1Test"].InfoModTime = time.Date(2019, time.August, 21, 17, 40, 04, 0, time.UTC)
	tests["file1Test"].InfoMode = 0

	tests["file2Test"].InfoModTime = time.Date(2019, time.August, 21, 17, 40, 04, 0, time.UTC)
	tests["file2Test"].InfoMode = 0

	fslibtest.RunTest(t, "NTFS", "testdata/filesystem/ntfs.dd", func(f fsio.ReadSeekerAt) (fs.FS, error) { return New(f) }, tests)
}

func TestNew(t *testing.T) {
	r := bytes.NewReader([]byte{})
	type args struct {
		r io.ReaderAt
	}
	tests := []struct {
		name    string
		args    args
		wantFs  *FS
		wantErr bool
	}{
		{"no ntfs", args{r}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFs, err := New(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(gotFs, tt.wantFs) {
				t.Errorf("New() gotFs = %v, want %v", gotFs, tt.wantFs)
			}
		})
	}
}
