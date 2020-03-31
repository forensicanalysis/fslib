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

package hfsplus

// Information taken from: https://developer.apple.com/legacy/library/technotes/tn/tn1150.html

/*
// HFSUniStr255 defines file and folder names on HFS Plus which consist of up
// to 255 Unicode characters with a preceding 16-bit length.
type HFSUniStr255 struct {
	Length  uint16
	UniChar [255]uint16
}

// For each file and folder, HFS Plus maintains a record containing access
// permissions, defined by the HFSPlusBSDInfo structure.
type HFSPlusBSDInfo struct {
	OwnerID    uint32
	GroupID    uint32
	AdminFlags uint8
	OwnerFlags uint8
	FileMode   uint16
	Special    uint32
}

// HFS Plus maintains information about the contents of a file using the HFSPlusForkData
// structure.
type HFSPlusForkData struct {
	LogicalSize uint64
	ClumpSize   uint32
	TotalBlocks uint32
	Extents     HFSPlusExtentRecord
}

// A HFSPlusExtentRecord contains descriptors for eight extents.
type HFSPlusExtentRecord [8]HFSPlusExtentDescriptor

// The HFSPlusExtentDescriptor structure is used to hold information about a specific
// extent.
type HFSPlusExtentDescriptor struct {
	StartBlock uint32
	BlockCount uint32
}

// The HFSPlusVolumeHeader contains information about the volume as a whole,
// including the location of other key structures in the volume.
type HFSPlusVolumeHeader struct {
	Signature          uint16
	Version            uint16
	Attributes         uint32
	LastMountedVersion uint32
	JournalInfoBlock   uint32
	CreateDate         uint32
	ModifyDate         uint32
	BackupDate         uint32
	CheckedDate        uint32
	FileCount          uint32
	FolderCount        uint32
	BlockSize          uint32
	TotalBlocks        uint32
	FreeBlocks         uint32
	NextAllocation     uint32
	RsrcClumpSize      uint32
	DataClumpSize      uint32
	NextCatalogID      HFSCatalogNodeID
	WriteCount         uint32
	EncodingsBitmap    uint64
	FinderInfo         [8]uint32
	AllocationFile     HFSPlusForkData
	ExtentsFile        HFSPlusForkData
	CatalogFile        HFSPlusForkData
	AttributesFile     HFSPlusForkData
	StartupFile        HFSPlusForkData
}

// B-tree

// The node descriptor contains basic information about the node as well as
// forward and backward links to other nodes. The BTNodeDescriptor data type
// describes this structure.
type BTNodeDescriptor struct {
	FLink      uint32
	BLink      uint32
	Kind       int8
	Height     uint8
	NumRecords uint16
	Reserved   uint16
}

// The B-tree header record contains general information about the B-tree such
// as its size, maximum key length, and the location of the first and last leaf
// nodes. The data type BTHeaderRec describes the structure of a header record.
type BTHeaderRec struct {
	TreeDepth      uint16
	RootNode       uint32
	LeafRecords    uint32
	FirstLeafNode  uint32
	LastLeafNode   uint32
	NodeSize       uint16
	MaxKeyLength   uint16
	TotalNodes     uint32
	FreeNodes      uint32
	Reserved1      uint16
	ClumpSize      uint32 // misaligned
	BtreeType      uint8
	KeyCompareType uint8
	Attributes     uint32 // long aligned again
	Reserved3      [16]uint32
}

// Catalog File

// A HFSCatalogNodeID (CNID) is assigned to each file or folder in the catalog
// file.
type HFSCatalogNodeID uint32

// For a given file, folder, or thread record, the catalog file key consists
// of the parent folder's CNID and the name of the file or folder. This
// structure is described using the HFSPlusCatalogKey type.
type HFSPlusCatalogKey struct {
	KeyLength uint16
	ParentID  HFSCatalogNodeID
	NodeName  HFSUniStr255
}

// The catalog folder record is used in the catalog B-tree file to hold information about a particular folder on the volume. The data of the record is described by the HFSPlusCatalogFolder type.
type HFSPlusCatalogFolder struct {
	RecordType       int16
	Flags            uint16
	Valence          uint32
	FolderID         HFSCatalogNodeID
	CreateDate       uint32
	ContentModDate   uint32
	AttributeModDate uint32
	AccessDate       uint32
	BackupDate       uint32
	Permissions      HFSPlusBSDInfo
	UserInfo         FolderInfo
	FinderInfo       ExtendedFolderInfo
	TextEncoding     uint32
	Reserved         uint32
}

// The catalog file record is used in the catalog B-tree file to hold information about a particular file on the volume. The data of the record is described by the HFSPlusCatalogFile type.
type HFSPlusCatalogFile struct {
	RecordType       int16
	Flags            uint16
	Reserved1        uint32
	FileID           HFSCatalogNodeID
	CreateDate       uint32
	ContentModDate   uint32
	AttributeModDate uint32
	AccessDate       uint32
	BackupDate       uint32
	Permissions      HFSPlusBSDInfo
	UserInfo         FileInfo
	FinderInfo       ExtendedFileInfo
	TextEncoding     uint32
	Reserved2        uint32
	DataFork         HFSPlusForkData
	ResourceFork     HFSPlusForkData
}

// The catalog thread record is used in the catalog B-tree file to link a CNID to the file or folder record using that CNID. The data of the record is described by the HFSPlusCatalogThread type.
type HFSPlusCatalogThread struct {
	RecordType int16
	Reserved   int16
	ParentID   HFSCatalogNodeID
	NodeName   HFSUniStr255
}

// Point defines a location of a file or folder in a window.
type Point struct {
	V int16
	H int16
}

// Rect defines a folder's window rectangle
type Rect struct {
	Top    int16
	Left   int16
	Bottom int16
	Right  int16
}

// FourCharCode is four 1-byte character packed together.
type FourCharCode uint32

// OSType is a 32-bit value made by packing four 1-byte characters together.
type OSType FourCharCode

// FileInfo describes attributes for files used for Finder.
type FileInfo struct {
	FileType      OSType /* The type of the file * /
	FileCreator   OSType /* The file's creator * /
	FinderFlags   uint16
	Location      Point /* File's location in the folder. * /
	ReservedField uint16
}

// ExtendedFileInfo describes additional attributes for files used for Finder.
type ExtendedFileInfo struct {
	Reserved1           [4]int16
	ExtendedFinderFlags uint16
	Reserved2           int16
	PutAwayFolderID     int32
}

// FolderInfo describes attributes for folders used for Finder.
type FolderInfo struct {
	WindowBounds Rect /* The position and dimension of the * /
	/* folder's window * /
	FinderFlags uint16
	Location    Point /* Folder's location in the parent * /
	/* folder. If set to {0, 0}, the Finder * /
	/* will place the item automatically * /
	ReservedField uint16
}

// ExtendedFolderInfo describes additional attributes for folders used for Finder.
type ExtendedFolderInfo struct {
	ScrollPosition      Point /* Scroll position (for icon views) * /
	Reserved1           int32
	ExtendedFinderFlags uint16
	Reserved2           int16
	PutAwayFolderID     int32
}

// Extents Overflow File

// The HFSPlusExtentKey describes the structure of the key for the extents overflow file.
type HFSPlusExtentKey struct {
	KeyLength  uint16
	ForkType   uint8
	Pad        uint8
	FileID     HFSCatalogNodeID
	StartBlock uint32
}

// Attributes File

// HFSPlusAttrForkData defines a fork data attribute.
type HFSPlusAttrForkData struct {
	RecordType uint32
	Reserved   uint32
	TheFork    HFSPlusForkData
}

// HFSPlusAttrExtents defines an extension attribute.
type HFSPlusAttrExtents struct {
	RecordType uint32
	Reserved   uint32
	Extents    HFSPlusExtentRecord
}

// Journal

// The journal info block describes where the journal header and journal buffer are stored. The journal info block is stored at the beginning of the allocation block whose number is stored in the journalInfoBlock field of the volume header. The journal info block is described by the data type JournalInfoBlock.
type JournalInfoBlock struct {
	Flags           uint32
	DeviceSignature [8]uint32
	Offset          uint64
	Size            uint64
	Reserved        [32]uint32
}

// The journal begins with a journal header, whose main purpose is to describe the location of transactions in the journal buffer. The journal header is stored using the JournalHeader data type.
type JournalHeader struct {
	Magic     uint32
	Endian    uint32
	Start     uint64
	End       uint64
	Size      uint64
	BlhdrSize uint32
	Checksum  uint32
	JhdrSize  uint32
}

// The block list header describes a list of blocks included in a transaction. A transaction may include several block lists if it modifies more blocks than can be represented in a single block list. The block list header is stored in a structure of type blockListHeader.
type BlockListHeader struct {
	MaxBlocks uint16
	NumBlocks uint16
	BytesUsed uint32
	Checksum  uint32
	Pad       uint32
	Binfo     [1]BlockInfo
}

// The first element of the binfo array is used to indicate whether the transaction contains additional block lists. Each of the other elements of the binfo array represent a single block of data in the journal buffer which must be copied to its correct location on disk. The fields have the following meaning:
type BlockInfo struct {
	Bnum  uint64
	Bsize uint32
	Next  uint32
}

// Hot Files

// The B-tree's user data record contains information about hot file recording. The format of the user data is described by the HotFilesInfo structure:
type HotFilesInfo struct {
	Magic       uint32
	Version     uint32
	Duration    uint32 /* duration of sample period * /
	Timebase    uint32 /* recording period start time * /
	Timeleft    uint32 /* recording period stop time * /
	Threshold   uint32
	Maxfileblks uint32
	Maxfilecnt  uint32
	Tag         [32]uint8
}

// HotFileKey is a key in the hot file B-tree.
type HotFileKey struct {
	KeyLength   uint16
	ForkType    uint8
	Pad         uint8
	Temperature uint32
	FileID      uint32
}
*/
