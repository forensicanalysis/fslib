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
	"io/fs"
	"os"

	"github.com/forensicanalysis/fslib/ntfs"
)

func Example() {
	// Read the root directory on an NTFS disk image.

	// open the disk image
	image, _ := os.Open("testdata/filesystem/ntfs.dd")

	// parse the file system
	fsys, _ := ntfs.New(image, 0, 0)

	// get filenames
	entries, _ := fs.ReadDir(fsys, ".")

	var filenames []string
	for _, entry := range entries {
		filenames = append(filenames, entry.Name())
	}

	// print filenames
	fmt.Println(filenames)
	// Output: [$AttrDef $BadClus $Bitmap $Boot $Extend $LogFile $MFT $MFTMirr $Secure $UpCase $Volume README.md container document evidence.json folder image]
}
