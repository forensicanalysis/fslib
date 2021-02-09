package mbr

import (
	"io"
	"io/fs"
	"syscall"
	"time"
)

// Root is a pseudo root directory for a Master Boot Record.
type Root struct {
	mbr       *MbrPartitionTable
	dirOffset int
}

func (r *Root) Read([]byte) (int, error) {
	return 0, syscall.EPERM
}

// Name always returns / for MBR roots.
func (r *Root) Name() string { return "." }

// ReadDir lists all partitions in the MBR.
func (r *Root) ReadDir(n int) ([]fs.DirEntry, error) {
	var partitionInfos []fs.DirEntry
	partitions := r.mbr.Partitions()
	for index, partition := range partitions {
		if partition.NumSectors() != 0 {
			p := NewPartition(index, &partitions[index])
			partitionInfos = append(partitionInfos, p)
		}
	}

	// directory already exhausted
	if n <= 0 && r.dirOffset >= len(partitionInfos) {
		return nil, nil
	}

	var err error
	// read till end
	if n > 0 && r.dirOffset+n > len(partitionInfos) {
		err = io.EOF
		if r.dirOffset > len(partitionInfos) {
			return nil, err
		}
	}

	if n > 0 && r.dirOffset+n <= len(partitionInfos) {
		partitionInfos = partitionInfos[r.dirOffset : r.dirOffset+n]
		r.dirOffset += n
	} else {
		partitionInfos = partitionInfos[r.dirOffset:]
		r.dirOffset += len(partitionInfos)
	}

	return partitionInfos, err
}

// Size returns 0 for MBR pseudo roots.
func (r *Root) Size() int64 { return 0 }

// Mode returns fs.ModeDir for MBR pseudo roots.
func (r *Root) Mode() fs.FileMode { return fs.ModeDir }

// ModTime returns the zero time (0001-01-01 00:00) for MBR pseudo roots.
func (r *Root) ModTime() time.Time { return time.Time{} }

// IsDir returns true for MBR pseudo roots.
func (r *Root) IsDir() bool { return true }

// Sys returns nil for MBR pseudo roots.
func (r *Root) Sys() interface{} { return nil }

// Close does not do anything for MBR pseudo roots.
func (r *Root) Close() error { return nil }

// Stat returns the MBR pseudo roots itself as fs.FileMode.
func (r *Root) Stat() (fs.FileInfo, error) { return r, nil }
