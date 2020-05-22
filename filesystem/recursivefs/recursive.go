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

package recursivefs

import (
	"errors"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem/fat16"
	"github.com/forensicanalysis/fslib/filesystem/gpt"
	"github.com/forensicanalysis/fslib/filesystem/mbr"
	"github.com/forensicanalysis/fslib/filesystem/ntfs"
	"github.com/forensicanalysis/fslib/filesystem/osfs"
	"github.com/forensicanalysis/fslib/filesystem/zip"
	"github.com/forensicanalysis/fslib/filetype"
	"github.com/forensicanalysis/fslib/fsio"
)

func parseRealPath(sample string) (rpath []element, err error) {
	sample = path.Clean(sample)
	parts := strings.Split(sample, "/")

	if len(parts) == 0 {
		return []element{{"", ""}}, nil
	}

	key := "/"
	var fs fslib.FS = osfs.New()
	for len(parts) > 0 {
		key = path.Join(key, parts[0])
		parts = parts[1:]
		info, err := fs.Stat(key)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			file, err := fs.Open(key)
			if err != nil {
				return nil, err
			}

			rpath = append(rpath, element{fs.Name(), key})
			key = "/"
			isFs, fsName, err := detectFsFromFile(file)
			if err != nil {
				return nil, fmt.Errorf("error detection fs %s: %w", key, err)
			}
			fs, err = fsFromName(fsName, file)
			if err != nil {
				return nil, fmt.Errorf("could not get fs from name %s: %w", fsName, err)
			}
			if !isFs && len(parts) > 0 {
				return nil, errors.New("could not resolve path")
			}
		} else if len(parts) == 0 {
			rpath = append(rpath, element{fs.Name(), key})
		}
	}
	return rpath, nil
}

func detectFsFromFile(base fslib.Item) (isFs bool, fs string, err error) {
	ext := strings.TrimLeft(path.Ext(base.Name()), ".")

	t, err := filetype.DetectReaderByExtension(base, ext)
	if err != nil && err != io.EOF {
		return
	}
	base.Seek(0, 0) // nolint: errcheck

	switch t {
	case filetype.GPT:
		fs = "GPT"
	case filetype.MBR:
		fs = "MBR"
	case filetype.Zip, filetype.Xlsx, filetype.Pptx, filetype.Docx:
		fs = "ZIP"
	case filetype.FAT16:
		fs = "FAT16"
	case filetype.NTFS:
		fs = "NTFS"
	default:
		return false, "", nil
	}
	return true, fs, err
}

func fsFromName(name string, f fsio.ReadSeekerAt) (fs fslib.FS, err error) {
	switch name {
	case "OsFs":
		fs = osfs.New()
	case "ZIP":
		fs, err = zip.New(f)
	case "FAT16":
		fs, err = fat16.New(f)
	case "MBR":
		fs, err = mbr.New(f)
	case "GPT":
		fs, err = gpt.New(f)
	case "NTFS":
		fs, err = ntfs.New(f)
	}
	return
}
