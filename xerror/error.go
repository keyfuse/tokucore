// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xerror

import (
	"bytes"
	"fmt"
	"io"
)

// Error -- the error structure returned from calling.
type Error struct {
	Num     int
	State   string
	Message string
	stack   *stack
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
	err.stack = caller()
	return err
}

// Error -- implements the error interface.
func (e *Error) Error() string {
	buf := &bytes.Buffer{}
	buf.WriteString(e.Message)
	buf.WriteString(fmt.Sprintf(" (errno %d) (state %s)", e.Num, e.State))
	return buf.String()
}

// Format -- implements the error interface.
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, e.Message)
			io.WriteString(s, fmt.Sprintf(" (errno %d) (state %s)\n", e.Num, e.State))
			io.WriteString(s, e.stack.trace())
		} else {
			io.WriteString(s, e.Message)
			io.WriteString(s, fmt.Sprintf(" (errno %d) (state %s)", e.Num, e.State))
		}
	}
}

// Error type
const (
	ER_UNKNOWN int = 1000
)

// Errors -- the jump table of error.
var Errors = map[int]*Error{
	ER_UNKNOWN: {Num: ER_UNKNOWN, State: "T0000", Message: "unknown.error"},
}
