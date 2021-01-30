package fslib_test

import (
	"github.com/forensicanalysis/fslib"
	"io/fs"
	"os"
	"runtime"
	"testing"
	"testing/fstest"

	"github.com/forensicanalysis/fslib/fat16"
	"github.com/forensicanalysis/fslib/gpt"
	"github.com/forensicanalysis/fslib/mbr"
	"github.com/forensicanalysis/fslib/ntfs"
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
	f, err := os.Open("testdata/filesystem/fat16.dd")
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
	f, err := os.Open("testdata/filesystem/mbr_fat16.dd")
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
	f, err := os.Open("testdata/filesystem/gpt_apfs.dd")
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
	f, err := os.Open("testdata/filesystem/ntfs.dd")
	if err != nil {
		t.Fatal(err)
	}

	fsys, err := ntfs.New(f)
	if err != nil {
		t.Fatal(err)
	}
	return fsys
}


func TestToForensicPath(t *testing.T) {
	type args struct {
		systemPath string
	}
	tests := []struct {
		name        string
		windowsTest bool
		args        args
		wantName    string
		wantErr     bool
	}{
		{"Windows Abs Path", true, args{"C:\\Windows"}, "C/Windows", false},
		// {"Windows Rel Path", true, args{"\\Windows"}, "/C/Windows", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.windowsTest && runtime.GOOS == "windows") || !tt.windowsTest {
				gotName, err := fslib.ToForensicPath(tt.args.systemPath)
				if (err != nil) != tt.wantErr {
					t.Errorf("ToForensicPath() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if gotName != tt.wantName {
					t.Errorf("ToForensicPath() gotName = %v, want %v", gotName, tt.wantName)
				}
			}
		})
	}
}
