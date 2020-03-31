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

// Package copy provides copy functions for files and directories for afero
// (https://github.com/spf13/afero) filesystems.
package copy

import (
	"io"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// Item copies a file or directory recursively between two file systems.
func Item(srcfs, destfs afero.Fs, src, dest string) error {
	isDir, err := afero.IsDir(srcfs, src)
	if err != nil {
		return err
	}
	if isDir {
		return Directory(srcfs, destfs, src, dest)
	}
	return File(srcfs, destfs, src, dest)
}

// File copies a singe file between two file systems.
func File(srcfs, destfs afero.Fs, src, dest string) error {
	srcfile, err := srcfs.Open(src)
	if err != nil {
		return errors.Wrap(err, "open failed")
	}
	defer srcfile.Close()

	if err := destfs.MkdirAll(path.Dir(dest), 0700); err != nil {
		return errors.Wrap(err, "mkdir failed")
	}

	destfile, err := destfs.Create(dest)
	if err != nil {
		return errors.Wrap(err, "create failed")
	}
	defer destfile.Close()

	_, err = io.Copy(destfile, srcfile)
	return errors.Wrap(err, "copy failed")
}

// Directory copies a directory recursively between two file systems.
func Directory(srcfs, destfs afero.Fs, src, dest string) error {
	if err := destfs.MkdirAll(dest, 0700); err != nil {
		return err
	}

	children, err := afero.ReadDir(srcfs, src)
	if err != nil {
		return err
	}
	for _, child := range children {
		err := Item(srcfs, destfs, path.Join(src, child.Name()), path.Join(dest, child.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}
