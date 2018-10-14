// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"fmt"
	"strings"

	"github.com/tokublock/tokucore/xbase"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xerror"
	"github.com/tokublock/tokucore/xvm"
)

const (
	hashSize        = 32         // hashSize -- array used to store hashes.
	witnessMarker   = byte(0x00) // Witness Marker fiedl, 0x00.
	witnessFlag     = byte(0x01) // Witness Flag field, 0x01.
	defaultSequence = 0xffffffff // Default sequence.

	// witnessScaleFactor determines the level of "discount" witness data
	// receives compared to "base" data. A scale factor of 4, denotes that
	// witness data is 1/4 as cheap as regular non-witness data.
	witnessScaleFactor = 4
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

	// sigHashMask defines the number of bits of the hash type which is used
	// to identify which outputs are signed.
	sigHashMask = 0x1f
)

// TxIn -- the info of input transaction.
type TxIn struct {
	Hash              []byte   // Previous Tx ID(Hash).
	Index             uint32   // Previous tx output index.
	Value             uint64   // Previous tx output amount.
	Script            []byte   // The final unlocking script.
	RedeemScript      []byte   // Previous redeem script.
	PrevLockingScript []byte   // Previous tx output script(locking script).
	Witness           [][]byte // Witness script.
	HasWitness        bool     // Whether the input is witness.
	WitnessScriptCode []byte   // Witness rework locking script code.
	SignatureHash     []byte   // Signature hash of the input.
	Sequence          uint32
}

// NewTxIn -- build a TxIn.
func NewTxIn(txHash []byte, n uint32, value uint64, script []byte, redeemScript []byte) *TxIn {
	var hasWitness bool
	var witnessScriptCode []byte

	// If the locking script is a witness, we should build the witness script code for signature hash.
	switch GetScriptClass(script) {
	case WitnessV0PubKeyHashTy:
		instrs, err := xvm.NewScriptReader(script).AllInstructions()
		if err != nil {
			return nil
		}
		p2pkh := NewPayToPubKeyHashScript(instrs[1].Data())
		script, err := p2pkh.GetLockingScriptBytes()
		if err != nil {
			return nil
		}
		witnessScriptCode = script
		hasWitness = true
	}
	return &TxIn{
		Hash:              txHash,
		Index:             n,
		Value:             value,
		PrevLockingScript: script,
		RedeemScript:      redeemScript,
		WitnessScriptCode: witnessScriptCode,
		HasWitness:        hasWitness,
		Sequence:          defaultSequence,
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

// Transaction -- the bitcoin transaction.
type Transaction struct {
	signIdx      int
	version      uint32
	lockTime     uint32
	inputs       []*TxIn
	outputs      []*TxOut
	sigHashType  SigHashType
	hashPrevouts []byte
	hashSequence []byte
	hashOutputs  []byte
}

// NewTransaction -- creates a new Transaction.
func NewTransaction() *Transaction {
	return &Transaction{
		version:     1,
		sigHashType: SigHashAll,
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

// SetTxIn -- set the txin tuples.
// Used for deserialize the transaction and verify.
func (tx *Transaction) SetTxIn(idx int, amount uint64, locking []byte, redeem []byte) {
	tx.inputs[idx].Value = amount
	tx.inputs[idx].PrevLockingScript = locking
	tx.inputs[idx].RedeemScript = redeem

	// If the locking script is a witness, we should build the witness script code for signature hash.
	switch GetScriptClass(locking) {
	case WitnessV0PubKeyHashTy:
		instrs, _ := xvm.NewScriptReader(locking).AllInstructions()
		p2pkh := NewPayToPubKeyHashScript(instrs[1].Data())
		script, _ := p2pkh.GetLockingScriptBytes()
		tx.inputs[idx].WitnessScriptCode = script
		tx.inputs[idx].HasWitness = true
	}
}

// LockTime -- returns tx locktime.
func (tx *Transaction) LockTime() uint32 {
	return tx.lockTime
}

// SignIdx -- returns tx signIdx.
func (tx *Transaction) SignIdx() int {
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
	return xcrypto.DoubleSha256(tx.SerializeNoWitness())
}

// ID -- returns transaction hex with reversed format.
func (tx *Transaction) ID() string {
	return xbase.NewIDToString(tx.Hash())
}

// WitnessHash -- returns the tx withess format hash.
func (tx *Transaction) WitnessHash() []byte {
	return xcrypto.DoubleSha256(tx.Serialize())
}

// WitnessID -- returns transaction hex with reversed format.
func (tx *Transaction) WitnessID() string {
	return xbase.NewIDToString(tx.WitnessHash())
}

// SignIndex -- sign specified transaction input with pubkey format.
func (tx *Transaction) SignIndex(idx int, compressed bool, keys ...*xcrypto.PrivateKey) error {
	txIn := tx.inputs[idx]
	signs := make([]PubKeySign, 0)

	// Sanity check.
	if len(keys) > 1 && txIn.RedeemScript == nil {
		return xerror.NewError(Errors, ER_TRANSACTION_SIGN_REDEEM_EMPTY, idx, len(keys))
	}

	for _, key := range keys {
		var err error
		var pubkey []byte
		var signature []byte

		// Pubkey.
		if compressed {
			pubkey = key.PubKey().SerializeCompressed()
		} else {
			pubkey = key.PubKey().SerializeUncompressed()
		}

		if txIn.HasWitness {
			if signature, err = tx.WitnessSignature(idx, key); err != nil {
				return err
			}
		} else {
			if signature, err = tx.RawSignature(idx, key); err != nil {
				return err
			}
		}
		signs = append(signs, PubKeySign{
			PubKey:    pubkey,
			Signature: signature,
		})
	}
	return tx.EmbedIdxSignature(idx, signs)
}

// EmbedIdxSignature -- build the unlocking with signs for the idx.
func (tx *Transaction) EmbedIdxSignature(idx int, signs []PubKeySign) error {
	txIn := tx.inputs[idx]

	// Sanity check.
	if len(signs) > 1 && txIn.RedeemScript == nil {
		return xerror.NewError(Errors, ER_TRANSACTION_SIGN_REDEEM_EMPTY, idx, len(signs))
	}

	builder := xvm.NewScriptBuilder()
	class := GetScriptClass(txIn.PrevLockingScript)
	switch class {
	case PubKeyHashTy:
		// unlocking: <sig> <pubkey>
		// witness:   (empty)
		unlocking, err := builder.AddData(signs[0].Signature).AddData(signs[0].PubKey).Script()
		if err != nil {
			return err
		}
		txIn.Script = unlocking
	case ScriptHashTy:
		// unlocking: OP_0 <A sig> <C sig> <redeemScript>
		// witness:   (empty)
		if len(signs) > 0 {
			builder.AddOp(xvm.OP_0)
		}
		for _, sign := range signs {
			builder.AddData(sign.Signature)
		}
		unlocking, err := builder.AddData(txIn.RedeemScript).Script()
		if err != nil {
			return err
		}
		txIn.Script = unlocking
	case WitnessV0PubKeyHashTy:
		// unlocking: (empty)
		// witness:   [<sig>] [<pubkey>]
		txIn.Witness = append(txIn.Witness, signs[0].Signature)
		txIn.Witness = append(txIn.Witness, signs[0].PubKey)
		txIn.Script = nil
	default:
		return xerror.NewError(Errors, ER_SCRIPT_TYPE_UNKNOWN, xvm.DisasmString(txIn.PrevLockingScript))
	}
	return nil
}

// RawSignatureHash -- returns transaction hash used to get signed/verified.
func (tx *Transaction) RawSignatureHash(idx int, hashType SigHashType) []byte {
	buffer := xbase.NewBuffer()

	// version
	buffer.WriteU32(tx.version)
	switch hashType {
	case SigHashAll:
		buffer.WriteVarInt(uint64(len(tx.inputs)))
		for i, in := range tx.inputs {
			buffer.WriteBytes(in.Hash)
			buffer.WriteU32(in.Index)
			if i == idx {
				if in.RedeemScript != nil {
					buffer.WriteVarBytes(in.RedeemScript)
				} else {
					script := xvm.RemoveOpcode(in.PrevLockingScript, byte(xvm.OP_CODESEPARATOR))
					buffer.WriteVarBytes(script)
				}
			}
			buffer.WriteU32(in.Sequence)
		}

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

// WitnessSignatureHash -- returns transaction hash used to get signed/verified.
func (tx *Transaction) WitnessSignatureHash(idx int, hashType SigHashType) []byte {
	var zeroHash [32]byte
	txIn := tx.inputs[idx]

	// hashPrevouts.
	// If anyone can pay isn't active, then we can use the cached
	// hashPrevOuts, otherwise we just write zeroes for the prev outs.
	if (hashType & SigHashAnyOneCanPay) == 0 {
		if tx.hashPrevouts == nil {
			hashBuffer := xbase.NewBuffer()
			for _, in := range tx.inputs {
				hashBuffer.WriteBytes(in.Hash[:])
				hashBuffer.WriteU32(in.Index)
			}
			tx.hashPrevouts = xcrypto.DoubleSha256(hashBuffer.Bytes())
		}
	} else {
		tx.hashPrevouts = zeroHash[:]
	}

	// If the sighash isn't anyone can pay, single, or none, the use the
	// cached hash sequences, otherwise write all zeroes for the
	// hashSequence.
	if (hashType&SigHashAnyOneCanPay) == 0 && (hashType&sigHashMask) != SigHashSingle && (hashType&sigHashMask) != SigHashNone {
		if tx.hashSequence == nil {
			hashBuffer := xbase.NewBuffer()
			for _, in := range tx.inputs {
				hashBuffer.WriteU32(in.Sequence)
			}
			tx.hashSequence = xcrypto.DoubleSha256(hashBuffer.Bytes())
		}
	} else {
		tx.hashSequence = zeroHash[:]
	}

	// If the current signature mode isn't single, or none, then we can
	// re-use the pre-generated hashoutputs sighash fragment. Otherwise,
	// we'll serialize and add only the target output index to the signature
	// pre-image.
	if hashType&SigHashSingle != SigHashSingle && hashType&SigHashNone != SigHashNone {
		if tx.hashOutputs == nil {
			hashBuffer := xbase.NewBuffer()
			for _, out := range tx.outputs {
				hashBuffer.WriteU64(out.Value)
				hashBuffer.WriteVarBytes(out.Script)
			}
			tx.hashOutputs = xcrypto.DoubleSha256(hashBuffer.Bytes())
		}
	} else {
		tx.hashOutputs = zeroHash[:]
	}

	buffer := xbase.NewBuffer()
	buffer.WriteU32(tx.version)
	buffer.WriteBytes(tx.hashPrevouts)
	buffer.WriteBytes(tx.hashSequence)
	buffer.WriteBytes(txIn.Hash)
	buffer.WriteU32(txIn.Index)
	buffer.WriteVarBytes(txIn.WitnessScriptCode)
	buffer.WriteU64(txIn.Value)
	buffer.WriteU32(txIn.Sequence)
	buffer.WriteBytes(tx.hashOutputs)
	buffer.WriteU32(tx.lockTime)
	buffer.WriteU32(uint32(hashType))
	return xcrypto.DoubleSha256(buffer.Bytes())
}

// RawSignature -- sign the idx input and return the signature.
func (tx *Transaction) RawSignature(idx int, prv *xcrypto.PrivateKey) ([]byte, error) {
	// Sanity Check
	inputs := len(tx.inputs)
	if idx >= inputs {
		return nil, xerror.NewError(Errors, ER_TRANSACTION_SIGN_OUT_INDEX, idx, inputs)
	}

	txIn := tx.inputs[idx]
	txIn.SignatureHash = tx.RawSignatureHash(idx, SigHashAll)
	signature, err := xcrypto.Sign(txIn.SignatureHash, prv)
	if err != nil {
		return nil, err
	}
	return append(signature, byte(tx.sigHashType)), nil
}

// WitnessSignature -- sign the idx input and return the witness signature.
func (tx *Transaction) WitnessSignature(idx int, prv *xcrypto.PrivateKey) ([]byte, error) {
	// Sanity Check
	inputs := len(tx.inputs)
	if idx >= inputs {
		return nil, xerror.NewError(Errors, ER_TRANSACTION_SIGN_OUT_INDEX, idx, inputs)
	}

	txIn := tx.inputs[idx]
	txIn.SignatureHash = tx.WitnessSignatureHash(idx, SigHashAll)
	signature, err := xcrypto.Sign(txIn.SignatureHash, prv)
	if err != nil {
		return nil, err
	}
	return append(signature, byte(tx.sigHashType)), nil
}

// HasWitness -- returns whether the inputs contain witness datas.
func (tx *Transaction) HasWitness() bool {
	for _, in := range tx.inputs {
		if in.HasWitness {
			return true
		}
	}
	return false
}

// Serialize -- the new witness serialization defined in BIP0141 and BIP0144.
func (tx *Transaction) Serialize() []byte {
	buffer := xbase.NewBuffer()
	hasWitness := tx.HasWitness()

	// version
	buffer.WriteU32(tx.version)

	// Witness marker.
	if hasWitness {
		buffer.WriteU8(witnessMarker)
		buffer.WriteU8(witnessFlag)
	}

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

	if hasWitness {
		for _, in := range tx.inputs {
			wits := in.Witness
			buffer.WriteVarInt(uint64(len(wits)))
			for _, wit := range wits {
				buffer.WriteVarBytes(wit)
			}
		}
	}

	// Lock time
	buffer.WriteU32(tx.lockTime)
	return buffer.Bytes()
}

// Deserialize -- decode bytes to witness transaction struct.
func (tx *Transaction) Deserialize(data []byte) error {
	var err error
	buffer := xbase.NewBufferReader(data)

	// Version.
	if tx.version, err = buffer.ReadU32(); err != nil {
		return err
	}

	var marker byte
	if marker, err = buffer.ReadU8(); err != nil {
		return err
	}
	if marker != witnessMarker {
		return fmt.Errorf("witness.marker.error.want:%x.got:%x", witnessMarker, marker)
	}

	var flag byte
	if flag, err = buffer.ReadU8(); err != nil {
		return err
	}
	if flag != witnessFlag {
		return fmt.Errorf("witness.flag.error.want:%x.got:%x", witnessFlag, flag)
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
			if txIn.Hash, err = buffer.ReadBytes(hashSize); err != nil {
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

	// If the transaction's flag byte isn't 0x00 at this point, then one or
	// more of its inputs has accompanying witness data.
	if flag != 0x00 {
		for _, in := range tx.inputs {
			witCount, err := buffer.ReadVarInt()
			if err != nil {
				return err
			}

			in.Witness = make([][]byte, witCount)
			for j := uint64(0); j < witCount; j++ {
				if in.Witness[j], err = buffer.ReadVarBytes(); err != nil {
					return err
				}
			}
		}
	}

	// Lock time.
	if tx.lockTime, err = buffer.ReadU32(); err != nil {
		return err
	}
	return nil
}

// SerializeNoWitness -- normal serialization.
// https://en.bitcoin.it/wiki/Protocol_documentation#tx
func (tx *Transaction) SerializeNoWitness() []byte {
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

// DeserializeNoWitness -- decode bytes to raw transaction struct.
func (tx *Transaction) DeserializeNoWitness(data []byte) error {
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
			if txIn.Hash, err = buffer.ReadBytes(hashSize); err != nil {
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

// Verify -- verify the transaction with signature and pubkey.
func (tx *Transaction) Verify() error {
	for i, in := range tx.inputs {
		engine := xvm.NewEngine()

		// Hasher function.
		hasherFn := func(hashType byte) []byte {
			var sighash []byte
			if in.HasWitness {
				sighash = tx.WitnessSignatureHash(i, SigHashType(hashType))
			} else {
				sighash = tx.RawSignatureHash(i, SigHashType(hashType))
			}
			return sighash
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

		// Locking & Unlocking.
		var locking []byte
		var unlocking []byte
		if in.HasWitness {
			builder := xvm.NewScriptBuilder()
			script, err := builder.AddData(in.Witness[0]).AddData(in.Witness[1]).Script()
			if err != nil {
				return err
			}
			locking = in.WitnessScriptCode
			unlocking = script
		} else {
			locking = in.PrevLockingScript
			unlocking = in.Script
		}

		// Verify.
		err := engine.Verify(unlocking, locking)
		if err != nil {
			return err
		}
	}
	return nil
}

// BaseSize -- the size of the transaction serialised with the witness data stripped.
// https://github.com/bitcoin/bips/blob/master/bip-0141.mediawiki
func (tx *Transaction) BaseSize() int {
	size := 0

	// version
	size += 4

	// inputs
	size += xbase.VarIntSerializeSize(uint64(len(tx.inputs)))
	for _, in := range tx.inputs {
		size += len(in.Hash)
		size += 4
		size += xbase.VarIntSerializeSize(uint64(len(in.Script)))
		size += len(in.Script)
		size += 4
	}

	// outputs
	size += xbase.VarIntSerializeSize(uint64(len(tx.outputs)))
	for _, out := range tx.outputs {
		size += 8
		size += xbase.VarIntSerializeSize(uint64(len(out.Script)))
		size += len(out.Script)
	}

	// Lock time
	size += 4
	return size
}

// WitnessSize -- the witness datas serialised size.
func (tx *Transaction) WitnessSize() int {
	size := 0
	hasWitness := tx.HasWitness()

	if hasWitness {
		for _, in := range tx.inputs {
			wits := in.Witness
			size += xbase.VarIntSerializeSize(uint64(len(wits)))
			for _, wit := range wits {
				size += xbase.VarIntSerializeSize(uint64(len(wit)))
				size += len(wit)
			}
		}
	}
	return size
}

// Weight -- defined as Base transaction size * 3 + Total transaction size.
func (tx *Transaction) Weight() int {
	baseSize := tx.BaseSize()
	witnessSize := tx.WitnessSize()
	return baseSize*(witnessScaleFactor-1) + (baseSize + witnessSize)
}

// Vsize -- defined as Transaction weight / 4 (rounded up to the next integer).
func (tx *Transaction) Vsize() int {
	return tx.Weight() / witnessScaleFactor
}

// Size -- size in bytes serialized as described in BIP144, including base data and witness data.
func (tx *Transaction) Size() int {
	return tx.BaseSize() + tx.WitnessSize()
}

// Fees -- returns the transaction fees.
func (tx *Transaction) Fees() int {
	var totalIn uint64
	var totalOut uint64

	for _, in := range tx.inputs {
		totalIn += in.Value
	}

	for _, out := range tx.outputs {
		totalOut += out.Value
	}
	return int(totalIn - totalOut)
}

// ToString -- returns a human-readable representation of a transaction.
func (tx *Transaction) ToString() string {
	var lines []string

	lines = append(lines, "\n{")
	lines = append(lines, fmt.Sprintf("  \"inputs\":["))
	for i, in := range tx.inputs {
		lines = append(lines, "    {")
		lines = append(lines, fmt.Sprintf("      \"hash\":\t\"%s\",", xbase.NewIDToString(in.Hash)))
		lines = append(lines, fmt.Sprintf("      \"n\":\t%d,", in.Index))
		lines = append(lines, fmt.Sprintf("      \"Value\":\t%d,", in.Value))
		lines = append(lines, fmt.Sprintf("      \"prevlocking\":\t\"%s\",", xvm.DisasmString(in.PrevLockingScript)))
		if in.RedeemScript != nil {
			lines = append(lines, fmt.Sprintf("      \"redeemscript\":\t\"%s\",", xvm.DisasmString(in.RedeemScript)))
			lines = append(lines, fmt.Sprintf("      \"redeemhex\":\t\"%x\",", in.RedeemScript))
		}
		lines = append(lines, fmt.Sprintf("      \"script\":\t\"%s\",", xvm.DisasmString(in.Script)))
		lines = append(lines, fmt.Sprintf("      \"sighash\":\t\"%x\",", in.SignatureHash))
		if in.HasWitness {
			lines = append(lines, fmt.Sprintf("      \"witness\":\t\"%x\",", in.Witness[0]))
			lines = append(lines, fmt.Sprintf("      \"scriptcode\":\t\"%x\"", in.WitnessScriptCode))
		}
		if i == (len(tx.inputs) - 1) {
			lines = append(lines, "    }")
		} else {
			lines = append(lines, "    },")
		}
	}
	lines = append(lines, fmt.Sprintf("  ],"))

	lines = append(lines, fmt.Sprintf("  \"outputs\":["))
	for i, out := range tx.outputs {
		lines = append(lines, "    {")
		lines = append(lines, fmt.Sprintf("      \"value\":\t%d,", out.Value))
		lines = append(lines, fmt.Sprintf("      \"script\":\t\"%s\"", xvm.DisasmString(out.Script)))
		if i == (len(tx.outputs) - 1) {
			lines = append(lines, "    }")
		} else {
			lines = append(lines, "    },")
		}
	}

	fees := tx.Fees()
	size := tx.Size()
	witnessSize := tx.WitnessSize()
	lines = append(lines, fmt.Sprintf("  ],"))
	lines = append(lines, fmt.Sprintf("  \"basesize\":\t%v,", tx.BaseSize()))
	lines = append(lines, fmt.Sprintf("  \"witsize\":\t%v,", witnessSize))
	lines = append(lines, fmt.Sprintf("  \"vsize\":\t%v,", tx.Vsize()))
	lines = append(lines, fmt.Sprintf("  \"size\":\t%v,", size))
	lines = append(lines, fmt.Sprintf("  \"weight\":\t%v,", tx.Weight()))
	lines = append(lines, fmt.Sprintf("  \"fees\":\t%v sat,", fees))
	lines = append(lines, fmt.Sprintf("  \"feesperb\":\t%.2f sat/B,", float64(fees)/float64(size)))
	lines = append(lines, fmt.Sprintf("  \"saving\":\t%.2f %%", (float64(witnessSize)/float64(size))*100))
	lines = append(lines, "}\n")
	return strings.Join(lines, "\n")
}
