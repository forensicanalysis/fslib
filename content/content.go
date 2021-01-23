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

// Package content provides functions to extractString text from different file formats
// like docx, xslx, pptx or pdf.
package content

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"

	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filetype"
	"github.com/forensicanalysis/fslib/fsio"
)

// Content returns the binary contents as a string.
func Content(r fsio.ReadSeekerAt) (content io.Reader, err error) {
	detectedType, err := filetype.DetectReader(r)
	if err != nil {
		return nil, err
	}

	switch detectedType {
	case filetype.Docx:
		return docxContent(r)
	case filetype.Pptx:
		return pptxContent(r)
	case filetype.Xlsx:
		return xlsxContent(r)
	case filetype.Pdf:
		return PDFContent(r)
	}

	_, err = r.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err = StringsReader(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func docxContent(r fsio.ReadSeekerAt) (io.Reader, error) {
	size, err := fsio.GetSize(r)
	if err != nil {
		return nil, err
	}
	fsys, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	f, err := fslib.Open(fsys, "/word/document.xml")
	return bytes.NewBufferString(xmlContent(f)), err
}

func xlsxContent(r fsio.ReadSeekerAt) (io.Reader, error) {
	// unzip -> xl/worksheets/sheetX.xml
	size, err := fsio.GetSize(r)
	if err != nil {
		return nil, err
	}
	fsys, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	s := &bytes.Buffer{}

	r, err = fslib.Open(fsys, "/xl/sharedStrings.xml")
	if err == nil {
		s.WriteString(xmlContent(r))
	}

	i := 1
	for {
		r, err := fslib.Open(fsys, fmt.Sprintf("/xl/worksheets/sheet%d.xml", i))
		if err != nil {
			break
		}
		s.WriteString(xmlContent(r))
		i++
	}
	return s, nil
}

func pptxContent(r fsio.ReadSeekerAt) (io.Reader, error) {
	// unzip -> ppt/slides/slideX.xml
	size, err := fsio.GetSize(r)
	if err != nil {
		return nil, err
	}
	fsys, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	s := &bytes.Buffer{}
	i := 1
	for {
		r, err := fslib.Open(fsys, fmt.Sprintf("/ppt/slides/slide%d.xml", i))
		if err != nil {
			break
		}
		s.WriteString(xmlContent(r))
		i++
	}
	return s, nil
}
