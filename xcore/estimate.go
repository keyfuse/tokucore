// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"github.com/tokublock/tokucore/xbase"
)

const (
	// inputP2PKHSigScriptSize is the worst case (largest) serialize size
	// of a transaction input script that redeems a compressed P2PKH output.
	// It is calculated as:
	//
	//   - OP_DATA_73
	//   - 72 bytes DER signature + 1 byte sighash
	//   - OP_DATA_33
	//   - 33 bytes serialized compressed pubkey
	inputP2PKHSigScriptSize = 1 + 73 + 1 + 33

	// inputP2PKHSize is the worst case (largest) serialize size of a
	// transaction input redeeming a compressed P2PKH output.  It is
	// calculated as:
	//
	//   - 32 bytes previous tx
	//   - 4 bytes output index
	//   - 1 byte compact int encoding value 107
	//   - 107 bytes signature script
	//   - 4 bytes sequence
	inputP2PKHSize = 32 + 4 + 1 + inputP2PKHSigScriptSize + 4

	// Witness datas.
	inputWitnessSignatureSize = 73 // <signature> - 72 bytes DER signature + 1 byte sighash
	inputWitnessPubKeySize    = 33 // <pubkey>    - 33 bytes serialized compressed pubkey
)

// EstimateSize --
// returns a worst case serialize size estimate for a
// signed transaction that spends inputCount number of compressed P2PKH outputs
// and contains each transaction output from txOuts.  The estimated size is
// incremented for an additional P2PKH change output if addChangeOutput is true.
func EstimateSize(txins []*TxIn, txouts []*TxOut) int64 {
	baseSize := 0
	witnessSize := 0

	// Core size.
	{
		// Version.
		baseSize += 4
		// Input size.
		{
			txinNums := len(txins)
			baseSize += xbase.VarIntSerializeSize(uint64(txinNums))
			for _, in := range txins {
				if !in.HasWitness {
					baseSize += inputP2PKHSize
				}
			}
		}

		// Output size.
		{
			txoutNums := len(txouts)
			baseSize += xbase.VarIntSerializeSize(uint64(txoutNums))
			for _, out := range txouts {
				baseSize += (8 + len(out.Script))
			}
		}
		// Locktime.
		baseSize += 4
	}

	// Witness Size.
	{
		var hasWitness bool
		for _, in := range txins {
			if in.HasWitness {
				hasWitness = true
				break
			}
		}

		if hasWitness {
			for _, in := range txins {
				if in.HasWitness {
					// Witness slice varlen.
					witnessSize += xbase.VarIntSerializeSize(uint64(2))
					witnessSize += xbase.VarIntSerializeSize(uint64(inputWitnessSignatureSize))
					witnessSize += xbase.VarIntSerializeSize(uint64(inputWitnessPubKeySize))
				} else {
					witnessSize += xbase.VarIntSerializeSize(uint64(0))
				}
			}
		}
	}
	return int64((baseSize*(witnessScaleFactor-1) + (baseSize + witnessSize)) / witnessScaleFactor)
}

// EstimateFees -- estimate the fee.
func EstimateFees(estimateSize int64, relayFeePerKb int64) int64 {
	fees := 0.0
	if relayFeePerKb > 0 {
		fees = (float64(relayFeePerKb) / float64(1000)) * float64(estimateSize)
	}
	return int64(fees)
}
