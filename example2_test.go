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

package fslib_test

import (
	"fmt"
	"github.com/forensicanalysis/fslib"
	"io/ioutil"
	"os"
	"path"

	"github.com/forensicanalysis/fslib/filesystem/recursivefs"
)

func ExampleReadFile() {
	// Read the text file on an FAT16 disk image.

	// parse the file system
	fs := recursivefs.New()

	// create fslib path
	wd, _ := os.Getwd()
	fpath, _ := fslib.ToForensicPath(path.Join(wd, "test/data/filesystem/fat16.dd/README.md"))

	// get handle the README.md
	file, _ := fs.Open(fpath)

	// get content
	content, _ := ioutil.ReadAll(file)

	// print content
	fmt.Println(string(content[7:16]))
	// Output: evidence
}
