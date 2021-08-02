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

package fat16

import (
	"encoding/binary"
	"io"
	"io/fs"
	"os"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/forensicanalysis/fslib/fsio"
	fslibtest "github.com/forensicanalysis/fslib/fstest"
)

func Test_FS(t *testing.T) {
	file, err := os.Open("../testdata/filesystem/fat16.dd")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fsys, err := New(file)
	if err != nil {
		t.Fatal(err)
	}
	err = fstest.TestFS(fsys, "image/alps.jpg")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Walk(t *testing.T) {
	file, err := os.Open("../testdata/filesystem/mbr_fat16.dd")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	_, err = file.Seek(128*512, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}

	r := io.NewSectionReader(file, 128*512, 34816*512)

	fsys, err := New(r)
	if err != nil {
		t.Fatal(err)
	}

	err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if strings.Contains(path, "System Volume Information/System Volume Information") {
			t.Fatal("recursion")
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_FAT16Parser(t *testing.T) {
	file, err := os.Open("../testdata/filesystem/fat16.dd")
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
	tests := fslibtest.GetDefaultContainerTests()

	anyTime := time.Time{}
	tests["rootTest"].InfoModTime = anyTime
	tests["rootTest"].InfoMode = fs.ModeDir
	tests["rootTest"].InfoSize = 32 * 512

	tests["dir1Test"].InfoModTime = anyTime
	tests["dir1Test"].InfoMode = fs.ModeDir

	tests["dir2Test"].InfoModTime = anyTime
	tests["dir2Test"].InfoMode = fs.ModeDir

	tests["file1Test"].InfoModTime = anyTime
	tests["file1Test"].InfoMode = 0

	tests["file2Test"].InfoModTime = anyTime
	tests["file2Test"].InfoMode = 0

	fslibtest.RunTest(t, "FAT16", "testdata/filesystem/fat16.dd", func(f fsio.ReadSeekerAt) (fs.FS, error) { return New(f) }, tests)
}
