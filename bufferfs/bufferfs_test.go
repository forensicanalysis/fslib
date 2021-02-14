package bufferfs

import (
	"testing"
	"testing/fstest"
)

func TestFS(t *testing.T) {
	mem1 := fstest.MapFS{}

	mem1["foo"] = &fstest.MapFile{Data: []byte("fs1")}

	fallbackFS := New(mem1)

	err := fstest.TestFS(fallbackFS, "foo")
	if err != nil {
		t.Error(err)
	}
}
