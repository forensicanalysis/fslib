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

package osfs

import (
	"reflect"
	"runtime"
	"sort"
	"testing"

	"github.com/forensicanalysis/fslib"
)

func getOSFS(t *testing.T) (*FS, *Item, *Item) {
	fs := New()
	f, err := fs.OpenSystemPath("../../test/data/document/Digital forensics.txt")
	if err != nil {
		t.Fatal("Error opening file: ", err)
	}
	dir, err := fs.OpenSystemPath("../../test/data/")
	if err != nil {
		t.Fatal("Error opening file: ", err)
	}
	return fs, f.(*Item), dir.(*Item)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *FS
	}{
		{"New", &FS{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*func TestOSFSItem_Attributes(t *testing.T) {
	_, f, dir := getOSFS(t)
	tests := []struct {
		name string
		item *Item
		want map[string]interface{}
	}{
		{"File Attributes", f, map[string]interface{}{"mode": os.FileMode(0666), "modified": nil}},
		{"Folder Attributes", dir, map[string]interface{}{"mode": os.FileMode(0777) | os.ModeDir, "modified": nil}},
	}
	for _, tt := range tests {
		if runtime.GOOS == "windows" {
			t.Run(tt.name, func(t *testing.T) {
				attrs := tt.item.Attributes()
				attrs["modified"] = nil
				if got := attrs; !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Attributes() = %v, want %v", got, tt.want)
				}
			})
		}

	}
}*/

/*func TestOSFSItem_IsDir(t *testing.T) {
	_, f, dir := getOSFS(t)
	tests := []struct {
		name    string
		item    *Item
		wantDir bool
	}{
		{"File IsDir", f, false},
		{"Folder IsDir", dir, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDir := tt.item.IsDir(); gotDir != tt.wantDir {
				t.Errorf("IsDir() = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}
*/
func TestOSFSItem_Readdirnames(t *testing.T) {
	_, f, dir := getOSFS(t)
	type args struct {
		n int
	}
	tests := []struct {
		name      string
		item      *Item
		args      args
		wantItems []string
		wantErr   bool
	}{
		{"File Readdirnames", f, args{0}, []string{}, true},
		{"Folder Readdirnames", dir, args{0}, []string{"container", "document", "filesystem", "image", "creation.md", "_meta"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItems, err := tt.item.Readdirnames(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Readdirnames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			sort.Strings(gotItems)
			sort.Strings(tt.wantItems)
			if !reflect.DeepEqual(gotItems, tt.wantItems) {
				t.Errorf("Readdirnames() gotItems = %v, want %v", gotItems, tt.wantItems)
			}
		})
	}
}

func TestWindowsRoot(t *testing.T) {
	root, err := New().Open("/")
	if err != nil {
		t.Error(err)
		return
	}

	type args struct {
		n int
	}
	tests := []struct {
		name      string
		item      fslib.Item
		args      args
		wantItems []string
		wantErr   bool
	}{
		{"Root Readdirnames", root, args{0}, []string{"C"}, false},
	}
	for _, tt := range tests {
		if runtime.GOOS == "windows" {
			t.Run(tt.name, func(t *testing.T) {
				gotItems, err := tt.item.Readdirnames(tt.args.n)
				if (err != nil) != tt.wantErr {
					t.Errorf("Readdirnames() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !isSubset(gotItems, tt.wantItems) {
					t.Errorf("Readdirnames() gotItems = %v, want %v", gotItems, tt.wantItems)
				}
			})
		}
	}
}

/*func TestOSFSItem_Size(t *testing.T) {
	_, f, dir := getOSFS(t)
	tests := []struct {
		name  string
		item  *Item
		wantS int64
	}{
		{"File Size", f, 678},
		{"Folder Size", dir, 4096},
	}
	for _, tt := range tests {
		if runtime.GOOS == "windows" {
			t.Run(tt.name, func(t *testing.T) {
				if gotS := tt.item.Size(); gotS != tt.wantS {
					t.Errorf("Size() = %v, want %v", gotS, tt.wantS)
				}
			})
		}
	}
}*/

func TestOSFS_Open(t *testing.T) {
	fs, _, _ := getOSFS(t)
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		fs       *FS
		args     args
		wantItem fslib.Item
		wantErr  bool
	}{
		{"Open fail", fs, args{"foo"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItem, err := fs.OpenSystemPath(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotItem, tt.wantItem) {
				t.Errorf("Open() gotItem = %v, want %v", gotItem, tt.wantItem)
			}
		})
	}
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
		{"Windows Abs Path", true, args{"C:\\Windows"}, "/C/Windows", false},
		// {"Windows Rel Path", true, args{"\\Windows"}, "/C/Windows", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if (tt.windowsTest && runtime.GOOS == "windows") || !tt.windowsTest {
				gotName, err := ToForensicPath(tt.args.systemPath)
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

func isSubset(s []string, sub []string) bool {
	for _, e := range sub {
		if !contains(s, e) {
			return false
		}
	}
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
