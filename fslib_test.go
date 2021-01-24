package fslib_test

import (
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/forensicanalysis/fslib/gpt"
	"github.com/forensicanalysis/fslib/mbr"
	"github.com/forensicanalysis/fslib/ntfs"
	"github.com/forensicanalysis/fslib/systemfs"

	"github.com/forensicanalysis/fslib/fat16"
)

func TestFSs(t *testing.T) {
	tests := []struct {
		name string
		fsys fs.FS
		path string
	}{
		// {"FAT16", newFAT(t), "image/alps.jpg"}, // TODO
		// {"GPT", newGPT(t), "p0"},
		// {"MBR", newMBR(t), "p0"},
		// {"NTFS", newNTFS(t), "image/alps.jpg"}, // TODO
		// {"OSFS", osfs.New(), "."}, // TODO
		// {"System FS", newSystemFS(t), "."}, // TODO
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fstest.TestFS(tt.fsys, tt.path); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func newFAT(t *testing.T) fs.FS {
	f, err := os.Open("testdata/data/filesystem/fat16.dd")
	if err != nil {
		t.Fatal(err)
	}

	fsys, err := fat16.New(f)
	if err != nil {
		t.Fatal(err)
	}
	return fsys
}

func newMBR(t *testing.T) fs.FS {
	f, err := os.Open("testdata/data/filesystem/mbr_fat16.dd")
	if err != nil {
		t.Fatal(err)
	}

	fsys, err := mbr.New(f)
	if err != nil {
		t.Fatal(err)
	}
	return fsys
}

func newGPT(t *testing.T) fs.FS {
	f, err := os.Open("testdata/data/filesystem/gpt_apfs.dd")
	if err != nil {
		t.Fatal(err)
	}

	fsys, err := gpt.New(f)
	if err != nil {
		t.Fatal(err)
	}
	return fsys
}

func newNTFS(t *testing.T) fs.FS {
	f, err := os.Open("testdata/data/filesystem/ntfs.dd")
	if err != nil {
		t.Fatal(err)
	}

	fsys, err := ntfs.New(f)
	if err != nil {
		t.Fatal(err)
	}
	return fsys
}

func newSystemFS(t *testing.T) fs.FS {
	fsys, err := systemfs.New()
	if err != nil {
		t.Fatal(err)
	}
	return fsys
}
