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

package content

import (
	"bufio"
	"bytes"
	"io"
)

// StringsReader returns the text data from io.Reader.
func StringsReader(r io.Reader, w io.Writer) (err error) {
	buffer := make([]byte, 4096*4096)
	var currentString bytes.Buffer
	f := bufio.NewWriter(w)
	defer f.Flush()
	for {
		size, _ := r.Read(buffer)
		if size == 0 {
			break
		}
		if err = extractString(buffer[:size], &currentString, f); err != nil {
			return err
		}
	}

	if currentString.Len() >= 4 {
		if _, err = currentString.WriteTo(f); err != nil {
			return err
		}

		if err = f.WriteByte('\n'); err != nil {
			return err
		}
	}
	return
}

func extractString(data []byte, currentString *bytes.Buffer, f *bufio.Writer) (err error) {
	for _, c := range data {
		if (c >= ' ' && c <= '~') || c == '' {
			if err = currentString.WriteByte(c); err != nil {
				return err
			}
		} else {
			if currentString.Len() >= 4 {
				if _, err = currentString.WriteTo(f); err != nil {
					return err
				}
				if err = f.WriteByte('\n'); err != nil {
					return err
				}
			} else if currentString.Len() > 0 {
				currentString.Reset()
			}
		}
	}
	return nil
}
