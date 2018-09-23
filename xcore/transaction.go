// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/tokublock/tokucore/xbase"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xerror"
	"github.com/tokublock/tokucore/xvm"
)

const (
	// HashSize -- array used to store hashes.
	HashSize = 32
)

// SigHashType -- hash type bits at the end of a signature.
type SigHashType uint32

// Hash type bits from the end of a signature.
const (
	SigHashOld          SigHashType = 0x0
	SigHashAll          SigHashType = 0x1
	SigHashNone         SigHashType = 0x2
	SigHashSingle       SigHashType = 0x3
	SigTypeMask         SigHashType = 0x1f
	SigHashAnyOneCanPay SigHashType = 0x80
)

// TxIn -- the info of input transaction.
type TxIn struct {
	Hash              []byte // Previous Tx ID(Hash).
	Index             uint32 // Previous tx output index.
	Script            []byte // Unlocking script.
	RedeemScript      []byte // Previous redeem script.
	PrevLockingScript []byte // Previous tx output script(locking script).
	Sequence          uint32
}

// NewTxIn -- build a TxIn.
func NewTxIn(txHash []byte, n uint32, script []byte, redeemScript []byte) *TxIn {
	return &TxIn{
		Hash:              txHash,
		Index:             n,
		PrevLockingScript: script,
		RedeemScript:      redeemScript,
		Sequence:          0xffffffff,
	}
}

// TxOut -- the info of output transaction.
type TxOut struct {
	Value  uint64
	Script []byte
}

// NewTxOut -- create a new TxOut.
func NewTxOut(value uint64, script []byte) *TxOut {
	return &TxOut{
		Value:  value,
		Script: script,
	}
}

// TransactionStat -- transaction stats from builder.
type TransactionStat struct {
	TotalIn      int64
	TotalOut     int64
	Change       int64
	Fees         int64
	FeesPerKb    int64
	EstimateFees int64
	EstimateSize int64
}

// Transaction -- the bitcoin transaction.
type Transaction struct {
	signIdx     uint32
	version     uint32
	lockTime    uint32
	magic       []byte
	inputs      []*TxIn
	outputs     []*TxOut
	sigHashType SigHashType
	stats       TransactionStat
}

// NewTransaction -- creates a new Transaction.
func NewTransaction() *Transaction {
	return &Transaction{
		version:     1,
		sigHashType: SigHashAll,
		magic:       []byte{0x00, 0x01},
	}
}

// SetSigHashType -- set the sighash type(default SigHashAll).
func (tx *Transaction) SetSigHashType(typ SigHashType) {
	tx.sigHashType = typ
}

// SetVersion -- set the tx version(default 1).
func (tx *Transaction) SetVersion(ver uint32) {
	tx.version = ver
}

// SetLockTime -- set the tx locktime.
func (tx *Transaction) SetLockTime(time uint32) {
	tx.lockTime = time
}

// LockTime -- returns tx locktime.
func (tx *Transaction) LockTime() uint32 {
	return tx.lockTime
}

// SignIdx -- returns tx signIdx.
func (tx *Transaction) SignIdx() uint32 {
	return tx.signIdx
}

// AddInput -- add a TxIn.
func (tx *Transaction) AddInput(in *TxIn) {
	tx.inputs = append(tx.inputs, in)
}

// AddOutput -- add a TxOut.
func (tx *Transaction) AddOutput(out *TxOut) {
	tx.outputs = append(tx.outputs, out)
}

// Outputs -- returns all outputs.
func (tx *Transaction) Outputs() []*TxOut {
	return tx.outputs
}

// Inputs -- returns all outputs.
func (tx *Transaction) Inputs() []*TxIn {
	return tx.inputs
}

// Hash -- returns the tx hash.
func (tx *Transaction) Hash() []byte {
	return xcrypto.DoubleSha256(tx.Serialize())
}

// ID -- returns transaction hex with reversed format.
func (tx *Transaction) ID() string {
	return xbase.NewIDToString(tx.Hash())
}

// Stats -- returns the builder stats.
func (tx *Transaction) Stats() TransactionStat {
	return tx.stats
}

// RawSignature -- sign the idx input and return the signature.
func (tx *Transaction) RawSignature(idx uint32, prv *xcrypto.PrivateKey) ([]byte, error) {
	// Sanity Check
	inputs := uint32(len(tx.inputs))
	if idx >= inputs {
		return nil, xerror.NewError(Errors, ER_TRANSACTION_SIGN_OUT_INDEX, idx, inputs)
	}

	// Signature hash.
	signatureHash := tx.SignatureHash(idx, byte(SigHashAll))
	signature, err := xcrypto.Sign(signatureHash, prv)
	if err != nil {
		return nil, err
	}
	return append(signature, byte(tx.sigHashType)), nil
}

// SignIndex --
// sign specified transaction input.
// If keys more than 1, it will be a multisig.
func (tx *Transaction) SignIndex(idx uint32, keys ...*xcrypto.PrivateKey) error {
	var signs []PubKeySign

	// Sanity check.
	redeemScript := tx.inputs[idx].RedeemScript
	if len(keys) > 1 && redeemScript == nil {
		return xerror.NewError(Errors, ER_TRANSACTION_SIGN_REDEEM_EMPTY, idx, len(keys))
	}

	for _, key := range keys {
		signature, err := tx.RawSignature(idx, key)
		if err != nil {
			return err
		}
		signs = append(signs, PubKeySign{
			PubKey:    key.PubKey().Serialize(),
			Signature: signature,
		})
	}
	return tx.Fill(idx, signs)
}

// Fill -- set the unlocking of the idx input.
func (tx *Transaction) Fill(idx uint32, signs []PubKeySign) error {
	// Sanity check.
	redeemScript := tx.inputs[idx].RedeemScript
	if len(signs) > 1 && redeemScript == nil {
		return xerror.NewError(Errors, ER_TRANSACTION_SIGN_REDEEM_EMPTY, idx, len(signs))
	}

	locking := tx.inputs[idx].PrevLockingScript
	unlocking, err := BuildUnlockingScriptBytes(locking, redeemScript, signs)
	if err != nil {
		return err
	}
	tx.inputs[idx].Script = unlocking
	return nil
}

// SignatureHash --
// returns transaction hashm used to get signed/verified.
func (tx *Transaction) SignatureHash(idx uint32, hashType byte) []byte {
	buffer := xbase.NewBuffer()

	// version
	buffer.WriteU32(tx.version)
	switch SigHashType(hashType) {
	case SigHashAll:
		// inputs.
		buffer.WriteVarInt(uint64(len(tx.inputs)))
		for i, in := range tx.inputs {
			buffer.WriteBytes(in.Hash)
			buffer.WriteU32(in.Index)
			if i == int(idx) {
				if in.RedeemScript != nil {
					buffer.WriteVarBytes(in.RedeemScript)
				} else {
					script, err := xvm.RemoveOpcode(in.PrevLockingScript, byte(xvm.OP_CODESEPARATOR))
					if err != nil {
						panic(err)
					}
					buffer.WriteVarBytes(script)
				}
			}
			buffer.WriteU32(in.Sequence)
		}

		// outputs.
		buffer.WriteVarInt(uint64(len(tx.outputs)))
		for _, out := range tx.outputs {
			buffer.WriteU64(out.Value)
			buffer.WriteVarBytes(out.Script)
		}
	}
	// Lock time
	buffer.WriteU32(tx.lockTime)
	// Hash type.
	buffer.WriteU32(uint32(hashType))
	return xcrypto.DoubleSha256(buffer.Bytes())
}

// Serialize --
// encode the tx to raw format.
// https://en.bitcoin.it/wiki/Protocol_documentation#tx
func (tx *Transaction) Serialize() []byte {
	buffer := xbase.NewBuffer()

	// version
	buffer.WriteU32(tx.version)

	// inputs
	buffer.WriteVarInt(uint64(len(tx.inputs)))
	for _, in := range tx.inputs {
		buffer.WriteBytes(in.Hash)
		buffer.WriteU32(in.Index)
		buffer.WriteVarBytes(in.Script)
		buffer.WriteU32(in.Sequence)
	}

	// outputs
	buffer.WriteVarInt(uint64(len(tx.outputs)))
	for _, out := range tx.outputs {
		buffer.WriteU64(out.Value)
		buffer.WriteVarBytes(out.Script)
	}
	// Lock time
	buffer.WriteU32(tx.lockTime)
	return buffer.Bytes()
}

// Deserialize -- decode bytes to raw transaction struct.
func (tx *Transaction) Deserialize(data []byte) error {
	var err error
	buffer := xbase.NewBufferReader(data)

	// Version.
	if tx.version, err = buffer.ReadU32(); err != nil {
		return err
	}

	// Inputs.
	var ins uint64
	if ins, err = buffer.ReadVarInt(); err != nil {
		return err
	}
	if ins > 0 {
		tx.inputs = make([]*TxIn, ins)
		for i := 0; i < int(ins); i++ {
			txIn := &TxIn{}
			if txIn.Hash, err = buffer.ReadBytes(HashSize); err != nil {
				return err
			}
			if txIn.Index, err = buffer.ReadU32(); err != nil {
				return err
			}
			if txIn.Script, err = buffer.ReadVarBytes(); err != nil {
				return err
			}
			if txIn.Sequence, err = buffer.ReadU32(); err != nil {
				return err
			}
			tx.inputs[i] = txIn
		}
	}

	// Outputs.
	var outs uint64
	if outs, err = buffer.ReadVarInt(); err != nil {
		return err
	}
	if outs > 0 {
		tx.outputs = make([]*TxOut, outs)
		for i := 0; i < int(outs); i++ {
			txOut := &TxOut{}
			if txOut.Value, err = buffer.ReadU64(); err != nil {
				return err
			}
			if txOut.Script, err = buffer.ReadVarBytes(); err != nil {
				return err
			}
			tx.outputs[i] = txOut
		}
	}

	// Lock time.
	if tx.lockTime, err = buffer.ReadU32(); err != nil {
		return err
	}
	return nil

}

// SerializeForPartially -- encode the tx to partially format, include PrevLockingScript and RedeemScript, SignIdx.
func (tx *Transaction) SerializeForPartially(idx uint32) []byte {
	buffer := xbase.NewBuffer()

	// Header.
	buffer.WriteVarBytes(tx.magic)
	buffer.WriteU32(idx)
	buffer.WriteU32(tx.version)

	// inputs
	buffer.WriteVarInt(uint64(len(tx.inputs)))
	for _, in := range tx.inputs {
		buffer.WriteBytes(in.Hash)
		buffer.WriteU32(in.Index)
		buffer.WriteVarBytes(in.PrevLockingScript)
		buffer.WriteVarBytes(in.RedeemScript)
		buffer.WriteU32(in.Sequence)
	}

	// outputs
	buffer.WriteVarInt(uint64(len(tx.outputs)))
	for _, out := range tx.outputs {
		buffer.WriteU64(out.Value)
		buffer.WriteVarBytes(out.Script)
	}
	// Lock time
	buffer.WriteU32(tx.lockTime)
	tx.signIdx = idx
	return buffer.Bytes()
}

// DeserializeForPartially -- decode bytes to transaction struct with PrevLockingScript and RedeemScript.
func (tx *Transaction) DeserializeForPartially(data []byte) error {
	var err error
	buffer := xbase.NewBufferReader(data)

	// Header.
	var magic []byte
	if magic, err = buffer.ReadVarBytes(); err != nil {
		return err
	}
	if !bytes.Equal(tx.magic, magic) {
		return xerror.NewError(Errors, ER_TRANSACTION_PARTIALLY_MAGIC_MISMATCH, tx.magic, magic)
	}
	if tx.signIdx, err = buffer.ReadU32(); err != nil {
		return err
	}
	if tx.version, err = buffer.ReadU32(); err != nil {
		return err
	}

	// Inputs.
	var ins uint64
	if ins, err = buffer.ReadVarInt(); err != nil {
		return err
	}
	if ins > 0 {
		tx.inputs = make([]*TxIn, ins)
		for i := 0; i < int(ins); i++ {
			txIn := &TxIn{}
			if txIn.Hash, err = buffer.ReadBytes(HashSize); err != nil {
				return err
			}
			if txIn.Index, err = buffer.ReadU32(); err != nil {
				return err
			}
			if txIn.PrevLockingScript, err = buffer.ReadVarBytes(); err != nil {
				return err
			}
			if txIn.RedeemScript, err = buffer.ReadVarBytes(); err != nil {
				return err
			}
			if txIn.Sequence, err = buffer.ReadU32(); err != nil {
				return err
			}
			tx.inputs[i] = txIn
		}
	}

	// Outputs.
	var outs uint64
	if outs, err = buffer.ReadVarInt(); err != nil {
		return err
	}
	if outs > 0 {
		tx.outputs = make([]*TxOut, outs)
		for i := 0; i < int(outs); i++ {
			txOut := &TxOut{}
			if txOut.Value, err = buffer.ReadU64(); err != nil {
				return err
			}
			if txOut.Script, err = buffer.ReadVarBytes(); err != nil {
				return err
			}
			tx.outputs[i] = txOut
		}
	}

	// Lock time.
	if tx.lockTime, err = buffer.ReadU32(); err != nil {
		return err
	}
	return nil
}

// Verify -- verify the transaction with signature and pubkey.
func (tx *Transaction) Verify() error {
	return tx.verify(false)
}

// VerifyDebug -- verify the transaction with signature and pubkey.
func (tx *Transaction) VerifyDebug() error {
	return tx.verify(true)
}

// ToString -- returns a human-readable representation of a transaction.
func (tx *Transaction) ToString() string {
	var lines []string

	lines = append(lines, "\n{")
	lines = append(lines, fmt.Sprintf("  \"inputs\":["))
	for _, input := range tx.inputs {
		lines = append(lines, "    {")
		lines = append(lines, fmt.Sprintf("      \"hash\":\t\"%s\"", xbase.NewIDToString(input.Hash)))
		lines = append(lines, fmt.Sprintf("      \"n\":\t%d", input.Index))
		lines = append(lines, fmt.Sprintf("      \"prevlocking\":\t\"%s\"", xvm.DisasmString(input.PrevLockingScript)))
		if input.RedeemScript != nil {
			lines = append(lines, fmt.Sprintf("      \"redeemscript\":\t\"%s\"", xvm.DisasmString(input.RedeemScript)))
			lines = append(lines, fmt.Sprintf("      \"redeemhex\":\t\"%x\"", input.RedeemScript))
		}
		lines = append(lines, fmt.Sprintf("      \"script\":\t\"%s\"", xvm.DisasmString(input.Script)))
		lines = append(lines, "    }")
	}
	lines = append(lines, fmt.Sprintf("  ],"))

	lines = append(lines, fmt.Sprintf("  \"outputs\":["))
	for _, output := range tx.outputs {
		lines = append(lines, "    {")
		lines = append(lines, fmt.Sprintf("      \"value\":\t%d", output.Value))
		lines = append(lines, fmt.Sprintf("      \"script\":\t\"%s\"", xvm.DisasmString(output.Script)))
		lines = append(lines, "    }")
	}
	lines = append(lines, fmt.Sprintf("  ]"))
	lines = append(lines, "}\n")
	return strings.Join(lines, "\n")
}

func (tx *Transaction) verify(debug bool) error {
	for i, in := range tx.inputs {
		engine := xvm.NewEngine()

		// Hasher function.
		hasherFn := func(hashType byte) []byte {
			return tx.SignatureHash(uint32(i), hashType)
		}
		engine.SetHasher(hasherFn)

		// Verifier function.
		verifierFn := func(hash []byte, signature []byte, pubkey []byte) error {
			pub, err := xcrypto.PubKeyFromBytes(pubkey)
			if err != nil {
				return err
			}
			err = xcrypto.Verify(hash, signature, pub)
			return err
		}
		engine.SetVerifier(verifierFn)

		if debug {
			engine.EnableDebug()
		}
		err := engine.Verify(in.Script, in.PrevLockingScript)
		if debug {
			engine.PrintTrace()
		}
		if err != nil {
			return err
		}
	}
	return nil
}
