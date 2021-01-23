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
	"bytes"
	"errors"
	"fmt"
	"github.com/forensicanalysis/fslib/filesystem/buffer"
	"github.com/forensicanalysis/fslib/filesystem/zip"
	"io"
	"io/fs"
	"io/ioutil"
	"path"
	"strings"

	"github.com/forensicanalysis/fslib/filesystem/fat16"
	"github.com/forensicanalysis/fslib/filesystem/gpt"
	"github.com/forensicanalysis/fslib/filesystem/mbr"
	"github.com/forensicanalysis/fslib/filesystem/ntfs"
	"github.com/forensicanalysis/fslib/filesystem/osfs"
	"github.com/forensicanalysis/fslib/filetype"
	"github.com/forensicanalysis/fslib/fsio"
)

func parseRealPath(sample string) (rpath []element, err error) {
	parts := strings.Split(sample, "/")

	if len(parts) == 0 {
		return []element{{"", ""}}, nil
	}

	key := "."
	var fsys fs.FS = osfs.New()
	var fsName = "OsFs"
	var isFS bool
	for len(parts) > 0 {
		key = path.Join(key, parts[0])
		parts = parts[1:]
		info, err := fs.Stat(fsys, key)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			file, err := fsys.Open(key)
			if err != nil {
				return nil, err
			}

			rpath = append(rpath, element{fsName, key})
			isFS, fsName, err = detectFsFromFile(path.Ext(key), file)
			if err != nil {
				return nil, fmt.Errorf("error detection fsys %s: %w", key, err)
			}

			file, err = reopen(file, fsys, key)
			if err != nil {
				return nil, err
			}

			fsys, err = fsFromName(fsName, file)
			if err != nil {
				return nil, fmt.Errorf("could not get fsys from name %s: %w", fsName, err)
			}
			if !isFS && len(parts) > 0 {
				return nil, errors.New("could not resolve path")
			}

			key = "."
		} else if len(parts) == 0 {
			rpath = append(rpath, element{fsName, key})
		}
	}
	return rpath, nil
}

func reopen(file fs.File, fsys fs.FS, key string) (fs.File, error) {
	if seeker, ok := file.(io.Seeker); ok {
		_, err := seeker.Seek(0, 0)
		if err == nil {
			return file, nil
		}
	}

	_ = file.Close()

	return fsys.Open(key)
}

func detectFsFromFile(ext string, base io.Reader) (isFs bool, fs string, err error) {
	ext = strings.TrimLeft(ext, ".")

	t, err := filetype.DetectReaderByExtension(base, ext)
	if err != nil && err != io.EOF {
		return
	}

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

func fsFromName(name string, r io.Reader) (fsys fs.FS, err error) {
	readSeekerAt, ok := r.(fsio.ReadSeekerAt)
	if !ok {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		readSeekerAt = bytes.NewReader(b)
	}

	switch name {
	case "OsFs":
		fsys = osfs.New()
	case "ZIP":
		// size, err := fsio.GetSize(readSeekerAt)
		// if err != nil {
		// 	return nil, err
		// }
		// fsys, err = zip.NewReader(readSeekerAt, size)

		zipfs, err := zip.New(readSeekerAt)
		if err != nil {
			return nil, err
		}
		fsys = buffer.New(zipfs)
	case "FAT16":
		fsys, err = fat16.New(readSeekerAt)
	case "MBR":
		fsys, err = mbr.New(readSeekerAt)
	case "GPT":
		fsys, err = gpt.New(readSeekerAt)
	case "NTFS":
		fsys, err = ntfs.New(readSeekerAt)
	}
	return fsys, err
}
