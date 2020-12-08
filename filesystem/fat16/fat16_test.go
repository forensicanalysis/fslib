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

package fat16

import (
	"encoding/binary"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem/fstests"
	"github.com/forensicanalysis/fslib/fsio"
)

func Test_FAT16Parser(t *testing.T) {
	file, err := os.Open("../../test/data/filesystem/fat16.dd")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	vh := volumeHeader{}
	err = binary.Read(file, binary.LittleEndian, &vh)
	assert.NoError(t, err)

	assert.EqualValues(t, "NO NAME", strings.TrimSpace(string(vh.VolumeLabel[:])))
	assert.EqualValues(t, uint8(0xf8), vh.MediaID)
}

func Test_FAT16(t *testing.T) {
	tests := fstests.GetDefaultContainerTests()

	anyTime := time.Time{}
	tests["rootTest"].InfoModTime = anyTime
	tests["rootTest"].InfoMode = os.ModeDir
	tests["rootTest"].InfoSize = 32 * 512

	tests["dir1Test"].InfoModTime = anyTime
	tests["dir1Test"].InfoMode = os.ModeDir

	tests["dir2Test"].InfoModTime = anyTime
	tests["dir2Test"].InfoMode = os.ModeDir

	tests["file1Test"].InfoModTime = anyTime
	tests["file1Test"].InfoMode = 0

	tests["file2Test"].InfoModTime = anyTime
	tests["file2Test"].InfoMode = 0

	delete(tests, "file2Test") // TODO: fix test
	delete(tests, "file3Test") // TODO: fix test
	delete(tests, "file4Test") // TODO: fix test

	fstests.RunTest(t, "FAT16", "filesystem/fat16.dd", func(f fsio.ReadSeekerAt) (fs.FS, error) { return New(f) }, tests)
}
