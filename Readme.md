<h1 align="center">fslib</h1>
<h3 align="center">file system processing for forensics</h3>

<p  align="center">
 <a href="https://github.com/forensicanalysis/fslib/actions"><img src="https://github.com/forensicanalysis/fslib/workflows/CI/badge.svg" alt="build" /></a>
 <a href="https://codecov.io/gh/forensicanalysis/fslib"><img src="https://codecov.io/gh/forensicanalysis/fslib/branch/master/graph/badge.svg" alt="coverage" /></a>
 <a href="https://goreportcard.com/report/github.com/forensicanalysis/fslib"><img src="https://goreportcard.com/badge/github.com/forensicanalysis/fslib" alt="report" /></a>
 <a href="https://godoc.org/github.com/forensicanalysis/fslib"><img src="https://godoc.org/github.com/forensicanalysis/fslib?status.svg" alt="doc" /></a>
 <a href="https://app.fossa.io/projects/git%2Bgithub.com%2Fforensicanalysis%2Ffslib?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.io/api/projects/git%2Bgithub.com%2Fforensicanalysis%2Ffslib.svg?type=shield"/></a>
</p>


The fslib project contains a collection of tools and libraries to parse file
systems, archives and other data types. The included libraries can be used to
access disk images of with different partitioning and file systems.
Additionally, file systems for live access to the currently mounted file system
and registry (on Windows) are implemented.

Some of the file systems supported:

- Native OS file system 
- ZIP
- NTFS
- FAT16
- MBR
- GPT
- Windows Registry (live not from files)

Meta file systems:

- ⭐ **Recursive FS**: Access container files on file systems recursivly, e.g. `"ntfs.dd/container/Computer forensics - Wikipedia.zip/Computer forensics - Wikipedia.pdf"`
- Buffer FS: Buffer accessed files of an underlying file system
- System FS: Similar to the native OS file system, but falls back to NTFS on failing access on Windows

## Paths

Paths in fslib use [io/fs](https://tip.golang.org/pkg/io/fs/#ValidPath) conventions:

> - Path names passed to open are unrooted, slash-separated sequences of path elements, like “x/y/z”.
> - Path names must not contain a “.” or “..” or empty element, except for the special case that the root directory is named “.”.
> - Paths are slash-separated on all systems, even Windows.
> - Backslashes must not appear in path names.

## Installation

``` shell
go get -u github.com/forensicanalysis/fslib
```



## Examples

### NTFSReaddirnames
``` go
package main

import (
	"fmt"
	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem/ntfs"
	"os"
)

func main() {
	// Read the root directory on an NTFS disk image.

	// open the disk image
	image, _ := os.Open("testdata/data/filesystem/ntfs.dd")

	// parse the file system
	fs, _ := ntfs.New(image)

	// get handle for root
	root, _ := fs.Open(".")

	// get filenames
	filenames, _ := fslib.Readdirnames(root, 0)

	// print filenames
	fmt.Println(filenames)
}

```

### ReadFile
``` go
package main

import (
	"fmt"
	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem/recursivefs"
	"io/ioutil"
	"os"
	"path"
)

func main() {
	// Read the text file on an FAT16 disk image.

	// parse the file system
	fs := recursivefs.New()

	// create fslib path
	wd, _ := os.Getwd()
	fpath, _ := fslib.ToForensicPath(path.Join(wd, "testdata/data/filesystem/fat16.dd/README.md"))

	// get handle the README.md
	file, _ := fs.Open(fpath)

	// get content
	content, _ := ioutil.ReadAll(file)

	// print content
	fmt.Println(string(content[7:16]))
}

```






## fs command

The fs command line tool that has various subcommands which imitate unix commands
but for nested file system structures.

 - **fs cat**: Print files
 - **fs file**: Determine files types
 - **fs hashsum**: Print hashsums
 - **fs ls**: List directory contents
 - **fs stat**: Display file status
 - **fs tree**: List contents of directories in a tree-like format


#### Download

https://github.com/forensicanalysis/fslib/releases

#### Usage Examples

List all files in a zip file:
```
fs ls test.zip
```

Extract the Amcache.hve file from a NTFS image in a zip file:

```
fs cat case/evidence.zip/ntfs.dd/Windows/AppCompat/Programs/Amcache.hve > Amcache.hve
```

Hash all files in a zip file:
```
fs hashsum case/evidence.zip/*
```



## Subpackages

| Package | Description |
| --- | --- |
| **fallbackfs** | The fallbackfs project implements a meta filesystem that wraps a sequence of file systems. |
| **fat16** | The fat16 project provides an io/fs implementation of the FAT16 file systems. |
| **filetype** | The filetype project provides functions to identify file types of binary data. |
| **fs** | The fs project implements the fs command line tool that has various subcommands which imitate unix commands but for nested file system structures. |
| **fsio** | The fsio project provides IO interfaces and functions similar for file system operations. |
| **fstest** | The fstest project provides functions for testing implementations of the io/fs. |
| **glob** | The glob project provides a globing function for io/fs. |
| **gpt** | The gpt project provides an io/fs implementation of the GUID partition table (GPT). |
| **mbr** | The mbr project provides an io/fs implementation of the Master Boot Record (MBR) partition table. |
| **ntfs** | The ntfs project provides an io/fs implementation of the New Technology File System (NTFS). |
| **osfs** | The osfs project provides an io/fs implementation of the native OS file system. |
| **recursivefs** | The recursivefs project provides an io/fs implementation that can open paths in nested file systems recursively. |
| **registryfs** | The registryfs project provides an io/fs implementation to access the Windows Registry. |
| **systemfs** | The systemfs project provides an io/fs implementation that uses the osfs as default, while a ntfs for every partition as a fallback on Windows, on UNIX the behavior is the same as osfs. |
| **zipfs** | The zipfs project provides an io/fs implementation to access zip files. |


## Contact

For feedback, questions and discussions you can use the [Open Source DFIR Slack](https://github.com/open-source-dfir/slack).
