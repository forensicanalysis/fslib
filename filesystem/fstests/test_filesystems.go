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

// Package fstests provides functions for testing implementations of the
// forensicfs.
package fstests

import (
	"github.com/forensicanalysis/fslib"
	"io/fs"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
	file3Head := []byte{0x89, 0x50, 0x4e, 0x47, 0xd, 0xa, 0x1a, 0xa, 0x0, 0x0, 0x0, 0xd, 0x49, 0x48, 0x44, 0x52}
	file4Head := []byte{0x4d, 0x4d, 0x0, 0x2a, 0x0, 0x4c, 0x26, 0x8, 0x9d, 0x97, 0x92, 0xff, 0x9c, 0x99, 0x92, 0xff}

	anyTime := time.Time{}
	dirTime := time.Date(2019, time.August, 15, 23, 01, 02, 0, time.UTC)
	dir2Time := time.Date(2019, time.August, 15, 23, 01, 02, 0, time.UTC)
	fileTime := time.Date(2018, time.March, 31, 21, 48, 36, 0, time.UTC)

	var anySys interface{}

	// Path Test
	rootTest := PathTest{"Root Test", ".", ".", 0, os.ModeDir | 0777, anyTime, true, anySys, rootFiles, []byte{}}
	dir1Test := PathTest{"Dir 1 Test", "document", "document", 0, os.ModeDir | 0755, dirTime, true, anySys, dirFiles, []byte{}}
	dir2Test := PathTest{"Dir 2 Test", "container", "container", 0, os.ModeDir | 0755, dir2Time, true, anySys, dir2Files, []byte{}}
	file1Test := PathTest{"File 1 Test", "README.md", "README.md", 496, 0644, fileTime, false, anySys, []string{}, file1Head}
	file2Test := PathTest{"File 2 Test", "image/alps.jpg", "alps.jpg", 344415, 0644, fileTime, false, anySys, []string{}, file2Head}
	file3Test := PathTest{"File 3 Test", "image/alps.png", "alps.png", 1338018, 0644, fileTime, false, anySys, []string{}, file3Head}
	file4Test := PathTest{"File 4 Test", "image/alps.tiff", "alps.tiff", 4994190, 0644, fileTime, false, anySys, []string{}, file4Head}

	return map[string]*PathTest{
		"rootTest":  &rootTest,
		"dir1Test":  &dir1Test,
		"dir2Test":  &dir2Test,
		"file1Test": &file1Test,
		"file2Test": &file2Test,
		"file3Test": &file3Test,
		"file4Test": &file4Test,
	}
}

// RunTest executes a set of tests.
func RunTest(t *testing.T, name, file string, new func(fsio.ReadSeekerAt) (fs.FS, error), tests map[string]*PathTest) {
	t.Run(name, func(t *testing.T) {
		fsys := osfs.New()
		base, err := fsys.OpenSystemPath("../../test/data/" + file)
		assert.NoError(t, err)
		assert.NotNil(t, base)
		if readSeekerAtBase, ok := base.(fsio.ReadSeekerAt); ok {
			checkFS(t, readSeekerAtBase, new, name, tests)
		} else {
			assert.Fail(t, "File does not implement ReadAt and Seek")
		}
	})
}

type Named interface {
	Name() string
}

func checkFS(t *testing.T, base fsio.ReadSeekerAt, new func(fsio.ReadSeekerAt) (fs.FS, error), name string, tests map[string]*PathTest) {
	// test creation
	fsys, err := new(base)
	assert.NoError(t, err)

	log.Print("check FS ", name)
	log.Print("-------------------")
	assert.NotNil(t, fsys)

	if namedFS, ok := fsys.(Named); ok {
		assert.EqualValues(t, name, namedFS.Name())
	} else {
		assert.Fail(t, "FS must have a name")
	}

	// test no leading slash
	// _, err = fsys.Open("no_slash")
	// assert.Error(t, err)

	// test not existing path
	_, err = fsys.Open("/non_existing")
	assert.Error(t, err)

	for _, tt := range tests {
		t.Run(tt.TestName, checkPath(name, tt, fsys))
	}
}

func checkPath(name string, tt *PathTest, fsys fs.FS) func(t *testing.T) {
	return func(t *testing.T) {
		log.Print("------------------------------")
		log.Print(name, " ", tt.TestName)
		log.Print("------------------------------")
		log.Print("test fsys.Stat")
		stat, err := fs.Stat(fsys, tt.Path)
		if assert.NoError(t, err) {
			assert.EqualValues(t, tt.InfoSize, stat.Size())
			assert.EqualValues(t, tt.InfoIsDir, stat.IsDir())
		}

		log.Print("-------------------")
		log.Print("test fsys.Open")
		file, err := fsys.Open(tt.Path)
		if assert.NoError(t, err) {
			log.Print("-------------------")
			log.Print("test item.Name")
			if namedFile, ok := file.(Named); ok {
				assert.EqualValues(t, tt.FileName, namedFile.Name())
			} else {
				assert.Fail(t, "File must have a name")
			}

			// fileInfos, err := file.Readdir(0)
			// assert.NoError(t, err)
			// assert.EqualValues(t, test.FileReaddir, fileInfos)
			if tt.InfoIsDir {
				log.Print("-------------------")
				log.Print("test dir.Readdirnames(0)")
				filenames, err := fslib.Readdirnames(file, 0)
				if assert.NoError(t, err) {
					assert.ElementsMatch(t, tt.FileReaddirnames, filenames, "dirnames do not match %s %s", tt.FileReaddirnames, filenames)
				}
			}

			if !tt.InfoIsDir {
				log.Print("-------------------")
				log.Print("test file.Read")
				head := make([]byte, len(tt.Head))
				c, err := file.Read(head)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.Head), c)
				assert.EqualValues(t, tt.Head, head)
			}
		}
	}
}

