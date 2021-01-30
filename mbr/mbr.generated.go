// Copyright (c) 2019-2020 Siemens AG
// Copyright (c) 2019-2021 Jonas Plum
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

//nolint
// code generated at 2019-08-22T16:57:11Z

package mbr

import (
	"encoding/binary"
	"io"
	"os"
)

/* MBR (Master Boot Record) partition table is a traditional way of
MS-DOS to partition larger hard disc drives into distinct
partitions.

This table is stored in the end of the boot sector (first sector) of
the drive, after the bootstrap code. Original DOS 2.0 specification
allowed only 4 partitions per disc, but DOS 3.2 introduced concept
of "extended partitions", which work as nested extra "boot records"
which are pointed to by original ("primary") partitions in MBR.*/
type MbrPartitionTable struct {
	decoder       io.ReadSeeker
	parent        interface{}
	root          interface{}
	bootstrapCode []byte            `ks:"bootstrap_code,attribute"`
	partitions    [4]PartitionEntry `ks:"partitions,attribute"`
	bootSignature []byte            `ks:"boot_signature,attribute"`
}

func (k *MbrPartitionTable) Parent() *MbrPartitionTable {
	return k.parent.(*MbrPartitionTable)
}
func (k *MbrPartitionTable) Root() *MbrPartitionTable {
	return k.root.(*MbrPartitionTable)
}
func (k *MbrPartitionTable) Decode(reader io.ReadSeeker, ancestors ...interface{}) (err error) {
	if k.decoder == nil {
		if reader == nil {
			panic("Neither decoder nor reader are set.")
		}
		k.decoder = reader
	}
	if len(ancestors) == 2 {
		k.parent = ancestors[0]
		k.root = ancestors[1]
	} else if len(ancestors) == 0 {
		k.parent = k
		k.root = k
	} else {
		panic("To many ancestors are given.")
	}
	if err == nil {
		var elem []byte
		elem = make([]byte, 0x1be)
		pos, _ := k.decoder.Seek(0, os.SEEK_CUR)
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		pos = pos + int64(0x1be)
		_, err = k.decoder.Seek(pos, os.SEEK_SET)
		k.bootstrapCode = elem
	}
	if err == nil {
		var elem PartitionEntry
		for index := 0; index < int(4); index++ {
			err = elem.Decode(k.decoder, k, k.Root())
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				break
			}
			k.partitions[index] = elem
		}
	}
	if err == nil {
		var elem []byte
		elem = make([]byte, 2)
		pos, _ := k.decoder.Seek(0, os.SEEK_CUR)
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		pos = pos + int64(2)
		_, err = k.decoder.Seek(pos, os.SEEK_SET)
		k.bootSignature = elem
	}
	return
}
func (k *MbrPartitionTable) BootstrapCode() (value []byte) {
	return k.bootstrapCode
}
func (k *MbrPartitionTable) Partitions() (value [4]PartitionEntry) {
	return k.partitions
}
func (k *MbrPartitionTable) BootSignature() (value []byte) {
	return k.bootSignature
}

type Chs struct {
	decoder     io.ReadSeeker
	parent      interface{}
	root        interface{}
	head        uint8 `ks:"head,attribute"`
	b2          uint8 `ks:"b2,attribute"`
	b3          uint8 `ks:"b3,attribute"`
	sector      int64 `ks:"sector,instance"`
	sectorSet   bool
	cylinder    int64 `ks:"cylinder,instance"`
	cylinderSet bool
}

func (k *Chs) Parent() *PartitionEntry {
	return k.parent.(*PartitionEntry)
}
func (k *Chs) Root() *MbrPartitionTable {
	return k.root.(*MbrPartitionTable)
}
func (k *Chs) Decode(reader io.ReadSeeker, ancestors ...interface{}) (err error) {
	if k.decoder == nil {
		if reader == nil {
			panic("Neither decoder nor reader are set.")
		}
		k.decoder = reader
	}
	if len(ancestors) == 2 {
		k.parent = ancestors[0]
		k.root = ancestors[1]
	} else if len(ancestors) == 0 {
		k.parent = k
		k.root = k
	} else {
		panic("To many ancestors are given.")
	}
	if err == nil {
		var elem uint8
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.head = elem
	}
	if err == nil {
		var elem uint8
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.b2 = elem
	}
	if err == nil {
		var elem uint8
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.b3 = elem
	}
	return
}
func (k *Chs) Head() (value uint8) {
	return k.head
}
func (k *Chs) B2() (value uint8) {
	return k.b2
}
func (k *Chs) B3() (value uint8) {
	return k.b3
}
func (k *Chs) Sector() (value int64) {
	if !k.sectorSet {
		k.sector = int64(k.B2() & 63)
		k.sectorSet = true
	}
	return k.sector
}
func (k *Chs) Cylinder() (value int64) {
	if !k.cylinderSet {
		k.cylinder = int64(k.B3() + ((k.B2() & 192) << 2))
		k.cylinderSet = true
	}
	return k.cylinder
}

type PartitionEntry struct {
	decoder       io.ReadSeeker
	parent        interface{}
	root          interface{}
	status        uint8  `ks:"status,attribute"`
	chsStart      *Chs   `ks:"chs_start,attribute"`
	partitionType uint8  `ks:"partition_type,attribute"`
	chsEnd        *Chs   `ks:"chs_end,attribute"`
	lbaStart      uint32 `ks:"lba_start,attribute"`
	numSectors    uint32 `ks:"num_sectors,attribute"`
}

func (k *PartitionEntry) Parent() *MbrPartitionTable {
	return k.parent.(*MbrPartitionTable)
}
func (k *PartitionEntry) Root() *MbrPartitionTable {
	return k.root.(*MbrPartitionTable)
}
func (k *PartitionEntry) Decode(reader io.ReadSeeker, ancestors ...interface{}) (err error) {
	if k.decoder == nil {
		if reader == nil {
			panic("Neither decoder nor reader are set.")
		}
		k.decoder = reader
	}
	if len(ancestors) == 2 {
		k.parent = ancestors[0]
		k.root = ancestors[1]
	} else if len(ancestors) == 0 {
		k.parent = k
		k.root = k
	} else {
		panic("To many ancestors are given.")
	}
	if err == nil {
		var elem uint8
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.status = elem
	}
	if err == nil {
		var elem Chs
		err = elem.Decode(k.decoder, k, k.Root())
		k.chsStart = &elem
	}
	if err == nil {
		var elem uint8
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.partitionType = elem
	}
	if err == nil {
		var elem Chs
		err = elem.Decode(k.decoder, k, k.Root())
		k.chsEnd = &elem
	}
	if err == nil {
		var elem uint32
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.lbaStart = elem
	}
	if err == nil {
		var elem uint32
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.numSectors = elem
	}
	return
}
func (k *PartitionEntry) Status() (value uint8) {
	return k.status
}
func (k *PartitionEntry) ChsStart() (value *Chs) {
	return k.chsStart
}
func (k *PartitionEntry) PartitionType() (value uint8) {
	return k.partitionType
}
func (k *PartitionEntry) ChsEnd() (value *Chs) {
	return k.chsEnd
}
func (k *PartitionEntry) LbaStart() (value uint32) {
	return k.lbaStart
}
func (k *PartitionEntry) NumSectors() (value uint32) {
	return k.numSectors
}
