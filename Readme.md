<h1 align="center">fslib</h1>
<h3 align="center">file system processing for forensics</h3>

<p  align="center">
 <a href="https://godocs.io/github.com/forensicanalysis/fslib"><img src="https://godocs.io/github.com/forensicanalysis/fslib?status.svg" alt="doc" /></a>
</p>


The fslib project contains a collection of packages to parse file
systems, archives and similar data. The included packages can be used to
access disk images of with different partitioning and file systems.
Additionally, file systems for live access to the currently mounted file system
and registry (on Windows) are implemented.

All filesystems implement [io/fs](https://golang.org/pkg/io/fs/#FS).

### Included File systems

- **Native OS file system** (directory listing for Windows root provides list of drives)
- **Windows Registry** (live not from files)
- **NTFS**
- **FAT16**
- **MBR**
- **GPT**

### Meta file systems

- **Buffer FS**: Buffer accessed files of an underlying file system
- **System FS**: Similar to the native OS file system, but falls back to NTFS on failing access on Windows

### See also

- **[zipfs](http://github.com/forensicanalysis/zipfs)**: A zip file system
- ⭐ **[Recursive FS](http://github.com/forensicanalysis/recursivefs)**: Access container files on file systems recursively, e.g. `"ntfs.dd/forensic.zip/Computer forensics - Wikipedia.pdf"`

## Installation

``` shell
go get -u github.com/forensicanalysis/fslib
```

## Example

``` go
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

## Contact

For feedback, questions and discussions you can use the [Open Source DFIR Slack](https://github.com/open-source-dfir/slack).
