// +build !go1.8

package content

import (
	"errors"
	"io"

	"github.com/forensicanalysis/fslib/fsio"
)

// PDFContent returns the text data from a pdf file.
func PDFContent(r fsio.ReadSeekerAt) (io.Reader, error) {
	return nil, errors.New("Old Go Version")
}
