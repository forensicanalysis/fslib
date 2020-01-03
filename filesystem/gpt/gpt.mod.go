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

//nolint
// code generated at 2019-08-22T18:06:44Z

package gpt

import (
	"encoding/binary"
	"io"
)

type GptPartitionTable struct {
	decoder       io.ReadSeeker
	parent        interface{}
	root          interface{}
	sectorSize    int64 `ks:"sector_size,instance"`
	sectorSizeSet bool
	primary       *PartitionHeader `ks:"primary,instance"`
	primarySet    bool
	backup        *PartitionHeader `ks:"backup,instance"`
	backupSet     bool
}

func (k *GptPartitionTable) Parent() *GptPartitionTable {
	return k.parent.(*GptPartitionTable)
}
func (k *GptPartitionTable) Root() *GptPartitionTable {
	return k.root.(*GptPartitionTable)
}
func (k *GptPartitionTable) Decode(reader io.ReadSeeker, ancestors ...interface{}) (err error) {
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
	return
}
func (k *GptPartitionTable) SectorSize() (value int64) {
	if !k.sectorSizeSet {
		k.sectorSize = int64(0x200)
		k.sectorSizeSet = true
	}
	return k.sectorSize
}
func (k *GptPartitionTable) Primary() (value *PartitionHeader) {
	if !k.primarySet {
		var err error
		var elem PartitionHeader
		posprimary, _ := k.decoder.Seek(0, io.SeekCurrent) // Cannot fail
		_, err = k.decoder.Seek(int64(k.Root().SectorSize()), io.SeekStart)
		if err != nil {
			return
		}
		err = elem.Decode(k.decoder, k, k.Root())
		k.primary = &elem
		k.decoder.Seek(posprimary, io.SeekStart)
		k.primarySet = true
	}
	return k.primary
}
func (k *GptPartitionTable) Backup() (value *PartitionHeader) {
	if !k.backupSet {
		var err error
		var elem PartitionHeader
		posbackup, _ := k.decoder.Seek(0, io.SeekCurrent) // Cannot fail
		_, err = k.decoder.Seek(-int64(k.Root().SectorSize()), io.SeekStart)
		if err != nil {
			return
		}
		err = elem.Decode(k.decoder, k, k.Root())
		k.backup = &elem
		k.decoder.Seek(posbackup, io.SeekStart)
		k.backupSet = true
	}
	return k.backup
}

type PartitionEntry struct {
	decoder    io.ReadSeeker
	parent     interface{}
	root       interface{}
	typeGuid   []byte `ks:"type_guid,attribute"`
	guid       []byte `ks:"guid,attribute"`
	firstLba   uint64 `ks:"first_lba,attribute"`
	lastLba    uint64 `ks:"last_lba,attribute"`
	attributes uint64 `ks:"attributes,attribute"`
	name       []byte `ks:"name,attribute"`
}

func (k *PartitionEntry) Parent() *PartitionHeader {
	return k.parent.(*PartitionHeader)
}
func (k *PartitionEntry) Root() *GptPartitionTable {
	return k.root.(*GptPartitionTable)
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
		var elem []byte
		elem = make([]byte, 0x10)
		pos, _ := k.decoder.Seek(0, io.SeekCurrent)
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		pos = pos + int64(0x10)
		_, err = k.decoder.Seek(pos, io.SeekStart)
		k.typeGuid = elem
	}
	if err == nil {
		var elem []byte
		elem = make([]byte, 0x10)
		pos, _ := k.decoder.Seek(0, io.SeekCurrent)
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		pos = pos + int64(0x10)
		_, err = k.decoder.Seek(pos, io.SeekStart)
		k.guid = elem
	}
	if err == nil {
		var elem uint64
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.firstLba = elem
	}
	if err == nil {
		var elem uint64
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.lastLba = elem
	}
	if err == nil {
		var elem uint64
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.attributes = elem
	}
	if err == nil {
		var elem []byte
		elem = make([]byte, 0x48)
		pos, _ := k.decoder.Seek(0, io.SeekCurrent)
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		pos = pos + int64(0x48)
		_, err = k.decoder.Seek(pos, io.SeekStart)
		k.name = elem
	}
	return
}
func (k *PartitionEntry) TypeGuid() (value []byte) {
	return k.typeGuid
}
func (k *PartitionEntry) Guid() (value []byte) {
	return k.guid
}
func (k *PartitionEntry) FirstLba() (value uint64) {
	return k.firstLba
}
func (k *PartitionEntry) LastLba() (value uint64) {
	return k.lastLba
}
func (k *PartitionEntry) Attributes() (value uint64) {
	return k.attributes
}
func (k *PartitionEntry) Name() (value []byte) {
	return k.name
}

type PartitionHeader struct {
	decoder        io.ReadSeeker
	parent         interface{}
	root           interface{}
	signature      []byte           `ks:"signature,attribute"`
	revision       uint32           `ks:"revision,attribute"`
	headerSize     uint32           `ks:"header_size,attribute"`
	crc32Header    uint32           `ks:"crc32_header,attribute"`
	reserved       uint32           `ks:"reserved,attribute"`
	currentLba     uint64           `ks:"current_lba,attribute"`
	backupLba      uint64           `ks:"backup_lba,attribute"`
	firstUsableLba uint64           `ks:"first_usable_lba,attribute"`
	lastUsableLba  uint64           `ks:"last_usable_lba,attribute"`
	diskGuid       []byte           `ks:"disk_guid,attribute"`
	entriesStart   uint64           `ks:"entries_start,attribute"`
	entriesCount   uint32           `ks:"entries_count,attribute"`
	entriesSize    uint32           `ks:"entries_size,attribute"`
	crc32Array     uint32           `ks:"crc32_array,attribute"`
	entries        []PartitionEntry `ks:"entries,instance"`
	entriesSet     bool
}

func (k *PartitionHeader) Parent() *GptPartitionTable {
	return k.parent.(*GptPartitionTable)
}
func (k *PartitionHeader) Root() *GptPartitionTable {
	return k.root.(*GptPartitionTable)
}
func (k *PartitionHeader) Decode(reader io.ReadSeeker, ancestors ...interface{}) (err error) {
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
		elem = make([]byte, 8)
		pos, _ := k.decoder.Seek(0, io.SeekCurrent)
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		pos = pos + int64(8)
		_, err = k.decoder.Seek(pos, io.SeekStart)
		k.signature = elem
	}
	if err == nil {
		var elem uint32
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.revision = elem
	}
	if err == nil {
		var elem uint32
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.headerSize = elem
	}
	if err == nil {
		var elem uint32
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.crc32Header = elem
	}
	if err == nil {
		var elem uint32
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.reserved = elem
	}
	if err == nil {
		var elem uint64
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.currentLba = elem
	}
	if err == nil {
		var elem uint64
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.backupLba = elem
	}
	if err == nil {
		var elem uint64
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.firstUsableLba = elem
	}
	if err == nil {
		var elem uint64
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.lastUsableLba = elem
	}
	if err == nil {
		var elem []byte
		elem = make([]byte, 0x10)
		pos, _ := k.decoder.Seek(0, io.SeekCurrent)
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		pos = pos + int64(0x10)
		_, err = k.decoder.Seek(pos, io.SeekStart)
		k.diskGuid = elem
	}
	if err == nil {
		var elem uint64
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.entriesStart = elem
	}
	if err == nil {
		var elem uint32
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.entriesCount = elem
	}
	if err == nil {
		var elem uint32
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.entriesSize = elem
	}
	if err == nil {
		var elem uint32
		err = binary.Read(k.decoder, binary.LittleEndian, &elem)
		k.crc32Array = elem
	}
	return
}
func (k *PartitionHeader) Signature() (value []byte) {
	return k.signature
}
func (k *PartitionHeader) Revision() (value uint32) {
	return k.revision
}
func (k *PartitionHeader) HeaderSize() (value uint32) {
	return k.headerSize
}
func (k *PartitionHeader) Crc32Header() (value uint32) {
	return k.crc32Header
}
func (k *PartitionHeader) Reserved() (value uint32) {
	return k.reserved
}
func (k *PartitionHeader) CurrentLba() (value uint64) {
	return k.currentLba
}
func (k *PartitionHeader) BackupLba() (value uint64) {
	return k.backupLba
}
func (k *PartitionHeader) FirstUsableLba() (value uint64) {
	return k.firstUsableLba
}
func (k *PartitionHeader) LastUsableLba() (value uint64) {
	return k.lastUsableLba
}
func (k *PartitionHeader) DiskGuid() (value []byte) {
	return k.diskGuid
}
func (k *PartitionHeader) EntriesStart() (value uint64) {
	return k.entriesStart
}
func (k *PartitionHeader) EntriesCount() (value uint32) {
	return k.entriesCount
}
func (k *PartitionHeader) EntriesSize() (value uint32) {
	return k.entriesSize
}
func (k *PartitionHeader) Crc32Array() (value uint32) {
	return k.crc32Array
}
func (k *PartitionHeader) Entries() (value []PartitionEntry) {
	if !k.entriesSet {
		var err error
		var elem PartitionEntry
		posentries, _ := k.decoder.Seek(0, io.SeekCurrent) // Cannot fail
		_, err = k.decoder.Seek(int64(k.EntriesStart())*int64(k.Root().SectorSize()), io.SeekStart)
		if err != nil {
			return
		}
		k.entries = []PartitionEntry{}
		for index := 0; index < int(k.EntriesCount()); index++ {
			pos, _ := k.decoder.Seek(0, io.SeekCurrent)
			err = elem.Decode(k.decoder, k, k.Root())
			pos = pos + int64(k.EntriesSize())
			_, err = k.decoder.Seek(pos, io.SeekStart)
			if err != nil {
				if err == io.EOF {
					err = nil
				}
				break
			}
			k.entries = append(k.entries, elem)
		}
		k.decoder.Seek(posentries, io.SeekStart)
		k.entriesSet = true
	}
	return k.entries
}
