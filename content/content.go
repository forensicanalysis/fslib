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
	"bufio"
	"bytes"
	"fmt"
	"github.com/forensicanalysis/fslib/filesystem/buffer"
	"github.com/forensicanalysis/fslib/filesystem/zip"
	"github.com/forensicanalysis/fslib/filetype"
	"github.com/forensicanalysis/fslib/fsio"
	"io"
)

// Content returns the binary contents as a string.
func Content(r io.Reader) (content io.Reader, err error) {
	bufReader := bufio.NewReaderSize(r, 8192)

	detectedType, err := filetype.DetectReader(bufReader)
	if err != nil {
		return nil, err
	}

	switch detectedType {
	case filetype.Docx:
		return docxContent(bufReader)
	case filetype.Pptx:
		return pptxContent(bufReader)
	case filetype.Xlsx:
		return xlsxContent(bufReader)
	case filetype.Pdf:
		return PDFContent(bufReader)
	}

	buf := &bytes.Buffer{}
	if err = StringsReader(bufReader, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func docxContent(r io.Reader) (io.Reader, error) {
	fsys, err := bufferedZipFS(r)
	if err != nil {
		return nil, err
	}
	f, err := fsys.Open("word/document.xml")

	if err != nil {
		return nil, err
	}
	br, err := fsio.MaybeBufferedReader(f)
	if err != nil {
		return nil, err
	}
	return bytes.NewBufferString(xmlContent(br)), nil
}

func xlsxContent(r io.Reader) (io.Reader, error) {
	// unzip -> xl/worksheets/sheetX.xml
	fsys, err := bufferedZipFS(r)
	if err != nil {
		return nil, err
	}

	s := &bytes.Buffer{}
	f, err := fsys.Open("xl/sharedStrings.xml")
	if err == nil {
		br, err := fsio.MaybeBufferedReader(f)
		if err != nil {
			return nil, err
		}
		s.WriteString(xmlContent(br))
	}

	i := 1
	for {
		f, err := fsys.Open(fmt.Sprintf("xl/worksheets/sheet%d.xml", i))
		if err != nil {
			break
		}

		br, err := fsio.MaybeBufferedReader(f)
		if err != nil {
			return nil, err
		}
		s.WriteString(xmlContent(br))
		i++
	}
	return s, nil
}

func pptxContent(r io.Reader) (io.Reader, error) {
	// unzip -> ppt/slides/slideX.xml
	fsys, err := bufferedZipFS(r)
	if err != nil {
		return nil, err
	}
	s := &bytes.Buffer{}
	i := 1
	for {
		f, err := fsys.Open(fmt.Sprintf("ppt/slides/slide%d.xml", i))
		if err != nil {
			break
		}
		br, err := fsio.MaybeBufferedReader(f)
		if err != nil {
			return nil, err
		}
		s.WriteString(xmlContent(br))
		i++
	}
	return s, nil
}

func bufferedZipFS(r io.Reader) (*buffer.FS, error) {
	br, err := fsio.MaybeBufferedReader(r)
	if err != nil {
		return nil, err
	}
	zipfs, err := zip.New(br)
	if err != nil {
		return nil, err
	}
	fsys := buffer.New(zipfs)

	return fsys, nil
}
