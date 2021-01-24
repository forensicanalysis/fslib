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

package filetype

import (
	"bufio"
	"io"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/forensicanalysis/fslib/fsio"
)

func TestIdentify(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Filetype
		wantErr bool
	}{
		{"container/7z.7z", args{"container/7z.7z"}, Sevenz, false},
		{"container/Computer forensics - Wikipedia.pdf.gz", args{"container/Computer forensics - Wikipedia.pdf.gz"}, Gz, false},
		{"container/tar.tar", args{"container/tar.tar"}, Tar, false},
		{"container/zip.zip", args{"container/zip.zip"}, Zip, false},
		{"document/Computer forensics - Wikipedia.pdf", args{"document/Computer forensics - Wikipedia.pdf"}, Pdf, false},
		{"document/Design_of_the_FAT_file_system.xlsx", args{"document/Design_of_the_FAT_file_system.xlsx"}, Xlsx, false},
		{"document/Digital forensics.docx", args{"document/Digital forensics.docx"}, Docx, false},
		{"document/Digital forensics.txt", args{"document/Digital forensics.txt"}, Text, false},
		{"document/NTFS.pptx", args{"document/NTFS.pptx"}, Pptx, false},
		{"filesystem/exfat.dd", args{"filesystem/exfat.dd"}, ExFAT, false},
		{"filesystem/fat16.dd", args{"filesystem/fat16.dd"}, FAT16, false},
		{"filesystem/ntfs.dd", args{"filesystem/ntfs.dd"}, NTFS, false},
		{"filesystem/hfs+.dd", args{"filesystem/hfs+.dd"}, HFSPlus, false},
		{"filesystem/gpt_apfs.dd", args{"filesystem/gpt_apfs.dd"}, GPT, false},
		{"filesystem/mbr_fat16.dd", args{"filesystem/mbr_fat16.dd"}, MBR, false},
		{"image/alps.jpg", args{"image/alps.jpg"}, Jpeg, false},
		{"image/alps.png", args{"image/alps.png"}, Png, false},
		{"image/alps.tiff", args{"image/alps.tiff"}, Tiff, false},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open("../testdata/data/" + tt.args.filename)
			if err != nil {
				t.Fatalf("Could not open file %s", tt.args.filename)
			}
			defer file.Close()

			got, err := DetectReader(bufio.NewReaderSize(file, 8*1024))
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.ID != tt.want.ID {
				t.Errorf("DetectReader() = %v, want %v", got.ID, tt.want.ID)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open("../testdata/data/" + tt.args.filename)
			if err != nil {
				t.Fatalf("Could not open file %s", tt.args.filename)
			}
			defer file.Close()
			got, err := DetectReaderByExtension(file, path.Ext(tt.args.filename))
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectReaderByExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.ID != tt.want.ID {
				t.Errorf("DetectReaderByExtension() = %v, want %v", got.ID, tt.want.ID)
			}
		})

		t.Run(tt.name, func(t *testing.T) {

			file, err := os.Open("../testdata/data/" + tt.args.filename)
			if err != nil {
				t.Fatalf("Could not open file %s", tt.args.filename)
			}
			defer file.Close()
			head := make([]byte, 8192)
			_, err = file.Read(head)
			if err != nil {
				t.Fatalf("Could not read file %s", tt.args.filename)
			}
			got := DetectByExtension(head, path.Ext(tt.args.filename))
			if got.ID != tt.want.ID {
				t.Errorf("DetectByExtension() = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

func TestDetect(t *testing.T) {
	type args struct {
		buf []byte
	}
	tests := []struct {
		name string
		args args
		want *Filetype
	}{
		{"empty", args{nil}, Empty},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Detect(tt.args.buf); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectByExtension(t *testing.T) {
	type args struct {
		buf   []byte
		guess string
	}
	tests := []struct {
		name string
		args args
		want *Filetype
	}{
		{"empty", args{nil, ""}, Empty},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectByExtension(tt.args.buf, tt.args.guess); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectByExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectReader(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Filetype
		wantErr bool
	}{
		{"error reader", args{&fsio.ErrorReader{}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectReader(bufio.NewReader(tt.args.r))
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectReader() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectReaderByExtension(t *testing.T) {
	type args struct {
		r     io.Reader
		guess string
	}
	tests := []struct {
		name    string
		args    args
		want    *Filetype
		wantErr bool
	}{
		{"error reader", args{&fsio.ErrorReader{}, ""}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectReaderByExtension(tt.args.r, tt.args.guess)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectReaderByExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DetectReaderByExtension() got = %v, want %v", got, tt.want)
			}
		})
	}
}
