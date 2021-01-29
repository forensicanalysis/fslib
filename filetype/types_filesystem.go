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

package filetype

import "github.com/h2non/filetype/types"

var filesystemTypes = []*Filetype{MBR, GPT, FAT16, HFSPlus, ExFAT, NTFS}

var (
	// MBR is the file type for the MBR partition system.
	MBR = &Filetype{"mbr", types.NewMIME("filesystem/mbr"), []string{"dd"}, MBRMatch, 1}
	// GPT is the file type for the GPT partition system.
	GPT = &Filetype{"gpt", types.NewMIME("filesystem/gpt"), []string{"dd"}, GPTMatch, 0}
	// FAT16 is the file type for the FAT16 file system.
	FAT16 = &Filetype{"fat16", types.NewMIME("filesystem/fat16"), []string{"dd"}, FAT16Match, 0}
	// HFSPlus is the file type for the HFSPlus file system.
	HFSPlus = &Filetype{"hfsplus", types.NewMIME("filesystem/hfsplus"), []string{"dd"}, HFSPlusMatch, 0}
	// ExFAT is the file type for the ExFAT file system.
	ExFAT = &Filetype{"exfat", types.NewMIME("filesystem/exfat"), []string{"dd"}, ExFATMatch, 0}
	// NTFS is the file type for the NTFS file system.
	NTFS = &Filetype{"ntfs", types.NewMIME("filesystem/ntfs"), []string{"dd"}, NTFSMatch, 0}
)

// MBRMatch checks if the buffer matches a signature for MBR file systems.
func MBRMatch(buf []byte) bool {
	return len(buf) >= 512 && buf[510] == 0x55 && buf[511] == 0xAA
}

// GPTMatch checks if the buffer matches a signature for GPT file systems.
func GPTMatch(buf []byte) bool {
	return len(buf) >= 512+9 && buf[512] == 'E' && buf[513] == 'F' && buf[514] == 'I' &&
		buf[515] == ' ' && buf[516] == 'P' && buf[517] == 'A' &&
		buf[518] == 'R' && buf[519] == 'T'
}

// FAT16Match checks if the buffer matches a signature for FAT16 file systems.
func FAT16Match(buf []byte) bool {
	return len(buf) >= 0x36+5 && buf[0x36] == 'F' && buf[0x36+1] == 'A' && buf[0x36+2] == 'T' &&
		buf[0x36+3] == '1' && buf[0x36+4] == '6'
}

// HFSPlusMatch checks if the buffer matches a signature for HFSPlus file systems.
func HFSPlusMatch(buf []byte) bool {
	return len(buf) >= 0x400+2 && buf[0x400] == 'H' && buf[0x401] == '+'
}

// ExFATMatch checks if the buffer matches a signature for ExFAT file systems.
func ExFATMatch(buf []byte) bool {
	return len(buf) >= 11 && buf[3] == 'E' && buf[4] == 'X' && buf[5] == 'F' && buf[6] == 'A' && buf[7] == 'T' &&
		buf[8] == ' ' && buf[9] == ' ' && buf[10] == ' '
}

// NTFSMatch checks if the buffer matches a signature for NTFS file systems.
func NTFSMatch(buf []byte) bool {
	return len(buf) >= 7 && buf[3] == 'N' && buf[4] == 'T' && buf[5] == 'F' && buf[6] == 'S'
}
