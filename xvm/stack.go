// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"fmt"

	"github.com/tokublock/tokucore/xerror"
)

// Stack -- a stack of immutable objects to be used with bitcoin scripts.
type Stack struct {
	stk [][]byte
}

// asBool -- gets the boolean value of the byte array.
func asBool(t []byte) bool {
	for i := range t {
		if t[i] != 0 {
			// Negative 0 is also considered false.
			if i == len(t)-1 && t[i] == 0x80 {
				return false
			}
			return true
		}
	}
	return false
}

// BoolBytes -- converts a boolean into the appropriate byte array.
func BoolBytes(v bool) []byte {
	if v {
		return []byte{1}
	}
	return nil
}

// NewStack -- create new Stack.
func NewStack() *Stack {
	stack := &Stack{
		stk: make([][]byte, 0),
	}
	return stack
}

// Depth -- returns the number of items on the stack.
func (s *Stack) Depth() int {
	return len(s.stk)
}

// PushByteArray -- adds the given back array to the top of the stack.
// Stack: [... x1 x2] -> [... x1 x2 data]
func (s *Stack) PushByteArray(val []byte) {
	s.stk = append(s.stk, val)
}

// PopByteArray -- pops the value off the top of the stack and returns it.
// Stack: [... x1 x2 x3] -> [... x1 x2]
func (s *Stack) PopByteArray() ([]byte, error) {
	return s.nipN(0)
}

// PushInt -- converts the ScriptNum to a suiteable byte array
// then pushes it onto the top of the stack.
// Stack: [... x1 x2] -> [... x1 x2 int]
func (s *Stack) PushInt(val ScriptNum) {
	s.PushByteArray(val.Bytes())
}

// PopInt -- pops the value off the top of the stack, converts it into a script
// num, and returns it.  The act of converting to a script num enforces the
// consensus rules imposed on data interpreted as numbers.
// Stack: [... x1 x2 x3] -> [... x1 x2]
func (s *Stack) PopInt() (ScriptNum, error) {
	so, err := s.PopByteArray()
	if err != nil {
		return 0, err
	}
	return MakeScriptNum(so, defaultScriptNumMaxLen)
}

// PushBool --
// converts the provided boolean to a suitable byte array then pushes it onto the top of the stack.
// Stack: [... x1 x2] -> [... x1 x2 bool]
func (s *Stack) PushBool(val bool) {
	s.PushByteArray(BoolBytes(val))
}

// PopBool --
// pops the value off the top of the stack, converts it into a bool, and returns it.
// Stack: [... x1 x2 x3] -> [... x1 x2]
func (s *Stack) PopBool() (bool, error) {
	so, err := s.PopByteArray()
	if err != nil {
		return false, err
	}
	return asBool(so), nil
}

// PeekByteArray -- returns the Nth item on the stack without removing it.
func (s *Stack) PeekByteArray(idx int) ([]byte, error) {
	sz := len(s.stk)
	pos := sz - idx - 1
	if idx < 0 || pos < 0 {
		return nil, xerror.NewError(Errors, ER_SCRIPT_STACK_INDEX_INVALID, idx, sz)
	}
	return s.stk[sz-idx-1], nil
}

// nipN -- removes the nth item on the stack and returns it.
// Stack :
// nipN(0): [... x1 x2 x3] -> [... x1 x2]
// nipN(1): [... x1 x2 x3] -> [... x1 x3]
// nipN(2): [... x1 x2 x3] -> [... x2 x3]
func (s *Stack) nipN(idx int) ([]byte, error) {
	sz := len(s.stk)
	pos := sz - idx - 1
	if idx < 0 || pos < 0 {
		return nil, xerror.NewError(Errors, ER_SCRIPT_STACK_INDEX_INVALID, idx, sz)
	}
	so := s.stk[pos]
	s.stk = append(s.stk[:pos], s.stk[pos+1:]...)
	return so, nil
}

// DupN duplicates the top N items on the stack.
//
// Stack :
// DupN(1): [... x1 x2] -> [... x1 x2 x2]
// DupN(2): [... x1 x2] -> [... x1 x2 x1 x2]
func (s *Stack) DupN(n int) error {
	if n < 1 {
		return xerror.NewError(Errors, ER_SCRIPT_STACK_OPERATION_INVALID, "DupN", n)
	}
	// Iteratively duplicate the value n-1 down the stack n times.
	// This leaves an in-order duplicate of the top n items on the stack.
	for i := n; i > 0; i-- {
		so, err := s.PeekByteArray(n - 1)
		if err != nil {
			return err
		}
		s.PushByteArray(so)
	}
	return nil
}

// String -- returns the stack in a readable format.
func (s *Stack) String() string {
	var result string

	if len(s.stk) == 0 {
		result += " <empty> "
	}
	for _, stack := range s.stk {
		if len(stack) == 0 {
			result += " <empty> "
		} else {
			result += fmt.Sprintf(" <%x> ", stack)
		}
	}
	return result
}

// Copy -- clone the stack to new.
func (s *Stack) Copy() *Stack {
	copy := NewStack()
	copy.stk = append(copy.stk, s.stk...)
	return copy
}

// Clean -- clean the stack.
func (s *Stack) Clean() {
	s.stk = s.stk[:0]
}
