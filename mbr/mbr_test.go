// +build go1.7

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

package mbr

import (
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/forensicanalysis/fslib/fsio"
	"github.com/forensicanalysis/fslib/fstest"
)

func TestMBREvidence(t *testing.T) {
	var err error

	file, err := os.Open("../testdata/data/filesystem/mbr_fat16.dd")
	assert.NoError(t, err)
	defer file.Close()

	mbr := MbrPartitionTable{}
	err = mbr.Decode(file)
	assert.NoError(t, err)

	p0 := mbr.Partitions()[0]
	assert.EqualValues(t, 128, p0.LbaStart())
	assert.EqualValues(t, 34816, p0.NumSectors())
	assert.EqualValues(t, 14, p0.PartitionType())
}

func Test_MBR(t *testing.T) {
	mbrPathTests := map[string]*fstest.PathTest{
		// Root
		"Root": {
			TestName:    "Root",
			Path:        ".",
			InfoSize:    0,
			InfoMode:    os.ModeDir,
			InfoModTime: time.Time{},
			InfoIsDir:   true,
			FileName:    ".",
			// FileReaddir:      []fs.FileInfo{},
			FileReaddirnames: []string{"p0"},
			Head:             []byte{},
		},
		// Partition
		"Partition 0": {
			TestName:    "Partition 0",
			Path:        "p0",
			InfoSize:    34816 * 512,
			InfoMode:    0,
			InfoModTime: time.Time{},
			InfoIsDir:   false,
			FileName:    "p0",
			// FileReaddir:      []fs.FileInfo{},
			FileReaddirnames: []string{},
			Head:             []byte{0xeb, 0x3c, 0x90, 0x4d, 0x53, 0x44, 0x4f, 0x53, 0x35, 0x2e, 0x30, 0x00, 0x02, 0x01, 0x02, 0x00},
		},
	}

	fstest.RunTest(t, "MBR", "filesystem/mbr_fat16.dd", func(f fsio.ReadSeekerAt) (fs.FS, error) { return New(f) }, mbrPathTests)
}

func BenchmarkMBR(b *testing.B) {
	for n := 0; n < b.N; n++ {
		file, _ := os.Open("../testdata/data/filesystem/mbr_fat16.dd")
		mbr := MbrPartitionTable{}
		err := mbr.Decode(file)
		if err != nil {
			b.Fatal(err)
		}
		file.Close()
	}
}
