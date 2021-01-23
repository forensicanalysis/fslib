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

package gpt

import (
	"encoding/binary"
	"io/fs"
	"os"
	"strings"
	"testing"
	"time"
	"unicode/utf16"

	"github.com/stretchr/testify/assert"

	"github.com/forensicanalysis/fslib/filesystem/fstests"
	"github.com/forensicanalysis/fslib/fsio"
)

func TestGPT(t *testing.T) {
	file, err := os.Open("../../test/data/filesystem/gpt_apfs.dd")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	gpt := GptPartitionTable{}
	err = gpt.Decode(file)
	if err != nil {
		t.Fatal(err)
	}

	primary := gpt.Primary()
	assert.EqualValues(t, []uint8{0x45, 0x46, 0x49, 0x20, 0x50, 0x41, 0x52, 0x54}, primary.Signature())
	partitions := primary.Entries()[0]
	name := partitions.Name()
	var u16 []uint16
	for i := 0; i < len(name); i += 2 {
		u16 = append(u16, binary.LittleEndian.Uint16(name[i:i+2]))
	}
	assert.EqualValues(t, "disk image", strings.Trim(string(utf16.Decode(u16)), "\x00"))
	assert.EqualValues(t, 40, partitions.FirstLba())
}

func Test_GPT(t *testing.T) {
	pptPathTests := map[string]*fstests.PathTest{
		// Root
		"Root": {
			TestName:    "Root",
			Path:        ".",
			InfoSize:    0,
			InfoMode:    os.ModeDir,
			InfoModTime: time.Time{},
			InfoIsDir:   true,
			FileName:    ".",
			// FileReaddir:      []os.FileInfo{},
			FileReaddirnames: []string{"p0"},
			Head:             []byte{},
		},
		// Partition
		"Partition 0": {
			TestName:    "Partition 0",
			Path:        "p0",
			InfoSize:    39024 * 512,
			InfoMode:    0,
			InfoModTime: time.Time{},
			InfoIsDir:   false,
			FileName:    "p0",
			// FileReaddir:      []os.FileInfo{},
			FileReaddirnames: []string{},
			Head:             []byte{0x3, 0x6, 0x6d, 0x2e, 0x74, 0x5, 0x3e, 0xea, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		},
	}

	fstests.RunTest(t, "GPT", "filesystem/gpt_apfs.dd", func(f fsio.ReadSeekerAt) (fs.FS, error) { return New(f) }, pptPathTests)
}

func BenchmarkGPT(b *testing.B) {
	for n := 0; n < b.N; n++ {
		file, _ := os.Open("../../test/data/filesystem/gpt_apfs.dd")
		gpt := GptPartitionTable{}
		err := gpt.Decode(file)
		if err != nil {
			b.Fatal(err)
		}
		file.Close()
	}
}
