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

package ntfs

import (
	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/fsio"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/forensicanalysis/fslib/filesystem/fstests"
)

func Test_NTFSImage(t *testing.T) {
	tests := fstests.GetDefaultContainerTests()

	extra := []string{
		"$AttrDef", "$BadClus", "$BadClus:$Bad", "$Bitmap", "$Boot", "$Extend", "$LogFile", "$MFT", "$MFTMirr",
		"$Secure", "$Secure:$SDS", "$UpCase", "$UpCase:$Info", "$Volume",
	}

	tests["rootTest"].InfoModTime = time.Date(2019, time.August, 21, 17, 40, 04, 0, time.UTC)
	tests["rootTest"].InfoMode = os.ModeDir
	tests["rootTest"].FileReaddirnames = append(tests["rootTest"].FileReaddirnames, extra...)
	sort.Strings(tests["rootTest"].FileReaddirnames)

	tests["dir1Test"].InfoModTime = time.Date(2019, time.August, 21, 17, 40, 04, 0, time.UTC)
	tests["dir1Test"].InfoMode = os.ModeDir

	tests["dir2Test"].InfoModTime = time.Date(2019, time.August, 21, 17, 40, 04, 0, time.UTC)
	tests["dir2Test"].InfoMode = os.ModeDir

	tests["file1Test"].InfoModTime = time.Date(2019, time.August, 21, 17, 40, 04, 0, time.UTC)
	tests["file1Test"].InfoMode = 0

	tests["file2Test"].InfoModTime = time.Date(2019, time.August, 21, 17, 40, 04, 0, time.UTC)
	tests["file2Test"].InfoMode = 0

	fstests.RunTest(t, "NTFS", "filesystem/ntfs.dd", func(f fsio.ReadSeekerAt) (fslib.FS, error) { return New(f) }, tests)
}
