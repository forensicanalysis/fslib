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

// Package fstests provides functions for testing implementations of the
// forensicfs.
package fstests

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem/osfs"
	"github.com/forensicanalysis/fslib/fsio"
)

// PathTest is a single test for file systems.
type PathTest struct {
	TestName         string
	Path             string
	FileName         string
	InfoSize         int64
	InfoMode         os.FileMode
	InfoModTime      time.Time
	InfoIsDir        bool
	InfoSys          interface{}
	FileReaddirnames []string
	Head             []byte
}

// GetDefaultContainerTests returns a set of default tests for the test data.
func GetDefaultContainerTests() map[string]*PathTest {
	// Test Data
	rootFiles := []string{"container", "document", "evidence.json", "image", "README.md", "folder"}
	dirFiles := []string{"Computer forensics - Wikipedia.pdf", "NTFS.pptx", "Design_of_the_FAT_file_system.xlsx", "Digital forensics.docx", "Digital forensics.txt"}
	dir2Files := []string{"Computer forensics - Wikipedia.7z", "Computer forensics - Wikipedia.tar", "Computer forensics - Wikipedia.pdf.gz", "Computer forensics - Wikipedia.zip"}
	file1Head := []byte{0x23, 0x20, 0x3a, 0x6d, 0x61, 0x67, 0x3a, 0x20, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65}
	file2Head := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01, 0x00, 0x00, 0x48}

	anyTime := time.Time{}
	dirTime := time.Date(2019, time.August, 15, 23, 01, 02, 0, time.UTC)
	dir2Time := time.Date(2019, time.August, 15, 23, 01, 02, 0, time.UTC)
	fileTime := time.Date(2018, time.March, 31, 21, 48, 36, 0, time.UTC)

	var anySys interface{}

	// Path Test
	rootTest := PathTest{"Root Test", "/", "/", 0, os.ModeDir | 0777, anyTime, true, anySys, rootFiles, []byte{}}
	dir1Test := PathTest{"Dir 1 Test", "/document", "document", 0, os.ModeDir | 0755, dirTime, true, anySys, dirFiles, []byte{}}
	dir2Test := PathTest{"Dir 2 Test", "/container", "container", 0, os.ModeDir | 0755, dir2Time, true, anySys, dir2Files, []byte{}}
	file1Test := PathTest{"File 1 Test", "/README.md", "README.md", 496, 0644, fileTime, false, anySys, []string{}, file1Head}
	file2Test := PathTest{"File 2 Test", "/image/alps.jpg", "alps.jpg", 344415, 0644, fileTime, false, anySys, []string{}, file2Head}

	return map[string]*PathTest{
		"rootTest":  &rootTest,
		"dir1Test":  &dir1Test,
		"dir2Test":  &dir2Test,
		"file1Test": &file1Test,
		"file2Test": &file2Test,
	}
}

// RunTest executes a set of tests.
func RunTest(t *testing.T, name, file string, new func(fsio.ReadSeekerAt) (fslib.FS, error), tests map[string]*PathTest) {
	t.Run(name, func(t *testing.T) {
		fs := osfs.New()
		base, err := fs.OpenSystemPath("../../test/data/" + file)
		assert.NoError(t, err)
		assert.NotNil(t, base)
		checkFS(t, base, new, name, tests)
	})
}

func checkFS(t *testing.T, base fsio.ReadSeekerAt, new func(fsio.ReadSeekerAt) (fslib.FS, error), name string, tests map[string]*PathTest) {
	// test creation
	fs, err := new(base)
	assert.NoError(t, err)

	log.Print("check FS ", name)
	log.Print("-------------------")
	assert.NotNil(t, fs)
	assert.EqualValues(t, name, fs.Name())

	// test no leading slash
	// _, err = fs.Open("no_slash")
	// assert.Error(t, err)

	// test not existing path
	_, err = fs.Open("/non_existing")
	assert.Error(t, err)

	for _, test := range tests {
		t.Run(test.TestName, func(t *testing.T) {
			log.Print("------------------------------")
			log.Print(name, " ", test.TestName)
			log.Print("------------------------------")
			log.Print("test fs.Stat")
			stat, err := fs.Stat(test.Path)
			if assert.NoError(t, err) {
				assert.EqualValues(t, test.InfoSize, stat.Size())
				assert.EqualValues(t, test.InfoIsDir, stat.IsDir())
			}

			log.Print("-------------------")
			log.Print("test fs.Open")
			file, err := fs.Open(test.Path)
			if assert.NoError(t, err) {
				log.Print("-------------------")
				log.Print("test item.Name")
				assert.EqualValues(t, test.FileName, file.Name())
				// fileInfos, err := file.Readdir(0)
				// assert.NoError(t, err)
				// assert.EqualValues(t, test.FileReaddir, fileInfos)
				if test.InfoIsDir {
					log.Print("-------------------")
					log.Print("test dir.Readdirnames(0)")
					filenames, err := file.Readdirnames(0)
					if assert.NoError(t, err) {
						assert.ElementsMatch(t, test.FileReaddirnames, filenames, "dirnames do not match %s %s", test.FileReaddirnames, filenames)
					}

					min := func(a, b int) int {
						if a < b {
							return a
						}
						return b
					}
					for _, i := range []int{1, 3, 1000} {
						log.Print("-------------------")
						log.Printf("test dir.Readdirnames(%d)", i)
						filenames, err = file.Readdirnames(i)
						if assert.NoError(t, err) {
							assert.Equal(t, min(i, len(test.FileReaddirnames)), len(filenames), "dirnames do not match %s %s", test.FileReaddirnames, filenames)
						}
					}
				}

				if !test.InfoIsDir {
					log.Print("-------------------")
					log.Print("test file.Read")
					head := make([]byte, len(test.Head))
					_, err = file.Read(head)
					assert.NoError(t, err)
					assert.EqualValues(t, test.Head, head)
				}
			}
		})
	}
}
