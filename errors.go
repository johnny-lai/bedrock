package bedrock

import (
	"fmt"
	"runtime"
)

// Error object
type Error struct {
	Message     string `json:"message"`
	Description string `json:"description"`
}

// Creates a new Error object and sets the message
func Errorf(format string, params ...interface{}) *Error {
	ret := new(Error)
	ret.Message = fmt.Sprintf(format, params...)
  if Mode() == DebugMode {
    ret.Description = StackTrace()
  }
	return ret
}

// This function ensure Error now meets the error interface
func (e *Error) Error() string {
  if e.Description == "" {
    return e.Message
  } else {
    return fmt.Sprintf("%v: %v", e.Message, e.Description)
  }
}

// Sets the Message
func (e *Error) M(msg string) *Error {
	e.Message = msg
	return e
}

// Sets the Description
func (e *Error) D(desc string) *Error {
	e.Description = desc
	return e
}

// Returns the stack trace as a string
func StackTrace() string {
	maxStackTraceSize := 4096

	trace := make([]byte, maxStackTraceSize)
	len := runtime.Stack(trace, false)

	return string(trace[:len])
}
