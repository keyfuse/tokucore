// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xerror

import (
	"bytes"
	"fmt"
)

// Error -- the error structure returned from calling.
type Error struct {
	Num     int
	State   string
	Message string
}

// NewError -- creates new Error.
func NewError(table map[int]*Error, number int, args ...interface{}) *Error {
	err := &Error{}
	errn, ok := table[number]
	if !ok {
		return Errors[ER_UNKNOWN]
	}
	err.Num = errn.Num
	err.State = errn.State
	err.Message = fmt.Sprintf(errn.Message, args...)
	return err
}

// Error -- implements the error interface.
func (e *Error) Error() string {
	buf := &bytes.Buffer{}
	buf.WriteString(e.Message)
	fmt.Fprintf(buf, " (errno %d) (state %s)", e.Num, e.State)
	return buf.String()
}

// Error type
const (
	ER_UNKNOWN int = 1000
)

// Errors -- the jump table of error.
var Errors = map[int]*Error{
	ER_UNKNOWN: {Num: ER_UNKNOWN, State: "T0000", Message: "unknown.error"},
}
