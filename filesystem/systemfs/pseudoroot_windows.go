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

package systemfs

import (
	"fmt"
	"github.com/forensicanalysis/fslib/filesystem/osfs"
	"github.com/forensicanalysis/go-vss"
	"os"
	"strings"
)

// Readdirnames lists all partitions in the window pseudo root.
func (r *Root) Readdirnames(n int) (partitions []string, err error) {
	root := osfs.Root{}
	partitions, err = root.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	for name := range r.fs.vss {
		partitions = append(partitions, name)
	}

	return partitions, nil
}

func getVSSStores(partition string) (map[string]*vss.VSS, error) {
	f, err := os.Open(fmt.Sprintf("\\\\.\\%s:", strings.ToLower(partition)))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stores, err := vss.NewByReader(f)
	if err != nil {
		return nil, nil
	}

	vssStores := map[string]*vss.VSS{}
	for i, store := range stores {
		vssStores[fmt.Sprintf("%svss%d", partition, i)] = store
	}
	return vssStores, nil
}
