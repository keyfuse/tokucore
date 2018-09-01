// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/tokublock/tokucore/xerror"
)

// Hasher -- hash function for checksig.
type Hasher func(hashType byte) []byte

// Verifier -- verify function for checksig.
type Verifier func(hash []byte, signature []byte, pubkey []byte) error

// Engine -- the virtual matchine to execute the bitcoin scripts.
type Engine struct {
	pc          uint64 // Program counter.
	debug       bool
	dstack      *Stack // Data stack.
	cstack      []int  // Control stack.
	hasher      Hasher
	verifier    Verifier
	reader      *ScriptReader
	instruction *Instruction // Current instruction
	traces      []Trace
}

// NewEngine -- creates new Engine.
func NewEngine() *Engine {
	return &Engine{
		dstack: NewStack(),
	}
}

// EnableDebug -- enable the debug.
func (vm *Engine) EnableDebug() {
	vm.debug = true
}

// DisableDebug -- disable the debug.
func (vm *Engine) DisableDebug() {
	vm.debug = false
}

// SetHasher -- set hasher function.
func (vm *Engine) SetHasher(fn Hasher) {
	vm.hasher = fn
}

// SetVerifier -- set verifier function.
func (vm *Engine) SetVerifier(fn Verifier) {
	vm.verifier = fn
}

// Step --
// will execute the next instruction and move the program counter to the
// next opcode in the script, or the next script if the current has ended.
func (vm *Engine) Step() (bool, error) {
	var err error

	vm.pc++
	if vm.instruction, err = vm.reader.NextInstruction(); err != nil || vm.instruction == nil {
		return true, err
	}
	if vm.branchShouldSkip() && !vm.instruction.isConditional() {
		return false, nil
	}
	if err = vm.instruction.op.opfunc(vm); err != nil {
		return true, err
	}
	return false, nil
}

// Execute -- execute the scripts.
func (vm *Engine) Execute(program []byte) error {
	return vm.execute(program, true)
}

// Verify -- verify the unlocking and locking.
func (vm *Engine) Verify(unlocking []byte, locking []byte) error {
	// Unlocking.
	if err := vm.execute(unlocking, false); err != nil {
		return err
	}
	copyStk := vm.dstack.Copy()

	// Reset.
	vm.reset()
	script := append(unlocking, locking...)
	if err := vm.execute(script, true); err != nil {
		return err
	}

	// P2SH.
	if vm.dstack.Depth() > 1 {
		vm.dstack = copyStk
		redeem, err := vm.dstack.PopByteArray()
		if err != nil {
			return err
		}
		return vm.execute(redeem, true)
	}
	return nil
}

// Traces -- returns the trace records.
func (vm *Engine) Traces() []Trace {
	return vm.traces
}

// PrintTrace -- pretty print the vm trace.
func (vm *Engine) PrintTrace() {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.Debug)
	for i, trace := range vm.traces {
		if i == 0 {
			fmt.Fprintln(w, "\n#Step\tExecuted OP Code\tResulted Stack\tRemaining OP Codes\t")
		}
		row := fmt.Sprintf("%04d\t%s\t%s\t%s\t", trace.Step, trace.Executed, trace.Stack, trace.Remaining)
		fmt.Fprintln(w, row)
	}
	w.Flush()
}

func (vm *Engine) execute(program []byte, final bool) error {
	vm.reader = NewScriptReader(program)
	for {
		done, err := vm.Step()
		if err != nil {
			return err
		}
		if done {
			break
		}
		if vm.debug && final && vm.instruction != nil {
			trace := Trace{
				Step:      vm.pc,
				Executed:  vm.instruction.op.name,
				Stack:     vm.dstack.String(),
				Remaining: vm.reader.DisasmRemaining(),
			}
			vm.traces = append(vm.traces, trace)
		}
	}

	if final {
		if vm.dstack.Depth() > 0 {
			v, err := vm.dstack.PopBool()
			if err != nil {
				return err
			}
			if !v {
				return xerror.NewError(Errors, ER_VM_EXEC_OPCODE_FAILED, "false.stack.entry.at.end.of.script.execution")
			}
		}
	}
	return nil
}

// branchShouldSkip --
// returns whether or not the current conditional branch is skip or not(actively executing).
// For example, when the data stack has an OP_FALSE on it
// and an OP_IF is encountered, the branch is inactive until an OP_ELSE or
// OP_ENDIF is encountered.  It properly handles nested conditionals.
func (vm *Engine) branchShouldSkip() bool {
	if len(vm.cstack) == 0 {
		return false
	}
	return vm.cstack[len(vm.cstack)-1] != OpCondTrue
}

func (vm *Engine) reset() {
	vm.pc = 0
	vm.dstack.Clean()
}
