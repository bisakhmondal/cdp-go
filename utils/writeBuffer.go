package utils

import (
	"io/ioutil"
	"path/filepath"
)

type Message struct {
	Fname, Body string
}

type WriteBuffer struct {
	list []Message
}

func NewWriteBuffer() *WriteBuffer {
	return &WriteBuffer{list: []Message{}}
}

func (w *WriteBuffer) AppendContent(message, filename string) {
	w.list = append(w.list, Message{
		Fname: filename,
		Body:  message,
	})
}

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
