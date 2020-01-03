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
	"bytes"
	"testing"
)

func Test_xmlContent(t *testing.T) {

	data := []byte(`
<content>
    <p>this is content area</p>
    <animal>
        <p>This id dog</p>
        <dog>
           <p>tommy</p>
        </dog>
    </animal>
    <birds>
        <p>this is birds</p>
        <p>this is birds</p>
    </birds>
    <animal>
        <p>this is animals</p>
    </animal>
</content>
`)

	type args struct {
		r *bytes.Reader
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"xml Content", args{bytes.NewReader(data)}, "this is content area This id dog tommy this is birds this is birds this is animals"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := xmlContent(tt.args.r); got != tt.want {
				t.Errorf("xmlContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
