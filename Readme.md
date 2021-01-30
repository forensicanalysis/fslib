<h1 align="center">fslib</h1>
<h3 align="center">file system processing for forensics</h3>

<p  align="center">
 <a href="https://github.com/forensicanalysis/fslib/actions"><img src="https://github.com/forensicanalysis/fslib/workflows/CI/badge.svg" alt="build" /></a>
 <a href="https://codecov.io/gh/forensicanalysis/fslib"><img src="https://codecov.io/gh/forensicanalysis/fslib/branch/master/graph/badge.svg" alt="coverage" /></a>
 <a href="https://goreportcard.com/report/github.com/forensicanalysis/fslib"><img src="https://goreportcard.com/badge/github.com/forensicanalysis/fslib" alt="report" /></a>
 <a href="https://godocs.io/github.com/forensicanalysis/fslib"><img src="https://godocs.io/github.com/forensicanalysis/fslib?status.svg" alt="doc" /></a>
 <a href="https://app.fossa.io/projects/git%2Bgithub.com%2Fforensicanalysis%2Ffslib?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.io/api/projects/git%2Bgithub.com%2Fforensicanalysis%2Ffslib.svg?type=shield"/></a>
</p>


The fslib project contains a collection of packages to parse file
systems, archives and other data types. The included packages can be used to
access disk images of with different partitioning and file systems.
Additionally, file systems for live access to the currently mounted file system
and registry (on Windows) are implemented.

## File systems supported

- **Native OS file system** (directory listing for Windows root provides list of drives)
- **Windows Registry** (live not from files)
- **NTFS**
- **FAT16**
- **MBR**
- **GPT**

## Meta file systems

- **Buffer FS**: Buffer accessed files of an underlying file system
- **System FS**: Similar to the native OS file system, but falls back to NTFS on failing access on Windows

## More file systems

- [zipfs](http://github.com/forensicanalysis/zipfs)
- ⭐ **[Recursive FS](http://github.com/forensicanalysis/recursivefs)**: Access container files on file systems recursively, e.g. `"ntfs.dd/forensic.zip/Computer forensics - Wikipedia.pdf"`

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
	"github.com/forensicanalysis/fslib/ntfs"
	"io/fs"
	"os"
)

func main() {
	// Read the root directory on an NTFS disk image.

	// open the disk image
	image, _ := os.Open("filesystem/ntfs.dd")

	// parse the file system
	fsys, _ := ntfs.New(image)

	// get filenames
	entries, _ := fs.ReadDir(fsys, ".")

	var filenames []string
	for _, entry := range entries {
		filenames = append(filenames, entry.Name())
	}

	// print filenames
	fmt.Println(filenames)
}

```


## Subpackages

| Package | Description |
| --- | --- |
| **fallbackfs** | The fallbackfs project implements a meta filesystem that wraps a sequence of file systems. |
| **fat16** | The fat16 project provides an io/fs implementation of the FAT16 file systems. |
| **filetype** | The filetype project provides functions to identify file types of binary data. |
| **fsio** | The fsio project provides IO interfaces and functions similar for file system operations. |
| **fstest** | The fstest project provides functions for testing implementations of the io/fs. |
| **glob** | The glob project provides a globing function for io/fs. |
| **gpt** | The gpt project provides an io/fs implementation of the GUID partition table (GPT). |
| **mbr** | The mbr project provides an io/fs implementation of the Master Boot Record (MBR) partition table. |
| **ntfs** | The ntfs project provides an io/fs implementation of the New Technology File System (NTFS). |
| **osfs** | The osfs project provides an io/fs implementation of the native OS file system. |
| **registryfs** | The registryfs project provides an io/fs implementation to access the Windows Registry. |
| **systemfs** | The systemfs project provides an io/fs implementation that uses the osfs as default, while a ntfs for every partition as a fallback on Windows, on UNIX the behavior is the same as osfs. |


## Contact

For feedback, questions and discussions you can use the [Open Source DFIR Slack](https://github.com/open-source-dfir/slack).
