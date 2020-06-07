package errors

import (
	"bytes"
	"fmt"
)

// application error codes
const (
	ENoError      string = ""
	EInternal     string = "E_INTERNAL"
	EFileNotFound string = "E_NOT_FOUND"
	EReadError    string = "E_READ_ERROR"
	EWriteError   string = "E_WRITE_ERROR"
)

// E is our error type
type E struct {
	Code    string
	Message string
	Context string
	Op      string
	Err     error
}

// Error returns the string representation of the error message.
func (e *E) Error() string {
	var buf bytes.Buffer

	if e.Op != "" {
		fmt.Fprintf(&buf, "%s", e.Op)
		if e.Context != "" {
			buf.WriteString(fmt.Sprintf("/%s", e.Context))
		}
		buf.WriteString(": ")
	}

	if e.Err != nil {
		buf.WriteString(e.Err.Error())
	} else {
		if e.Code != ENoError {
			fmt.Fprintf(&buf, "<%s> ", e.Code)
		}
		buf.WriteString(e.Message)
	}
	return buf.String()
}

// ErrorCode returns relevant error code
func ErrorCode(err error) string {
	if err == nil {
		return ENoError
	}
	e, ok := err.(*E)
	if ok {
		if e.Code != ENoError {
			return e.Code
		} else if e.Err != nil {
			return ErrorCode(e.Err)
		}
	}
	return EInternal
}

// ErrorMessage returns relevant error message
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}
	e, ok := err.(*E)
	if ok {
		if e.Message != "" {
			return e.Message
		} else if e.Err != nil {
			return ErrorMessage(e.Err)
		}
	}
	return fmt.Sprintf("An internal error has occurred: %s", err)
}

// ErrorContext returns relevant error context
func ErrorContext(err error) string {
	if err == nil {
		return ""
	}
	e, ok := err.(*E)
	if ok {
		if e.Context != "" {
			return e.Context
		} else if e.Err != nil {
			return ErrorContext(e.Err)
		}
	}
	return "NoContext"
}
