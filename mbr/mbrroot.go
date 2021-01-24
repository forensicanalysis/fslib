package mbr

import (
	"io/fs"
	"os"
	"syscall"
	"time"
)

// Root is a pseudo root directory for a Master Boot Record.
type Root struct {
	mbr *MbrPartitionTable
}

func (r *Root) Read([]byte) (int, error) {
	return 0, syscall.EPERM
}

// Name always returns / for MBR roots.
func (r *Root) Name() string { return "." }

// ReadDir lists all partitions in the MBR.
func (r *Root) ReadDir(count int) ([]fs.DirEntry, error) {
	var partitionInfos []fs.DirEntry
	partitions := r.mbr.Partitions()
	for index, partition := range partitions {
		if count != 0 && index == count {
			return partitionInfos, nil
		}
		if partition.NumSectors() != 0 {
			p := NewPartition(index, &partitions[index])
			partitionInfos = append(partitionInfos, p)
		}
	}
	return partitionInfos, nil
}

// Size returns 0 for MBR pseudo roots.
func (r *Root) Size() int64 { return 0 }

// Mode returns os.ModeDir for MBR pseudo roots.
func (r *Root) Mode() os.FileMode { return os.ModeDir }

// ModTime returns the zero time (0001-01-01 00:00) for MBR pseudo roots.
func (r *Root) ModTime() time.Time { return time.Time{} }

// IsDir returns true for MBR pseudo roots.
func (r *Root) IsDir() bool { return true }

// Sys returns nil for MBR pseudo roots.
func (r *Root) Sys() interface{} { return nil }

// Close does not do anything for MBR pseudo roots.
func (r *Root) Close() error { return nil }

// Stat returns the MBR pseudo roots itself as os.FileMode.
func (r *Root) Stat() (fs.FileInfo, error) { return r, nil }
