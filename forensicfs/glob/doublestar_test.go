// Copyright (c) 2014-2019 Bob Matcuk
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
// Author(s): Bob Matcuk, Jonas Plum
//
// This code was adapted from
// https://github.com/bmatcuk/doublestar
// for use with forensic filesystems.

package glob

import (
	"path"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/forensicanalysis/fslib"
	"github.com/forensicanalysis/fslib/filesystem/testfs"
)

type MatchTest struct {
	pattern, testPath string // a pattern and path to test the pattern on
	shouldMatch       bool   // true if the pattern should match the path
	expectedErr       error  // an expected error
	testOnDisk        bool   // true: test pattern against files in "test" directory
}

var matchTests = []MatchTest{
	{"/*", "/", true, nil, false},
	{"/*", "//", false, nil, false},
	{"//*", "//", true, nil, false},
	{"//*", "//debug/", false, nil, false},
	{"//*", "///", false, nil, false},
	{"/abc", "/abc", true, nil, true},
	{"/*", "/abc", true, nil, true},
	{"/*c", "/abc", true, nil, true},
	{"/a*", "/a", true, nil, true},
	{"/a*", "/abc", true, nil, true},
	{"/a*", "/ab/c", false, nil, true},
	{"/a*/b", "/abc/b", true, nil, true},
	{"/a*/b", "/a/c/b", false, nil, true},
	{"/a*b*c*d*e*/f", "/axbxcxdxe/f", true, nil, true},
	{"/a*b*c*d*e*/f", "/axbxcxdxexxx/f", true, nil, true},
	{"/a*b*c*d*e*/f", "/axbxcxdxe/xxx/f", false, nil, true},
	{"/a*b*c*d*e*/f", "/axbxcxdxexxx/fff", false, nil, true},
	{"/a*b?c*x", "/abxbbxdbxebxczzx", true, nil, true},
	{"/a*b?c*x", "/abxbbxdbxebxczzy", false, nil, true},
	{"/ab[c]", "/abc", true, nil, true},
	{"/ab[b-d]", "/abc", true, nil, true},
	{"/ab[e-g]", "/abc", false, nil, true},
	{"/ab[^c]", "/abc", false, nil, true},
	{"/ab[^b-d]", "/abc", false, nil, true},
	{"/ab[^e-g]", "/abc", true, nil, true},
	{"/a\\*b", "/ab", false, nil, true},
	{"/a?b", "/a☺b", true, nil, true},
	{"/a[^a]b", "/a☺b", true, nil, true},
	{"/a???b", "/a☺b", false, nil, true},
	{"/a[^a][^a][^a]b", "/a☺b", false, nil, true},
	{"/[a-ζ]*", "/α", true, nil, true},
	{"/*[a-ζ]", "/A", false, nil, true},
	{"/a?b", "/a/b", false, nil, true},
	{"/a*b", "/a/b", false, nil, true},
	{"/[]a]", "/]", false, ErrBadPattern, true},
	{"/[-]", "/-", false, ErrBadPattern, true},
	{"/[x-]", "/x", false, ErrBadPattern, true},
	{"/[x-]", "/-", false, ErrBadPattern, true},
	{"/[x-]", "/z", false, ErrBadPattern, true},
	{"/[-x]", "/x", false, ErrBadPattern, true},
	{"/[-x]", "/-", false, ErrBadPattern, true},
	{"/[-x]", "/a", false, ErrBadPattern, true},
	{"/[a-b-c]", "/a", false, ErrBadPattern, true},
	{"/[", "/a", false, ErrBadPattern, true},
	{"/[^", "/a", false, ErrBadPattern, true},
	{"/[^bc", "/a", false, ErrBadPattern, true},
	{"/a[", "/a", false, nil, false},
	{"/a[", "/ab", false, ErrBadPattern, true},
	{"/*x", "/xxx", true, nil, true},
	{"/[abc]", "/b", true, nil, true},
	{"/a/**", "/a", false, nil, true},
	{"/a/**", "/a/b", true, nil, true},
	{"/a/**", "/a/b/c", true, nil, true},
	{"/**/c", "/c", true, nil, true},
	{"/**/c", "/b/c", true, nil, true},
	{"/**/c", "/a/b/c", true, nil, true},
	{"/**/c", "/a/b", false, nil, true},
	{"/**/c", "/abcd", false, nil, true},
	{"/**/c", "/a/abc", false, nil, true},
	{"/a/**/b", "/a/b", true, nil, true},
	{"/a/**/c", "/a/b/c", true, nil, true},
	{"/a/**/d", "/a/b/c/d", true, nil, true},
	// {"/a//b/c", "/a/b/c", true, nil, true},
	// {"/a/b/c", "/a/b//c", true, nil, true},
	{"/ab{c,d}", "/abc", true, nil, true},
	{"/ab{c,d,*}", "/abcde", true, nil, true},
	{"/ab{c,d}[", "/abcd", false, ErrBadPattern, true},
	{"/abc/**", "/abc/b", true, nil, true},
	{"/**/abc", "/abc", true, nil, true},
	{"/abc**", "/abc/b", false, nil, true},
	{"/abc**", "/abc/b", false, nil, true},
	{"/**2/d", "/a/b/c/d", false, nil, true},
	{"/a/**2/d", "/a/b/c/d", true, nil, true},
	{"/**3/d", "/a/b/c/d", true, nil, true},
	{"/**5/d", "/a/b/c/d", true, nil, true},
	{"/**/d", "/f/g/h/i/j/k/d", false, nil, true},
	{"/**5/d", "/f/g/h/i/j/k/d", false, nil, true},
	{"/**6/d", "/f/g/h/i/j/k/d", true, nil, true},
	{"/**7/d", "/f/g/h/i/j/k/d", true, nil, true},
}

func TestMatch(t *testing.T) {
	for idx, tt := range matchTests {
		// Since Match() always uses "/" as the separator, we
		// don't need to worry about the tt.testOnDisk flag
		testMatchWith(t, idx, tt)
	}
}

func testMatchWith(t *testing.T, idx int, tt MatchTest) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("#%v. Match(%#q, %#q) panicked: %#v", idx, tt.pattern, tt.testPath, r)
		}
	}()

	// Match() always uses "/" as the separator
	ok, err := Match(tt.pattern, tt.testPath)
	if ok != tt.shouldMatch || err != tt.expectedErr {
		t.Errorf("#%v. Match(%#q, %#q) = %v, %v want %v, %v", idx, tt.pattern, tt.testPath, ok, err, tt.shouldMatch, tt.expectedErr)
	}

	if isStandardPattern(tt.pattern) {
		stdOk, stdErr := path.Match(tt.pattern, tt.testPath)
		if ok != stdOk || !compareErrors(err, stdErr) {
			t.Errorf("#%v. Match(%#q, %#q) != path.Match(...). Got %v, %v want %v, %v", idx, tt.pattern, tt.testPath, ok, err, stdOk, stdErr)
		}
	}
}

func TestPathMatch(t *testing.T) {
	for idx, tt := range matchTests {
		// Even though we aren't actually matching paths on disk, we are using
		// PathMatch() which will use the system's separator. As a result, any
		// patterns that might cause problems on-disk need to also be avoided
		// here in this test.
		if tt.testOnDisk {
			testPathMatchWith(t, idx, tt)
		}
	}
}

func testPathMatchWith(t *testing.T, idx int, tt MatchTest) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("#%v. Match(%#q, %#q) panicked: %#v", idx, tt.pattern, tt.testPath, r)
		}
	}()

	ok, err := PathMatch(tt.pattern, tt.testPath)
	if ok != tt.shouldMatch || err != tt.expectedErr {
		t.Errorf("#%v. Match(%#q, %#q) = %v, %v want %v, %v", idx, tt.pattern, tt.testPath, ok, err, tt.shouldMatch, tt.expectedErr)
	}

	if isStandardPattern(tt.pattern) {
		stdOk, stdErr := path.Match(tt.pattern, tt.testPath)
		if ok != stdOk || !compareErrors(err, stdErr) {
			t.Errorf("#%v. PathMatch(%#q, %#q) != path.Match(...). Got %v, %v want %v, %v", idx, tt.pattern, tt.testPath, ok, err, stdOk, stdErr)
		}
	}
}

func TestGlob(t *testing.T) {
	for idx, tt := range matchTests {
		if tt.testOnDisk {
			// test both relative paths and absolute paths
			testGlobWith(t, idx, tt)
		}
	}
}

func testGlobWith(t *testing.T, idx int, tt MatchTest) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("#%v. Glob(%#q) panicked: %#v", idx, tt.pattern, r)
		}
	}()

	tfs := getTestFS()
	pattern := path.Join("/test", tt.pattern)
	testPath := path.Join("/test", tt.testPath)
	matches, err := Glob(tfs, pattern)
	if inSlice(testPath, matches) != tt.shouldMatch {
		if tt.shouldMatch {
			t.Errorf("#%v. Glob(%#q) = %#v - doesn't contain %v, but should", idx, pattern, matches, tt.testPath)
		} else {
			t.Errorf("#%v. Glob(%#q) = %#v - contains %v, but shouldn't", idx, pattern, matches, tt.testPath)
		}
	}
	if err != tt.expectedErr {
		t.Errorf("#%v. Glob(%#q) has error %v, but should be %v", idx, pattern, err, tt.expectedErr)
	}
}

func isStandardPattern(pattern string) bool {
	return !strings.Contains(pattern, "**") && indexRuneWithEscaping(pattern, '{') == -1
}

func compareErrors(a, b error) bool {
	if a == nil {
		return b == nil
	}
	return b != nil
}

func inSlice(s string, a []string) bool {
	for _, i := range a {
		if i == s {
			return true
		}
	}
	return false
}

func getTestFS() *testfs.FS {
	infs := &testfs.FS{}
	// create the test directory
	dirs := []string{"/test/a/b/c", "/test/a/c", "/test/abc", "/test/axbxcxdxe/xxx", "/test/axbxcxdxexxx", "/test/b", "/test/f/g/h/i/j/k"}
	for _, dir := range dirs {
		infs.CreateDir(dir)
	}

	// create test files
	files := []string{
		"/test/a/abc", "/test/a/b/c/d", "/test/a/c/b", "/test/abc/b", "/test/abcd", "/test/abcde", "/test/abxbbxdbxebxczzx",
		"/test/abxbbxdbxebxczzy", "/test/axbxcxdxe/f", "/test/axbxcxdxe/xxx/f", "/test/axbxcxdxexxx/f",
		"/test/axbxcxdxexxx/fff", "/test/a☺b", "/test/b/c", "/test/c", "/test/x", "/test/xxx", "/test/z",
		"/test/α", "/test/f/g/h/i/j/k/l", "/test/f/g/h/i/j/k/d", "/test/f/g/h/i/j/k/u.bin", "/test/f/g/h/i/j/k/v.bin",
	}

	for _, file := range files {
		infs.CreateFile(file, []byte("test"))
	}

	return infs
}

func getInFS() fslib.FS {
	infs := &testfs.FS{}
	content := []byte("test")
	dirs := []string{"/dir/", "/dir/a/", "/dir/b/", "/dir/a/a/", "/dir/a/b/", "/dir/b/a/", "/dir/b/b/"}
	for _, dir := range dirs {
		infs.CreateDir(dir)
	}
	files := []string{"/foo.bin", "/dir/bar.bin", "/dir/baz.bin", "/dir/a/a/foo.bin", "/dir/a/b/foo.bin", "/dir/b/a/foo.bin", "/dir/b/b/foo.bin"}
	for _, file := range files {
		infs.CreateFile(file, content)
	}
	return infs
}

func Test_expandPath(t *testing.T) {
	type args struct {
		fs fslib.FS
		in string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"Expand path 1", args{getInFS(), "/*/bar.bin"}, []string{"/dir/bar.bin"}},
		{"Expand path 2", args{getInFS(), "/dir/*.bin"}, []string{"/dir/bar.bin", "/dir/baz.bin"}},
		{"Expand path 3", args{getInFS(), "/dir/*/*/foo.bin"}, []string{"/dir/a/a/foo.bin", "/dir/a/b/foo.bin", "/dir/b/a/foo.bin", "/dir/b/b/foo.bin"}},
		{"Expand path 4", args{getInFS(), "/**"}, []string{"/dir", "/dir/a", "/dir/a/a", "/dir/a/b", "/dir/b", "/dir/b/a", "/dir/b/b", "/dir/bar.bin", "/dir/baz.bin", "/foo.bin"}},
		{"Expand path 5", args{getInFS(), "/dir/**2/foo.bin"}, []string{"/dir/a/a/foo.bin", "/dir/a/b/foo.bin", "/dir/b/a/foo.bin", "/dir/b/b/foo.bin"}},
		{"Expand path 6", args{getInFS(), "/dir/**1"}, []string{"/dir/a", "/dir/b", "/dir/bar.bin", "/dir/baz.bin"}},
		{"Expand path 7", args{getInFS(), "/dir/**10"}, []string{"/dir/a", "/dir/a/a", "/dir/a/a/foo.bin", "/dir/a/b", "/dir/a/b/foo.bin", "/dir/b", "/dir/b/a", "/dir/b/a/foo.bin", "/dir/b/b", "/dir/b/b/foo.bin", "/dir/bar.bin", "/dir/baz.bin"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Glob(tt.args.fs, tt.args.in)
			if err != nil {
				t.Fatal(err)
			}
			sort.Strings(tt.want)
			sort.Strings(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expandPath(%s) = %v, want %v", tt.args.in, got, tt.want)
			}
		})
	}
}
