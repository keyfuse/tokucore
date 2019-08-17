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

// Demo for sent coin to two-party-threshold P2PKH address and spending from the address to normal address.
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
	shared := xcore.NewPayToPubKeyHashAddress(sharepub.Hash160())
	fmt.Printf("shared.addr:%v\n", shared.ToString(network.TestNet))

	// Funding.
	// https://blockstream.info/testnet/tx/baf5ac09c6047e0f5e082afe37d4a368e40a40a4025ea60057877ea8998f0c60
	{
		bohuCoin := xcore.NewCoinBuilder().AddOutput(
			"12a2e64f5975c29a5f3b3fc79e59461cfa411e3ed0459a3872873265ad9f9979",
			1,
			6718017,
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
	// https://blockstream.info/testnet/tx/5784c155eccfbf0ec636698b99b65c5d32417134443d789edb82bbf877f33f90
	{
		shareCoin := xcore.NewCoinBuilder().AddOutput(
			"baf5ac09c6047e0f5e082afe37d4a368e40a40a4025ea60057877ea8998f0c60",
			0,
			4000,
			"76a914349209f102f196e4a1e4a1eeba1329825ed80e0b88ac",
		).ToCoins()[0]

		tx, err := xcore.NewTransactionBuilder().
			AddCoin(shareCoin).
			To(bohu, 3000).
			Then().
			BuildTransaction()
		assertNil(err)

		// SigHash of index 0.
		idx0sighash := tx.RawSignatureHash(0, xcore.SigHashAll)

		// Phase 2.
		_, _, scalarR1 := aliceParty.Phase2(idx0sighash)
		encpk2, encpub2, scalarR2 := bobParty.Phase2(idx0sighash)

		// Phase 3.
		shareR1 := aliceParty.Phase3(scalarR2)
		shareR2 := bobParty.Phase3(scalarR1)

		// Phase 4.
		sig2, err := bobParty.Phase4(encpk2, encpub2, shareR2)
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
		fmt.Printf("two-party-ecdsa.p2pkh.spend:%v\n", tx.ToString())
		fmt.Printf("two-party-ecdsa.p2pkh.spend.txid:%v\n", tx.ID())
		fmt.Printf("two-party-ecdsa.p2pkh.spend.tx:%x\n", tx.Serialize())
	}
}
