// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package main

import (
	"fmt"

	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xcore"
	"github.com/keyfuse/tokucore/xcore/bip32"
	"github.com/keyfuse/tokucore/xcrypto"
)

func assertNil(err error) {
	if err != nil {
		panic(err)
	}
}

// Demo for sent coin to two-party-threshold Native SegWit P2WPKH address and spending from the SegWit address to normal address.
func main() {
	// Bohu.
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := bip32.NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohuPub := bohuHDKey.PublicKey()
	bohu := xcore.NewPayToPubKeyHashAddress(bohuPub.Hash160())

	// Alice Party.
	aliceSeed := []byte("this.is.alice.seed.")
	aliceHDKey := bip32.NewHDKey(aliceSeed)
	alicePrv := aliceHDKey.PrivateKey()
	aliceParty := xcrypto.NewEcdsaParty(alicePrv)

	// Bob Party.
	bobSeed := []byte("this.is.bob.seed.")
	bobHDKey := bip32.NewHDKey(bobSeed)
	bobPrv := bobHDKey.PrivateKey()
	bobPub := bobHDKey.PublicKey()
	bobParty := xcrypto.NewEcdsaParty(bobPrv)

	// Phase 1.
	sharepub1 := aliceParty.Phase1(bobPub)
	sharepub := sharepub1

	// Shared address.
	shared := xcore.NewPayToWitnessV0PubKeyHashAddress(sharepub.Hash160())
	fmt.Printf("shared.addr:%v\n", shared.ToString(network.TestNet))

	// Funding.
	// https://blockstream.info/testnet/tx/16b26d2c4c7bf020641bb29eb89aa7c6a7aef7ed078c161c9d8d39453e31057c
	{
		bohuCoin := xcore.NewCoinBuilder().AddOutput(
			"baf5ac09c6047e0f5e082afe37d4a368e40a40a4025ea60057877ea8998f0c60",
			1,
			6713017,
			"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
		).ToCoins()[0]

		tx, err := xcore.NewTransactionBuilder().
			AddCoin(bohuCoin).
			AddKeys(bohuPrv).
			To(shared, 4000).
			Then().
			SetChange(bohu).
			Then().
			Sign().
			BuildTransaction()
		assertNil(err)

		// Verify.
		err = tx.Verify()
		assertNil(err)

		fmt.Printf("p2pkh.fund:%v\n", tx.ToString())
		fmt.Printf("p2pkh.fund.txid:%v\n", tx.ID())
		fmt.Printf("p2pkh.fund.tx:%x\n", tx.Serialize())
	}

	// Spending.
	// https://blockstream.info/testnet/tx/07fa69fe67f3aabe54a5a853b9b051eb3282b6ff673a5a98d277588de63d0fee
	{
		shareCoin := xcore.NewCoinBuilder().AddOutput(
			"16b26d2c4c7bf020641bb29eb89aa7c6a7aef7ed078c161c9d8d39453e31057c",
			0,
			4000,
			"0014349209f102f196e4a1e4a1eeba1329825ed80e0b",
		).ToCoins()[0]

		tx, err := xcore.NewTransactionBuilder().
			AddCoin(shareCoin).
			To(bohu, 3000).
			Then().
			BuildTransaction()
		assertNil(err)

		// Witness SigHash of index 0.
		idx0sighash := tx.WitnessV0SignatureHash(0, xcore.SigHashAll)

		// Phase 2.
		encpk1, encpub1, scalarR1 := aliceParty.Phase2(idx0sighash)
		_, _, scalarR2 := bobParty.Phase2(idx0sighash)

		// Phase 3.
		shareR1 := aliceParty.Phase3(scalarR2)
		shareR2 := bobParty.Phase3(scalarR1)

		// Phase 4.
		sig2, err := bobParty.Phase4(encpk1, encpub1, shareR2)
		assertNil(err)

		// Phase 5.
		fs1, err := aliceParty.Phase5(shareR1, sig2)
		assertNil(err)
		sharesig := fs1

		// EmbedIdxSignature.
		tx.EmbedIdxEcdsaSignature(0, sharepub, sharesig, xcore.SigHashAll)

		// Verify.
		err = tx.Verify()

		fmt.Printf("party1.prvkey:%x\n", alicePrv.Serialize())
		fmt.Printf("party2.prvkey:%x\n", bobPrv.Serialize())
		fmt.Printf("two-party-ecdsa.p2wpkh.spend:%v\n", tx.ToString())
		fmt.Printf("two-party-ecdsa.p2wpkh.spend.txid:%v\n", tx.ID())
		fmt.Printf("two-party-ecdsa.p2wpkh.spend.tx:%x\n", tx.Serialize())
	}
}
