package fsio

import (
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		wantErr bool
	}{
		{"ErrorReader", &ErrorReader{}, true},
		{"ErrorReadSeeker", &ErrorReadSeeker{}, true},
		{"ErrorReadSeekerAt", &ErrorReadSeekerAt{}, true},
		{"ErrorReader", &ErrorReader{Skip: 1}, false},
		{"ErrorReadSeeker", &ErrorReadSeeker{Skip: 1}, false},
		{"ErrorReadSeekerAt", &ErrorReadSeekerAt{Skip: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.Read(nil)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Read() error =  %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReaderAt(t *testing.T) {
	tests := []struct {
		name    string
		r       io.ReaderAt
		wantErr bool
	}{
		{"ErrorReaderAt", &ErrorReaderAt{}, true},
		{"ErrorReadSeekerAt", &ErrorReadSeekerAt{}, true},
		{"ErrorReaderAt", &ErrorReaderAt{Skip: 1}, false},
		{"ErrorReadSeekerAt", &ErrorReadSeekerAt{Skip: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.ReadAt(nil, 0)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadAt() error =  %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSeek(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Seeker
		wantErr bool
	}{
		{"ErrorSeeker", &ErrorSeeker{}, true},
		{"ErrorReadSeeker", &ErrorReadSeeker{}, true},
		{"ErrorReadSeekerAt", &ErrorReadSeekerAt{}, true},
		{"ErrorSeeker", &ErrorSeeker{Skip: 1}, false},
		{"ErrorReadSeeker", &ErrorReadSeeker{Skip: 1}, false},
		{"ErrorReadSeekerAt", &ErrorReadSeekerAt{Skip: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.Seek(0, 0)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Seek() error =  %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Writer
		wantErr bool
	}{
		{"ErrorWriter", &ErrorWriter{}, true},
		{"ErrorWriter", &ErrorWriter{Skip: 1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.r.Write(nil)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Write() error =  %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
