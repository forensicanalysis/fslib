package gpt

import (
	"io/fs"
	"os"
	"syscall"
	"time"
)

// Root is a pseudo root directory containing the partitions.
type Root struct {
	gpt *GptPartitionTable
}

func (r *Root) Read([]byte) (int, error) {
	return 0, syscall.EPERM
}

// Name always returns '/' for GPT roots.
func (r *Root) Name() string { return "." }

// ReadDir lists all partitions in the GPT.
func (r *Root) ReadDir(count int) ([]fs.DirEntry, error) {
	var partitionInfos []fs.DirEntry
	partitions := r.gpt.Primary().Entries()
	for index, partition := range partitions {
		if count != 0 && index == count {
			return partitionInfos, nil
		}
		if partition.FirstLba() != 0 || partition.LastLba() != 0 {
			p := NewPartition(index, &partitions[index])
			partitionInfos = append(partitionInfos, p)
		}
	}
	return partitionInfos, nil
}

// Size returns 0 for GPT pseudo roots.
func (r *Root) Size() int64 { return 0 }

// Mode returns os.ModeDir for GPT pseudo roots.
func (r *Root) Mode() os.FileMode { return os.ModeDir }

// ModTime returns the zero time (0001-01-01 00:00) for GPT pseudo roots.
func (r *Root) ModTime() time.Time { return time.Time{} }

// IsDir returns true for GPT pseudo roots.
func (r *Root) IsDir() bool { return true }

// Sys returns nil for GPT pseudo roots.
func (r *Root) Sys() interface{} { return nil }

// Close does not do anything for GPT pseudo roots.
func (r *Root) Close() error { return nil }

// Stat returns the GPT pseudo roots itself as os.FileMode.
func (r *Root) Stat() (fs.FileInfo, error) { return r, nil }
