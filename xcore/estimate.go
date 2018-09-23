// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"github.com/tokublock/tokucore/xbase"
)

// EstimateSize --
// returns a worst case serialize size estimate for a
// signed transaction that spends inputCount number of compressed P2PKH outputs
// and contains each transaction output from txOuts.  The estimated size is
// incremented for an additional P2PKH change output if addChangeOutput is true.
func EstimateSize(txins []*TxIn, txouts []*TxOut) int64 {
	// inputP2PKHSigScriptSize is the worst case (largest) serialize size
	// of a transaction input script that redeems a compressed P2PKH output.
	// It is calculated as:
	//
	//   - OP_DATA_73
	//   - 72 bytes DER signature + 1 byte sighash
	//   - OP_DATA_33
	//   - 33 bytes serialized compressed pubkey
	inputP2PKHSigScriptSize := 1 + 73 + 1 + 33

	// inputP2PKHSize is the worst case (largest) serialize size of a
	// transaction input redeeming a compressed P2PKH output.  It is
	// calculated as:
	//
	//   - 32 bytes previous tx
	//   - 4 bytes output index
	//   - 1 byte compact int encoding value 107
	//   - 107 bytes signature script
	//   - 4 bytes sequence
	inputP2PKHSize := 32 + 4 + 1 + inputP2PKHSigScriptSize + 4

	size := int64(0)

	// version.
	size += 4

	// input size.
	txinNums := len(txins)
	size += int64(xbase.VarIntSerializeSize(uint64(txinNums)))
	size += int64(inputP2PKHSize * txinNums)

	// output size.
	txoutNums := len(txouts)
	size += int64(xbase.VarIntSerializeSize(uint64(txoutNums)))
	for _, txout := range txouts {
		size += int64(8 + len(txout.Script))
	}
	return size
}

// EstimateFees -- estimate the fee.
func EstimateFees(estimateSize int64, relayFeePerKb int64) int64 {
	fees := int64(0)
	if relayFeePerKb > 0 {
		fees = (relayFeePerKb / 1000) * estimateSize
	}
	return fees
}
