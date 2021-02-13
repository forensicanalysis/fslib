// +build go1.16

package osfs

import (
	"github.com/forensicanalysis/fslib"
	"io/fs"
	"sort"
)

func (i *Item) ReadDir(n int) ([]fs.DirEntry, error) {
	entries, err := i.File.ReadDir(n)
	sort.Sort(fslib.ByName(entries))
	return entries, err
}
