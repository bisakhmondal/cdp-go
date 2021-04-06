package utils

import "testing"

func TestBuffer(t *testing.T) {
	w := NewWriteBuffer()
	w.AppendContent("hi there", "commit1")
	w.AppendContent("nice", "commit2")
	err := w.DumpContent("./commits")

	if err != nil {
		t.Fatal(err)
	}
}
