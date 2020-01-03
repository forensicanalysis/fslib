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

// Package fs implements the fs command line tool that has various subcommands
// which imitate unix commands but for nested file system structures.
//     cat      Print files
//     file     Determine files types
//     hashsum  Print hashsums
//     ls       List directory contents
//     stat     Display file status
//     strings  Find the printable strings in an object, or other binary, file
//     tree     List contents of directories in a tree-like format
//
// Usage Examples
//
// Extract the Amcache.hve file from a NTFS image in a zip file:
//     fs cat case/evidence.zip/ntfs.dd/Windows/AppCompat/Programs/Amcache.hve > Amcache.hve
// Hash all files in a zip file:
//     fs hashsum case/evidence.zip/*
package main

import "github.com/forensicanalysis/fslib/cmd/fs/subcommands"

func main() {
	subcommands.Execute()
}
