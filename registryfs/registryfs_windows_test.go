// +build go1.8

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

package registryfs

import (
	"io/fs"
	"reflect"
	"sort"
	"testing"
)

func TestRegistryFS_Open(t *testing.T) {

	computerName := &Key{
		name: "ComputerName",
		path: "HKEY_LOCAL_MACHINE/System/CurrentControlSet/Control/ComputerName/ComputerName",
	}

	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		wantItem fs.File
		wantErr  bool
	}{
		{"Open ComputerName", args{"HKEY_LOCAL_MACHINE/System/CurrentControlSet/Control/ComputerName/ComputerName"}, computerName, false},
		// {"Open XYZ", fs, args{"/HKEY_LOCAL_MACHINE/System/CurrentControlSet/Control/ComputerName/ComputerName"}, forensicfs.Item{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := New()
			gotItem, err := fs.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FS.Open() error = '%#v', wantErr '%#v'", err, tt.wantErr)
				return
			}

			gotKey := gotItem.(*Key)
			gotKey.fs = nil
			gotKey.Key = computerName.Key
			if !reflect.DeepEqual(gotKey, tt.wantItem) {
				t.Errorf("FS.Open() = '%#v', want '%#v'", gotKey, tt.wantItem)
			}
		})
	}
}

func TestRegistryKey_Readdir(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		wantItems []string
		wantErr   bool
	}{
		{"Open root", args{"."}, []string{"HKEY_CLASSES_ROOT", "HKEY_CURRENT_CONFIG", "HKEY_CURRENT_USER", "HKEY_LOCAL_MACHINE", "HKEY_USERS"}, false},
		{"Open HKEY_LOCAL_MACHINE", args{"HKEY_LOCAL_MACHINE"}, []string{"BCD00000000", "HARDWARE", "SAM", "SECURITY", "SOFTWARE", "SYSTEM"}, false},
		// {"Open SYSTEM", args{"HKEY_LOCAL_MACHINE/SYSTEM"}, []string{"HARDWARE", "SAM", "SOFTWARE", "SYSTEM"}, false},
		{"Open ComputerName", args{"HKEY_LOCAL_MACHINE/System/CurrentControlSet/Control/ComputerName/ComputerName"}, nil, false},
		{"Open ComputerName Parent", args{"HKEY_LOCAL_MACHINE/System/CurrentControlSet/Control/ComputerName"}, []string{"ComputerName", "ActiveComputerName"}, false},
		{"Open CurrentControlSet", args{"HKEY_LOCAL_MACHINE/System/CurrentControlSet"}, []string{"Control", "Enum", "Hardware Profiles", "Policies", "Services", "Software"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := New()
			entries, err := fs.ReadDir(fsys, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Key.Readdir() error = '%#v', wantErr '%#v'", err, tt.wantErr)
				return
			}

			var filenames []string
			for _, entry := range entries {
				filenames = append(filenames, entry.Name())
			}
			sort.Strings(tt.wantItems)
			if !reflect.DeepEqual(filenames, tt.wantItems) {
				t.Errorf("Key.Readdir() = '%#v', want '%#v'", filenames, tt.wantItems)
			}
		})
	}
}
