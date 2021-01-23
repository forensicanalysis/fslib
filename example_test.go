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
	"os"

	"github.com/forensicanalysis/fslib/filesystem/ntfs"
)

func ExampleNTFSReaddirnames() {
	// Read the root directory on an NTFS disk image.

	// open the disk image
	image, _ := os.Open("test/data/filesystem/ntfs.dd")

	// parse the file system
	fs, _ := ntfs.New(image)

	// get handle for root
	root, _ := fs.Open(".")

	// get filenames
	filenames, _ := fslib.Readdirnames(root, 0)

	// print filenames
	fmt.Println(filenames)
	// Output: [$AttrDef $BadClus $BadClus:$Bad $Bitmap $Boot $Extend $LogFile $MFT $MFTMirr $Secure $Secure:$SDS $UpCase $UpCase:$Info $Volume README.md container document evidence.json folder image]
}
