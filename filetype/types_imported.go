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

import (
	"github.com/h2non/filetype/matchers"
)

var importedTypes = []*Filetype{Wasm, Epub, Zip, Tar, Rar, Gz, Bz2, Sevenz, Xz,
	Pdf, Exe, Swf, Rtf, Eot, Ps, Sqlite, Nes, Crx, Cab, Deb, Ar, Z, Lz, Rpm, Elf,
	Dcm, Iso, Midi, Mp3, M4a, Ogg, Flac, Wav, Amr, Aac, Doc, Docx, Xls, Xlsx, Ppt,
	Pptx, Woff, Woff2, Ttf, Otf, Jpeg, Jpeg2000, Png, Gif, Webp, CR2, Tiff, Bmp,
	Jxr, Psd, Ico, Heif, Dwg, Mp4, M4v, Mkv, Webm, Mov, Avi, Wmv, Mpeg, Flv, Threegp,
}

var (
	// Application

	// Wasm is the file type for Wasm files.
	Wasm = newFiletype("wasm", matchers.TypeWasm, matchers.Wasm, 0)

	// Archive

	// Epub is the file type for Epub files.
	Epub = newFiletype("epub", matchers.TypeEpub, matchers.Epub, 0)
	// Zip is the file type for Zip files.
	Zip = newFiletype("zip", matchers.TypeZip, matchers.Zip, 10)
	// Tar is the file type for Tar files.
	Tar = newFiletype("tar", matchers.TypeTar, matchers.Tar, 0)
	// Rar is the file type for Rar files.
	Rar = newFiletype("rar", matchers.TypeRar, matchers.Rar, 0)
	// Gz is the file type for Gz files.
	Gz = newFiletype("gz", matchers.TypeGz, matchers.Gz, 0)
	// Bz2 is the file type for Bz2 files.
	Bz2 = newFiletype("bz2", matchers.TypeBz2, matchers.Bz2, 0)
	// Sevenz is the file type for Sevenz files.
	Sevenz = newFiletype("7z", matchers.Type7z, matchers.SevenZ, 0)
	// Xz is the file type for Xz files.
	Xz = newFiletype("xz", matchers.TypeXz, matchers.Xz, 0)
	// Pdf is the file type for Pdf files.
	Pdf = newFiletype("pdf", matchers.TypePdf, matchers.Pdf, 0)
	// Exe is the file type for Exe files.
	Exe = newFiletype("exe", matchers.TypeExe, matchers.Exe, 0)
	// Swf is the file type for Swf files.
	Swf = newFiletype("swf", matchers.TypeSwf, matchers.Swf, 0)
	// Rtf is the file type for Rtf files.
	Rtf = newFiletype("rtf", matchers.TypeRtf, matchers.Rtf, 0)
	// Eot is the file type for Eot files.
	Eot = newFiletype("eot", matchers.TypeEot, matchers.Eot, 0)
	// Ps is the file type for Ps files.
	Ps = newFiletype("ps", matchers.TypePs, matchers.Ps, 0)
	// Sqlite is the file type for Sqlite files.
	Sqlite = newFiletype("sqlite", matchers.TypeSqlite, matchers.Sqlite, 0)
	// Nes is the file type for Nes files.
	Nes = newFiletype("nes", matchers.TypeNes, matchers.Nes, 0)
	// Crx is the file type for Crx files.
	Crx = newFiletype("crx", matchers.TypeCrx, matchers.Crx, 0)
	// Cab is the file type for Cab files.
	Cab = newFiletype("cab", matchers.TypeCab, matchers.Cab, 0)
	// Deb is the file type for Deb files.
	Deb = newFiletype("deb", matchers.TypeDeb, matchers.Deb, 0)
	// Ar is the file type for Ar files.
	Ar = newFiletype("ar", matchers.TypeAr, matchers.Ar, 0)
	// Z is the file type for Z files.
	Z = newFiletype("Z      =", matchers.TypeZ, matchers.Z, 0)
	// Lz is the file type for Lz files.
	Lz = newFiletype("lz", matchers.TypeLz, matchers.Lz, 0)
	// Rpm is the file type for Rpm files.
	Rpm = newFiletype("rpm", matchers.TypeRpm, matchers.Rpm, 0)
	// Elf is the file type for Elf files.
	Elf = newFiletype("elf", matchers.TypeElf, matchers.Elf, 0)
	// Dcm is the file type for Dcm files.
	Dcm = newFiletype("dcm", matchers.TypeDcm, matchers.Dcm, 0)
	// Iso is the file type for Iso files.
	Iso = newFiletype("iso", matchers.TypeIso, matchers.Iso, 0)

	// Audio

	// Midi is the file type for Midi files.
	Midi = newFiletype("midi", matchers.TypeMidi, matchers.Midi, 0)
	// Mp3 is the file type for Mp3 files.
	Mp3 = newFiletype("mp3", matchers.TypeMp3, matchers.Mp3, 0)
	// M4a is the file type for M4a files.
	M4a = newFiletype("m4a", matchers.TypeM4a, matchers.M4a, 0)
	// Ogg is the file type for Ogg files.
	Ogg = newFiletype("ogg", matchers.TypeOgg, matchers.Ogg, 0)
	// Flac is the file type for Flac files.
	Flac = newFiletype("flac", matchers.TypeFlac, matchers.Flac, 0)
	// Wav is the file type for Wav files.
	Wav = newFiletype("wav", matchers.TypeWav, matchers.Wav, 0)
	// Amr is the file type for Amr files.
	Amr = newFiletype("amr", matchers.TypeAmr, matchers.Amr, 0)
	// Aac is the file type for Aac files.
	Aac = newFiletype("aac", matchers.TypeAac, matchers.Aac, 0)

	// Document

	// Doc is the file type for Doc files.
	Doc = newFiletype("doc", matchers.TypeDoc, matchers.Doc, 0)
	// Docx is the file type for Docx files.
	Docx = newFiletype("docx", matchers.TypeDocx, matchers.Docx, 0)
	// Xls is the file type for Xls files.
	Xls = newFiletype("xls", matchers.TypeXls, matchers.Xls, 0)
	// Xlsx is the file type for Xlsx files.
	Xlsx = newFiletype("xlsx", matchers.TypeXlsx, matchers.Xlsx, 0)
	// Ppt is the file type for Ppt files.
	Ppt = newFiletype("ppt", matchers.TypePpt, matchers.Ppt, 0)
	// Pptx is the file type for Pptx files.
	Pptx = newFiletype("pptx", matchers.TypePptx, matchers.Pptx, 0)

	// Font

	// Woff is the file type for Woff files.
	Woff = newFiletype("woff", matchers.TypeWoff, matchers.Woff, 0)
	// Woff2 is the file type for Woff2 files.
	Woff2 = newFiletype("woff2", matchers.TypeWoff2, matchers.Woff2, 0)
	// Ttf is the file type for Ttf files.
	Ttf = newFiletype("ttf", matchers.TypeTtf, matchers.Ttf, 0)
	// Otf is the file type for Otf files.
	Otf = newFiletype("otf", matchers.TypeOtf, matchers.Otf, 0)

	// Image

	// Jpeg is the file type for Jpeg files.
	Jpeg = newFiletype("jpeg", matchers.TypeJpeg, matchers.Jpeg, 0)
	// Jpeg2000 is the file type for Jpeg2000 files.
	Jpeg2000 = newFiletype("jpeg2000", matchers.TypeJpeg2000, matchers.Jpeg2000, 0)
	// Png is the file type for Png files.
	Png = newFiletype("png", matchers.TypePng, matchers.Png, 0)
	// Gif is the file type for Gif files.
	Gif = newFiletype("gif", matchers.TypeGif, matchers.Gif, 0)
	// Webp is the file type for Webp files.
	Webp = newFiletype("webp", matchers.TypeWebp, matchers.Webp, 0)
	// CR2 is the file type for CR2 files.
	CR2 = newFiletype("cr2", matchers.TypeCR2, matchers.CR2, 0)
	// Tiff is the file type for Tiff files.
	Tiff = newFiletype("tiff", matchers.TypeTiff, matchers.Tiff, 0)
	// Bmp is the file type for Bmp files.
	Bmp = newFiletype("bmp", matchers.TypeBmp, matchers.Bmp, 0)
	// Jxr is the file type for Jxr files.
	Jxr = newFiletype("jxr", matchers.TypeJxr, matchers.Jxr, 0)
	// Psd is the file type for Psd files.
	Psd = newFiletype("psd", matchers.TypePsd, matchers.Psd, 0)
	// Ico is the file type for Ico files.
	Ico = newFiletype("ico", matchers.TypeIco, matchers.Ico, 0)
	// Heif is the file type for Heif files.
	Heif = newFiletype("heif", matchers.TypeHeif, matchers.Heif, 0)
	// Dwg is the file type for Dwg files.
	Dwg = newFiletype("dwg", matchers.TypeDwg, matchers.Dwg, 0)

	// Video

	// Mp4 is the file type for Mp4 files.
	Mp4 = newFiletype("mp4", matchers.TypeMp4, matchers.Mp4, 0)
	// M4v is the file type for M4v files.
	M4v = newFiletype("m4v", matchers.TypeM4v, matchers.M4v, 0)
	// Mkv is the file type for Mkv files.
	Mkv = newFiletype("mkv", matchers.TypeMkv, matchers.Mkv, 0)
	// Webm is the file type for Webm files.
	Webm = newFiletype("webm", matchers.TypeWebm, matchers.Webm, 0)
	// Mov is the file type for Mov files.
	Mov = newFiletype("mov", matchers.TypeMov, matchers.Mov, 0)
	// Avi is the file type for Avi files.
	Avi = newFiletype("avi", matchers.TypeAvi, matchers.Avi, 0)
	// Wmv is the file type for Wmv files.
	Wmv = newFiletype("wmv", matchers.TypeWmv, matchers.Wmv, 0)
	// Mpeg is the file type for Mpeg files.
	Mpeg = newFiletype("mpeg", matchers.TypeMpeg, matchers.Mpeg, 0)
	// Flv is the file type for Flv files.
	Flv = newFiletype("flv", matchers.TypeFlv, matchers.Flv, 0)
	// Threegp is the file type for 3gp files.
	Threegp = newFiletype("3gp", matchers.Type3gp, matchers.Match3gp, 0)
)
