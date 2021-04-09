package persist

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// Message struct, stores message content and the filename where if's going to be saved
type Message struct {
	Fname, Body string
}

// A temporary buffer to store list of Message
type WriteBuffer struct {
	dir    string
	ticker *time.Ticker //snapshotting period
	queue  *Queue       //thread safe queue
	quit   chan struct{}
}

func NewWriteBuffer(directory string) (*WriteBuffer, error) {
	dirAbs, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}
	return &WriteBuffer{
		dir:    dirAbs,
		ticker: time.NewTicker(time.Second),
		queue:  newQueue(),
		quit:   make(chan struct{}),
	}, nil
}

// Init Initialize incremental snapshotting
func (w *WriteBuffer) Init() {
	go func() {
		fmt.Println("Snapshotting Initiated")
		for {
			select {
			case <-w.ticker.C:
				err := w.periodicDumpContent()
				if err != nil {
					fmt.Fprintf(os.Stderr, "error while writing commits: %s", err)
				}
			case <-w.quit:
				fmt.Println("Backup successful")
				return
			}
		}
	}()
}

// Add Message to the buffer
func (w *WriteBuffer) AppendContent(message, filename string) {
	w.queue.Push(Message{
		Fname: filename,
		Body:  message,
	})
}

// Flush the queue content by writing messages to the specific directory
func (w *WriteBuffer) periodicDumpContent() error {
	for {
		out := w.queue.Pop()
		if out == nil {
			return nil
		}
		mesg := out.(Message)
		path := filepath.Join(w.dir, mesg.Fname)

		err := ioutil.WriteFile(path, []byte(mesg.Body), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// Quit function to shut down Snapshotting
func (w *WriteBuffer) Quit() {
	w.ticker.Stop()
	close(w.quit)
	//dump the rest overs
	err := w.periodicDumpContent()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while writing commits: %s", err)
	}
}
