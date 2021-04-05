package utils

import "testing"

func TestBuffer(t *testing.T) {
	w := NewWriteBuffer()
	w.AppendContent("hi there", "lol")
	w.AppendContent("nice", "lol2")
	err := w.DumpContent("./temp")

	if err != nil {
		t.Fatal(err)
	}
}
