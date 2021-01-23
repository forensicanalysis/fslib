// Copyright (c) 2020 Siemens AG
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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/forensicanalysis/fslib/fsio"
)

func TestPDFContent(t *testing.T) {
	pdf, err := ioutil.ReadFile("../test/data/document/Computer forensics - Wikipedia.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}

	type args struct {
		r fsio.ReadSeekerAt
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"size error", args{&brokenSeeker{}}, "", true},
		{"no pdf", args{bytes.NewReader([]byte{})}, "Computer forensics", true},
		{"mini pdf", args{bytes.NewReader(pdf)}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PDFContent(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("PDFContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			b, err := io.ReadAll(got)
			if err != nil {
				log.Fatal(err)
			}
			if !strings.Contains(string(b), tt.want) {
				t.Errorf("PDFContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
