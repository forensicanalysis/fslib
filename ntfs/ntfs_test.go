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

package ntfs

import (
	"bytes"
	"io"
	"io/fs"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/forensicanalysis/fslib/fsio"
	"github.com/forensicanalysis/fslib/fstest"
)

func Test_NTFSImage(t *testing.T) {
	tests := fstest.GetDefaultContainerTests()

	extra := []string{
		"$AttrDef", "$BadClus", "$BadClus:$Bad", "$Bitmap", "$Boot", "$Extend", "$LogFile", "$MFT", "$MFTMirr",
		"$Secure", "$Secure:$SDS", "$UpCase", "$UpCase:$Info", "$Volume",
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

	fstest.RunTest(t, "NTFS", "filesystem/ntfs.dd", func(f fsio.ReadSeekerAt) (fs.FS, error) { return New(f) }, tests)
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
