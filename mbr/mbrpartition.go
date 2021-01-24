package mbr

import (
	"io"
	"io/fs"
	"strconv"
	"time"

	"github.com/forensicanalysis/fslib/fsio"
)

// Partition implements fs.File
type Partition struct {
	*io.SectionReader
	name      int
	partition *PartitionEntry
}

// NewPartition creates a new Partition object for parsing MBR partitions.
func NewPartition(name int, partition *PartitionEntry) *Partition {
	return &Partition{
		name:      name,
		partition: partition,
		SectionReader: io.NewSectionReader(
			&fsio.DecoderAtWrapper{ReadSeeker: partition.decoder},
			int64(partition.LbaStart()*512),
			int64(partition.NumSectors()*512),
		),
	}
}

// Name returns the name of a partition that consists of 'pX' where X is the
// number of the partition.
func (p *Partition) Name() string { return "p" + strconv.Itoa(p.name) }

// Size returns the partition size.
func (p *Partition) Size() int64 { return int64(p.partition.NumSectors() * 512) }

// Close does not do anything for MBR partitions.
func (p *Partition) Close() error { return nil }

// Stat return an fs.FileInfo object that describes a file.
func (p *Partition) Stat() (fs.FileInfo, error) { return p, nil }

// IsDir returns false for partition.
func (p *Partition) IsDir() bool { return false }

// Mode returns 0 for partitions.
func (p *Partition) Mode() fs.FileMode { return 0 }

// ModTime returns the zero time (0001-01-01 00:00) for partitions.
func (p *Partition) ModTime() time.Time { return time.Time{} }

// Sys returns the PartitionEntry.
func (p *Partition) Sys() interface{} { return p.partition }

func (p *Partition) Type() fs.FileMode { return p.Mode() }

func (p *Partition) Info() (fs.FileInfo, error) { return p, nil }
