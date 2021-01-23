package fslib

import (
	"io/fs"
	"os"
	"time"
)

func ReadDirFromNames(n int, readdirnames func(int) ([]string, error)) ([]fs.DirEntry, error) {
	names, err := readdirnames(n)
	if err != nil {
		return nil, err
	}
	return SimpleEntries(names), nil
}

func InfosToEntries(infos []os.FileInfo) (entries []fs.DirEntry) {
	for _, info := range infos {
		entries = append(entries, InfoToEntry(info))
	}
	return entries
}

func InfoToEntry(info os.FileInfo) fs.DirEntry {
	return &DirEntry{info}
}

type DirEntry struct {
	fs.FileInfo
}

func (e *DirEntry) Type() fs.FileMode {
	return e.FileInfo.Mode().Type()
}

func (e *DirEntry) Info() (fs.FileInfo, error) {
	return e.FileInfo, nil
}

type SimpleEntry struct {
	name string
}

func (e *SimpleEntry) Name() string               { return e.name }
func (e *SimpleEntry) IsDir() bool                { return true }
func (e *SimpleEntry) Type() fs.FileMode          { return fs.ModeDir }
func (e *SimpleEntry) Info() (fs.FileInfo, error) { return e, nil }
func (e *SimpleEntry) Size() int64                { return 0 }
func (e *SimpleEntry) Mode() fs.FileMode          { return fs.ModeDir }
func (e *SimpleEntry) ModTime() time.Time         { return time.Time{} }
func (e *SimpleEntry) Sys() interface{}           { return nil }

func SimpleEntries(names []string) (entries []fs.DirEntry) {
	for _, name := range names {
		entries = append(entries, &SimpleEntry{name: name})
	}
	return entries
}
