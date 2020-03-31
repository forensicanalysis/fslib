package fsio

import (
	"bytes"
	"io"
	"testing"
)

func TestDecoderAtWrapper_ReadAt(t *testing.T) {
	type fields struct {
		ReadSeeker io.ReadSeeker
	}
	type args struct {
		p   []byte
		off int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
	}{
		{"readat", fields{bytes.NewReader([]byte{1, 2})}, args{[]byte{1}, 0}, 1, false},
		{"readat eof", fields{bytes.NewReader([]byte{})}, args{nil, 0}, 0, true},
		{"fail 1. seek", fields{ReadSeeker: &ErrorReadSeeker{ErrorSeeker{Whence: -1}, ErrorReader{}}}, args{nil, 0}, 0, true},
		{"fail 2. seek", fields{ReadSeeker: &ErrorReadSeeker{}}, args{nil, 0}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			da := &DecoderAtWrapper{
				ReadSeeker: tt.fields.ReadSeeker,
			}
			gotN, err := da.ReadAt(tt.args.p, tt.args.off)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("ReadAt() gotN = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func TestGetSize(t *testing.T) {
	type args struct {
		seeker io.Seeker
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"get zero size", args{bytes.NewReader([]byte{})}, 0, false},
		{"get size", args{bytes.NewReader([]byte{0})}, 1, false},
		{"fail 1. seek", args{&ErrorSeeker{Whence: io.SeekCurrent}}, 0, true},
		{"fail 2. seek", args{&ErrorSeeker{Whence: io.SeekEnd}}, 0, true},
		{"fail 3. seek", args{&ErrorSeeker{Whence: io.SeekStart}}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSize(tt.args.seeker)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSize() got = %v, want %v", got, tt.want)
			}
		})
	}
}
