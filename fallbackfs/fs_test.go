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

package fallbackfs

import (
	"io/ioutil"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestFallbackFS_Open(t *testing.T) {
	mem1 := fstest.MapFS{}
	mem2 := fstest.MapFS{}

	mem1["foo"] = &fstest.MapFile{Data: []byte("fs1")}
	mem2["foo"] = &fstest.MapFile{Data: []byte("fs2")}
	mem2["bar"] = &fstest.MapFile{Data: []byte("bar2")}

	fallbackFS := New(mem1, mem2)

	type args struct {
		name string
	}
	tests := []struct {
		name     string
		fs       *FS
		args     args
		wantData []byte
		wantErr  bool
	}{
		{"Fallback Test 1", fallbackFS, args{"foo"}, []byte("fs1"), false},
		{"Fallback Test 2", fallbackFS, args{"bar"}, []byte("bar2"), false},
		{"Fallback Test failing", fallbackFS, args{"/bar"}, []byte("bar2"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItem, err := tt.fs.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FallbackFS.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				gotData, _ := ioutil.ReadAll(gotItem)

				if !reflect.DeepEqual(gotData, tt.wantData) {
					t.Errorf("FallbackFS.Open() = %v, want %v", gotData, tt.wantData)
				}
			}

			info, err := tt.fs.Stat(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FallbackFS.Stat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if info.IsDir() {
					t.Errorf("FallbackFS.IsDir() = %v, want %v", info.IsDir(), false)
				}
			}
		})
	}
}
