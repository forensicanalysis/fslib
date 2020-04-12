<h1 align="center">fslib</h1>

<p  align="center">
 <a href="https://github.com/forensicanalysis/fslib/actions"><img src="https://github.com/forensicanalysis/fslib/workflows/CI/badge.svg" alt="build" /></a>
 <a href="https://codecov.io/gh/forensicanalysis/fslib"><img src="https://codecov.io/gh/forensicanalysis/fslib/branch/master/graph/badge.svg" alt="coverage" /></a>
<a href="https://app.fossa.io/projects/git%2Bgithub.com%2Fforensicanalysis%2Ffslib?ref=badge_shield" alt="FOSSA Status"><img src="https://app.fossa.io/api/projects/git%2Bgithub.com%2Fforensicanalysis%2Ffslib.svg?type=shield"/></a>
 <a href="https://goreportcard.com/report/github.com/forensicanalysis/fslib"><img src="https://goreportcard.com/badge/github.com/forensicanalysis/fslib" alt="report" /></a>
 <a href="https://pkg.go.dev/github.com/forensicanalysis/fslib"><img src="https://img.shields.io/badge/go.dev-documentation-007d9c?logo=go&logoColor=white" alt="doc" /></a>
</p>


The fslib project contains a collection of tools and libraries to parse file
systems, archives and other data types. The included libraries can be used to
access disk images of with different partitioning and file systems.
Additionally, file systems for live access to the currently mounted file system
and registry (on Windows) are implemented.

## Paths
Paths in fslib use the following conventions:

- all paths are separated by forward slashes '/' (yes, even the windows registry)
- all paths need to start with forward slashes '/' (exception: the OSFS accepts relative paths)


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
	"github.com/forensicanalysis/fslib/filesystem/ntfs"
	"os"
)

func main() {
	// Read the root directory on an NTFS disk image.

	// open the disk image
	image, _ := os.Open("test/data/filesystem/ntfs.dd")

	// parse the file system
	fs, _ := ntfs.New(image)

	// get handle for root
	root, _ := fs.Open("/")

	// get filenames
	filenames, _ := root.Readdirnames(0)

	// print filenames
	fmt.Println(filenames)
}

```

### ReadFile
``` go
package main

import (
	"fmt"
	"github.com/forensicanalysis/fslib/filesystem/osfs"
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
	fpath, _ := osfs.ToForensicPath(path.Join(wd, "test/data/filesystem/fat16.dd/README.md"))

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
 - **fs strings**: Find the printable strings in an object, or other binary, file
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
| **content** | The content project provides functions to extractString text from different file formats like docx, xslx, pptx or pdf. |
| **copy** | The copy project provides copy functions for files and directories for afero (https://github.com/spf13/afero) filesystems. |
| **fallbackfs** | The fallbackfs project implements a meta filesystem that wraps a sequence of file systems. |
| **fat16** | The fat16 project provides a forensicfs implementation of the FAT16 file systems. |
| **filesystem** | The filesystem project contains implementations of the forensicfs interface for various file systems and similar data structures. |
| **filetype** | The filetype project provides functions to identify file types of binary data. |
| **forensicfs** | The forensicfs project provides defaults for read-only file system items. |
| **fs** | The fs project implements the fs command line tool that has various subcommands which imitate unix commands but for nested file system structures. |
| **fsio** | The fsio project provides IO interfaces and functions similar for file system operations. |
| **fstests** | The fstests project provides functions for testing implementations of the forensicfs. |
| **glob** | The glob project provides a globing function for forensicfs. |
| **gpt** | The gpt project provides a forensicfs implementation of the GUID partition table (GPT). |
| **mbr** | The mbr project provides a forensicfs implementation of the Master Boot Record (MBR) partition table. |
| **ntfs** | The ntfs project provides a forensicfs implementation of the New Technology File System (NTFS). |
| **osfs** | The osfs project provides a forensicfs implementation of the native OS file system. |
| **recursivefs** | The recursivefs project provides a forensicfs implementation that can open paths in nested forensicfs recursively. |
| **registryfs** | The registryfs project provides a forensicfs implementation to access the Windows Registry. |
| **systemfs** | The systemfs project provides a forensicfs implementation that uses the osfs as default, while a ntfs for every partition as a fallback on Windows, on UNIX the behavior is the same as osfs. |
| **testfs** | The testfs project provides a in memory forensicfs implementation for testing. |
| **zip** | The zip project provides a forensicfs implementation to access zip files. |
| **zipwrite** | The zipwrite project provides a write-only file system implementation for afero to create zip files. |


## Contact

For feedback, questions and discussions you can use the [Open Source DFIR Slack](https://github.com/open-source-dfir/slack).

## Acknowledgment

The development of this software was partially sponsored by Siemens CERT, but
is not an official Siemens product.


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fforensicanalysis%2Ffslib.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fforensicanalysis%2Ffslib?ref=badge_large)