package fsio

import "testing"

func TestErrorReaderAt_ReadAt(t *testing.T) {
	b := &ErrorReaderAt{}
	_, err := b.ReadAt(nil, 0)
	if err == nil {
		t.Fatalf("ReadAt() error = nil, wantErr true")
	}
}

func TestErrorReader_Read(t *testing.T) {
	b := &ErrorReader{}
	_, err := b.Read(nil)
	if err == nil {
		t.Fatalf("Read() error = nil, wantErr true")
	}
}

func TestErrorSeeker_Seek(t *testing.T) {
	b := &ErrorSeeker{}
	_, err := b.Seek(0, 0)
	if err == nil {
		t.Fatalf("Seek() error = nil, wantErr true")
	}
}

func TestErrorWriter_Write(t *testing.T) {
	b := &ErrorWriter{}
	_, err := b.Write([]byte{0x00})
	if err == nil {
		t.Fatalf("Write() error = nil, wantErr true")
	}
}

func TestErrorWriter_WriteOK(t *testing.T) {
	b := &ErrorWriter{Skip: 1}
	_, err := b.Write([]byte{0x00})
	if err != nil {
		t.Fatalf("Write() error != nil, wantErr false")
	}
}
