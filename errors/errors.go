package errors

import "fmt"

// application error codes
const (
	EFileNotFound = "notfound"
	EReadError    = "readerror"
	EWriteError   = "writeerror"
)

// E is our error type
type E struct {
	Code    string
	Message string
	Op      string
	Err     error
}

func (e E) Error() string {
	return fmt.Sprintf("%s: %s/%s + %s", e.Code, e.Message, e.Op, e.Err.Error())
}
