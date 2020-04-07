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
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

/*
type node struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
	Nodes   []node `xml:",any"`
}
*/

func xmlContent(r io.ReadSeeker) string {
	// dec := xml2map.NewDecoder(r)
	// m, _ := dec.Decode()
	// log.Printf("%#v", m)

	r.Seek(0, 0) // nolint: errcheck
	b, _ := ioutil.ReadAll(r)

	// re := regexp.MustCompile("<w:t[^>]*>([^>]*)</w:t>")
	// out := re.FindAllStringSubmatch(string(b), -1)
	// s := ""
	// for _, o := range out {
	// 	s += o[1]
	// }
	// fmt.Println(s)

	re := regexp.MustCompile("<[^>]*>")
	b = re.ReplaceAll(b, []byte(" "))

	s := string(b)
	re = regexp.MustCompile(`\s+`)
	s = re.ReplaceAllString(s, " ")

	// walkMap(reflect.ValueOf(m))
	return strings.TrimSpace(s)
}

/*
	dec := xml.NewDecoder(r)

	var n node
	err := dec.Decode(&n)
	if err != nil {
		panic(err)
	}

	buf := bytes.Buffer{}
	walk([]node{n}, func(n node) bool {
		if n.XMLName.Local == "p" {
			_, err := buf.Write(n.Content)
			if err != nil {
				return false
			}
			err = buf.WriteByte('\n')
			if err != nil {
				return false
			}
		}
		return true
	})
	return buf.String()
*/

// func walk(nodes []node, f func(node) bool) {
// 	for _, n := range nodes {
// 		if f(n) {
// 			walk(n.Nodes, f)
// 		}
// 	}
// }

// func walkMap(v reflect.Value) {
// 	// Indirect through pointers and interfaces
// 	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
// 		v = v.Elem()
// 	}
// 	switch v.Kind() {
// 	case reflect.Array, reflect.Slice:
// 		for i := 0; i < v.Len(); i++ {
// 			walkMap(v.Index(i))
// 		}
// 	case reflect.Map:
// 		for _, k := range v.MapKeys() {
// 			walkMap(v.MapIndex(k))
// 			if k.String() == "t" {
// 				fmt.Print(v.MapIndex(k))
// 			}
// 		}
// 	default:
// 		// handle other types
// 	}
// }
