#!/usr/bin/env bats
# Copyright (c) 2019 Siemens AG
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
# the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
# FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
# IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#
# Author(s): Jonas Plum


@test "test ls zip" {
  fs ls "test/data/container/zip.zip/"
}

@test "test cat zip" {
  fs cat "test/data/container/zip.zip/document/Computer forensics - Wikipedia.pdf" | cmp "test/data/document/Computer forensics - Wikipedia.pdf"
}

@test "test ls fat16" {
  fs ls "test/data/filesystem/fat16.dd/"
}

# @test "test ls fat16 folder" {
#   run fs ls "test/data/filesystem/fat16.dd/image"
#   [ "$status" -eq 0 ]
#   [ "$output" = "alps.jpg alps.png alps.tiff" ]
# }

# @test "test cat fat16" {
#   fs cat "test/data/filesystem/fat16.dd/image/alps.jpg" | cmp "test/data/image/alps.jpg"
# }