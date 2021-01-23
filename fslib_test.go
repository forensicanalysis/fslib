package fslib_test

import (
	"io/fs"
	"testing"

	"github.com/forensicanalysis/fslib/filesystem/fat16"
	"github.com/forensicanalysis/fslib/filesystem/gpt"
	"github.com/forensicanalysis/fslib/filesystem/mbr"
	"github.com/forensicanalysis/fslib/filesystem/ntfs"
	"github.com/forensicanalysis/fslib/filesystem/recursivefs"
	"github.com/forensicanalysis/fslib/filesystem/systemfs"
	"github.com/forensicanalysis/fslib/filesystem/zip"
)

func TestFiles(t *testing.T) {
	var _ fs.ReadDirFile = &fat16.Item{}
	var _ fs.ReadDirFile = &gpt.Root{}
	var _ fs.ReadDirFile = &mbr.Root{}
	var _ fs.ReadDirFile = &ntfs.Item{}
	var _ fs.ReadDirFile = &recursivefs.Item{}
	// var _ fs.ReadDirFile = &registryfs.Root{}
	// var _ fs.ReadDirFile = &registryfs.Key{}
	var _ fs.ReadDirFile = &systemfs.Root{}
	var _ fs.ReadDirFile = &zip.File{}
}
