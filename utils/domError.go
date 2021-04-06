package utils

import "fmt"

type DOMError struct {
	Err error
}

func (e *DOMError) Error() string {
	return fmt.Sprintf("DOM has changed after developing the program: %s", e.Err)
}
