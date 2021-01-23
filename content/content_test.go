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

package content

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/forensicanalysis/fslib/filesystem/osfs"
	"github.com/forensicanalysis/fslib/fsio"
)

type read struct{
	done bool
}

func (b *read) Read([]byte) (n int, err error) {
	if b.done {
		return 0, io.EOF
	}
	b.done = true
	return 0, nil
}

type readAt struct{}

func (b *readAt) ReadAt([]byte, int64) (n int, err error) { return 0, nil }

type seek struct{}

func (b *seek) Seek(int64, int) (int64, error) { return 0, nil }

type brokenSeeker struct {
	fsio.ErrorSeeker
	readAt
	read
}

type brokenReader struct {
	fsio.ErrorReader
	readAt
	seek
}

func TestGetContent(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"document/Computer forensics - Wikipedia.pdf", args{"document/Computer forensics - Wikipedia.pdf"}, "Computer forensics", false},
		{"document/Design_of_the_FAT_file_system.xlsx", args{"document/Design_of_the_FAT_file_system.xlsx"}, "Design of the FAT file system", false},
		{"document/Digital forensics.docx", args{"document/Digital forensics.docx"}, "Digital forensics", false},
		{"document/Digital forensics.txt", args{"document/Digital forensics.txt"}, "Digital forensics", false},
		{"document/NTFS.pptx", args{"document/NTFS.pptx"}, "NTFS", false},
	}
	for _, tt := range tests {
		fs := osfs.New()
		file, err := fs.OpenSystemPath("../test/data/" + tt.args.filename)
		if err != nil {
			t.Fatalf("Could not open file %s", tt.args.filename)
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := Content(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Content() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			b, err := ioutil.ReadAll(got)
			if err != nil {
				t.Error(err)
			}

			wantParts := strings.Split(tt.want, " ")
			for _, wantPart := range wantParts {

				if !strings.Contains(string(b), wantPart) {
					t.Errorf("Content() %s does not contain %v", tt.args.filename, wantPart)
				}
			}
		})
	}
}

func TestContent(t *testing.T) {
	br := &brokenReader{}
	br.ErrorReader = fsio.ErrorReader{Skip: 1}

	type args struct {
		r fsio.ReadSeekerAt
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// {"error 1", args{&brokenSeeker{}}, true},
		{"error 2", args{&brokenReader{}}, true},
		{"error 3", args{br}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Content(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Content() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
