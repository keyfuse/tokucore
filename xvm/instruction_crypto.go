// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"fmt"

	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xerror"
)

// Stack:
// [... signature pubkey] -> [... bool]
func opHash160(vm *Engine) error {
	x, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}
	vm.dstack.PushByteArray(xcrypto.Hash160(x))
	return nil
}

// Stack:
// [... signature pubkey] -> [... bool]
func opCheckSig(vm *Engine) error {
	if vm.hasher == nil {
		return xerror.NewError(Errors, ER_VM_EXEC_OPCODE_FAILED, "opCheckSig:vm.hasher.func.is.nil")
	}
	if vm.verifier == nil {
		return xerror.NewError(Errors, ER_VM_EXEC_OPCODE_FAILED, "opCheckSig:vm.verifier.func.is.nil")
	}

	pubkey, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}
	sig, err := vm.dstack.PopByteArray()
	if err != nil {
		return err
	}
	if len(sig) < 1 {
		vm.dstack.PushBool(false)
		return nil
	}

	hashType := sig[len(sig)-1]
	sigDER := sig[:len(sig)-1]
	hash := vm.hasher(hashType)
	if err := vm.verifier(hash, sigDER, pubkey); err != nil {
		vm.dstack.PushBool(false)
		return nil
	}
	vm.dstack.PushBool(true)
	return nil
}

func opCheckSigVerify(vm *Engine) error {
	err := opCheckSig(vm)
	if err != nil {
		return err
	}
	return equalVerify(vm, xerror.NewError(Errors, ER_VM_EXEC_OPCODE_FAILED, fmt.Sprintf("opCheckSigVerify")))
}

// m-of-n multisig
//
// Stack:
// [... dummy [sig ...] numsigs [pubkey ...] numpubkeys] -> [... bool]
func opCheckMultiSig(vm *Engine) error {
	// pubkeys.
	numKeys, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}
	n := int(numKeys.Int32())
	pubKeys := make([][]byte, n)
	for i := 0; i < n; i++ {
		pubKey, err := vm.dstack.PopByteArray()
		if err != nil {
			return err
		}
		pubKeys[i] = pubKey
	}

	// signatures.
	numSignatures, err := vm.dstack.PopInt()
	if err != nil {
		return err
	}
	m := int(numSignatures.Int32())
	if m > n {
		return xerror.NewError(Errors, ER_VM_EXEC_OPCODE_FAILED, fmt.Sprintf("opCheckMultiSig:more.signatures[%v].than.pubkeys[%v]", m, n))
	}
	signatures := make([][]byte, m)
	for i := 0; i < m; i++ {
		signature, err := vm.dstack.PopByteArray()
		if err != nil {
			return err
		}
		signatures[i] = signature
	}

	// A bug causes CHECKMULTISIG to consume one extra argument
	// whose contents were not checked in any way.
	//
	// Unfortunately this is a potential source of mutability,
	// so optionally verify it is exactly equal to zero prior
	// to removing it from the stack.
	_, err = vm.dstack.PopByteArray()
	if err != nil {
		return err
	}

	// Signature verify.
	success := true
	for sigIdx, pubKeyIdx := 0, 0; success && (sigIdx < m) && (pubKeyIdx < n); {
		sig := signatures[sigIdx]
		hashType := sig[len(sig)-1]
		sigDER := sig[:len(sig)-1]

		hash := vm.hasher(hashType)
		if err := vm.verifier(hash, sigDER, pubKeys[pubKeyIdx]); err == nil {
			sigIdx++
		}
		pubKeyIdx++
		if (n - pubKeyIdx) < (m - sigIdx) {
			success = false
		}
	}
	vm.dstack.PushBool(success)
	return nil
}

func opCheckMultiSigVerify(vm *Engine) error {
	err := opCheckMultiSig(vm)
	if err != nil {
		return err
	}
	return equalVerify(vm, xerror.NewError(Errors, ER_VM_EXEC_OPCODE_FAILED, fmt.Sprintf("opCheckMultiSigVerify")))
}
