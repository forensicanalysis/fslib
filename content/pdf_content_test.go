package content

/*
type read struct{}

func (b *read) Read([]byte) (n int, err error) { return 0, nil }

type readAt struct{}

func (b *readAt) ReadAt([]byte, int64) (n int, err error) { return 0, nil }

type brokenReadAt struct{}

func (b *brokenReadAt) ReadAt([]byte, int64) (n int, err error) {
	return 0, errors.New("broken readerAt")
}

type brokenSeek struct{}

func (b *brokenSeek) Seek(int64, int) (int64, error) {
	return 0, errors.New("broken seeker")
}

type brokenSeeker struct {
	brokenSeek
	readAt
	read
}

func plainTextError(io.ReaderAt, int64) (reader *pdf.Reader, err error) {
	return nil, errors.New("broken seeker")
}

func brokenReader(io.ReaderAt, int64) (reader *pdf.Reader, err error) {
	return pdf.NewReader(&brokenReadAt{}, 0)
}

func TestPDFContent(t *testing.T) {
	type args struct {
		r               fsio.ReadSeekerAt
		readerGenerator func(f io.ReaderAt, size int64) (reader *pdf.Reader, err error)
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"size error", args{&brokenSeeker{}, pdf.NewReader}, "", true},
		{"no pdf", args{bytes.NewReader([]byte{}), pdf.NewReader}, "", true},
		{"plaintext error", args{bytes.NewReader([]byte{}), plainTextError}, "", false},
		{"broken reader", args{bytes.NewReader([]byte{}), brokenReader}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readerGenerator = tt.args.readerGenerator
			got, err := PDFContent(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("PDFContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			if got != tt.want {
				t.Errorf("PDFContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
