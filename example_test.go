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

package fslib_test

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/recursivefs"
)

func ExampleReadFile() {
	// Read the pdf header from a zip file on an NTFS disk image.

	// parse the file system
	fsys := recursivefs.New()

	// create fslib path
	wd, _ := os.Getwd()
	fpath, _ := fslib.ToForensicPath(path.Join(wd, "testdata/data/filesystem/ntfs.dd/container/Computer forensics - Wikipedia.zip/Computer forensics - Wikipedia.pdf"))

	// get handle the README.md
	file, err := fsys.Open(fpath)
	if err != nil {
		panic(err)
	}

	// get content
	content, _ := io.ReadAll(file)

	// print content
	fmt.Println(string(content[0:4]))
	// Output: %PDF
}
