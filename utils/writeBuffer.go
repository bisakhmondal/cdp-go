package utils

import (
	"io/ioutil"
	"path/filepath"
)

// Message struct, stores message content and the filename where if's going to be saved
type Message struct {
	Fname, Body string
}

// A temporary buffer to store list of Message
type WriteBuffer struct {
	list []Message
}

func NewWriteBuffer() *WriteBuffer {
	return &WriteBuffer{list: []Message{}}
}

// Add Message to the buffer
func (w *WriteBuffer) AppendContent(message, filename string) {
	w.list = append(w.list, Message{
		Fname: filename,
		Body:  message,
	})
}

// Flush buffer by writing messages to the specific directory
func (w *WriteBuffer) DumpContent(directory string) error {
	for _, mesg := range w.list {
		dir, err := filepath.Abs(directory)
		if err != nil {
			return err
		}
		path := filepath.Join(dir, mesg.Fname)

		err = ioutil.WriteFile(path, []byte(mesg.Body), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
