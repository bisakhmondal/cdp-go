package persist

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
)

// Message struct, stores message content and the filename where if's going to be saved
type Message struct {
	Fname, Body string
}

// A temporary buffer to store list of Message
type WriteBuffer struct {
	dir    string
	queue  chan Message //thread safe queue
	wgroup sync.WaitGroup
}

const capacity = 16

func NewWriteBuffer(directory string) (*WriteBuffer, error) {
	dirAbs, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}
	return &WriteBuffer{
		dir:   dirAbs,
		queue: make(chan Message, capacity),
	}, nil
}

// Init Initialize incremental snapshotting
func (w *WriteBuffer) Init(quit chan error) {
	w.wgroup.Add(1)

	go func() {
		defer w.wgroup.Done()
		fmt.Println("Backup Initiated")
		for mesg := range w.queue {
			path := filepath.Join(w.dir, mesg.Fname)

			err := ioutil.WriteFile(path, []byte(mesg.Body), 0644)
			if err != nil {
				quit <- fmt.Errorf("error while writing content to disk: %s\n", err)
				return
			}
		}
	}()
}

// Add Message to the buffer
func (w *WriteBuffer) AppendContent(message, filename string) {
	w.queue <- Message{
		Fname: filename,
		Body:  message,
	}
}

// Quit function to shut down Snapshotting
func (w *WriteBuffer) Quit() {
	close(w.queue)
	// wait for the buffer to get fully empty
	w.wgroup.Wait()
	fmt.Println("Backup Successful")
}
