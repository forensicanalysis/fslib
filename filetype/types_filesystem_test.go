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
	"log"
	"os"
	"testing"
)

func TestGPT(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"gpt", args{"filesystem/gpt_apfs.dd"}, true},
		{"mbr", args{"filesystem/mbr_fat16.dd"}, false},
	}
	for _, tt := range tests {
		file, err := os.Open("../test/data/" + tt.args.filename)
		if err != nil {
			t.Fatalf("Could not open file %s", tt.args.filename)
		}
		head := make([]byte, 8192)
		_, err = file.Read(head)
		if err != nil {
			t.Fatalf("Could not read file %s", tt.args.filename)
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := GPTMatch(head); got != tt.want {
				t.Errorf("GPTMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMBR(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"gpt", args{"filesystem/gpt_apfs.dd"}, true},
		{"mbr", args{"filesystem/mbr_fat16.dd"}, true},
	}
	for _, tt := range tests {
		file, err := os.Open("../test/data/" + tt.args.filename)
		if err != nil {
			t.Fatalf("Could not open file %s", tt.args.filename)
		}
		head := make([]byte, 8192)
		_, err = file.Read(head)
		if err != nil {
			t.Fatalf("Could not read file %s", tt.args.filename)
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := MBRMatch(head); got != tt.want {
				t.Errorf("MBRMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNTFSMatch(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"ntfs", args{"filesystem/ntfs.dd"}, true},
		{"mbr", args{"filesystem/mbr_fat16.dd"}, false},
	}
	for _, tt := range tests {
		file, err := os.Open("../test/data/" + tt.args.filename)
		if err != nil {
			t.Fatalf("Could not open file %s", tt.args.filename)
		}
		head := make([]byte, 16)
		_, err = file.Read(head)
		if err != nil {
			t.Fatalf("Could not read file %s", tt.args.filename)
		}

		t.Run(tt.name, func(t *testing.T) {
			if got := NTFSMatch(head); got != tt.want {
				log.Printf("%x", head)
				t.Errorf("NTFSMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
