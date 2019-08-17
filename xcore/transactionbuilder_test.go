// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xcore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/keyfuse/tokucore/network"
	"github.com/keyfuse/tokucore/xcore/bip32"
	"github.com/keyfuse/tokucore/xcrypto"
	"github.com/keyfuse/tokucore/xerror"
)

func TestTransactionBuilderP2PKH(t *testing.T) {
	msg := []byte("666...satoshi")

	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := bip32.NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohuPub := bohuHDKey.PublicKey()
	bohu := NewPayToPubKeyHashAddress(bohuPub.Hash160())
	t.Logf("bohu.addr:%v", bohu.ToString(network.TestNet))

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := bip32.NewHDKey(seed)
	satoshiPub := satoshiHDKey.PublicKey()
	satoshi := NewPayToPubKeyHashAddress(satoshiPub.Hash160())

	// Prepare the UTXOs.
	bohuCoin := NewCoinBuilder().AddOutput(
		"bde974a17f9ab1cfbbfb00bb4561e27156ebd65a4163ea0f014e9114d5b65556",
		1,
		6762017,
		"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
	).ToCoins()[0]

	tx, err := NewTransactionBuilder().
		AddCoin(bohuCoin).
		AddKeys(bohuPrv).
		To(satoshi, 3000).
		Then().
		SetChange(bohu).
		SendFees(1000).
		Then().
		AddPushData(msg).
		Sign().
		BuildTransaction()
	assert.Nil(t, err)

	// Verify.
	err = tx.Verify()
	assert.Nil(t, err)

	t.Logf("%v", tx.ToString())
	t.Logf("txid:%v", tx.ID())

	assert.Equal(t, tx.BaseSize(), tx.Size())
	assert.Equal(t, tx.Vsize(), tx.Size())
	assert.Equal(t, "092ddeb0fa8205a06494f2cf83afda0377479c86065e60dea5ae347468b27361", tx.ID())

	t.Logf("basesize:%+v", tx.BaseSize())
	t.Logf("witnesssize:%+v", tx.WitnessSize())
	t.Logf("vsize:%+v", tx.Vsize())
	t.Logf("weight:%+v", tx.Weight())
	t.Logf("size:%+v", tx.Size())
	signedTx := tx.Serialize()
	t.Logf("actual.size:%v", len(signedTx))
	t.Logf("signed.tx:%x", signedTx)
}

func TestTransactionBuilderP2PKHMultiUTXO(t *testing.T) {
	seed1 := []byte("this.is.bohu.seed.")
	bohuHDKey1 := bip32.NewHDKey(seed1)
	bohuPrv1 := bohuHDKey1.PrivateKey()
	bohuPub1 := bohuHDKey1.PublicKey()
	bohu1 := NewPayToPubKeyHashAddress(bohuPub1.Hash160())
	t.Logf("bohu1.addr:%v", bohu1.ToString(network.TestNet))

	seed2 := []byte("this.is.bohu.seed.2.")
	bohuHDKey2 := bip32.NewHDKey(seed2)
	bohuPrv2 := bohuHDKey2.PrivateKey()
	bohuPub2 := bohuHDKey2.PublicKey()
	bohu2 := NewPayToPubKeyHashAddress(bohuPub2.Hash160())
	t.Logf("bohu2.addr:%v", bohu2.ToString(network.TestNet))

	// Satoshi.
	seed := []byte("this.is.satoshi.seed.")
	satoshiHDKey := bip32.NewHDKey(seed)
	satoshiPub := satoshiHDKey.PublicKey()
	satoshi := NewPayToPubKeyHashAddress(satoshiPub.Hash160())

	// Prepare the UTXOs.
	bohuCoin1 := NewCoinBuilder().AddOutput(
		"c37c3154ae611cfd9a57e684f0c12d51491d09060c643adc292565884e947b2b",
		1,
		126626962,
		"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
	).ToCoins()[0]

	bohuCoin2 := NewCoinBuilder().AddOutput(
		"77fb310add70fdbb37537ae91b07b7ca9e947d2acb48ceb1d4c01b4a5383d4fb",
		0,
		1000000,
		"76a9149a4308b3d5bd509bade50ff4e5fc69833142a85a88ac",
	).ToCoins()[0]

	tx, err := NewTransactionBuilder().
		AddCoin(bohuCoin2).
		AddKeys(bohuPrv2).
		SetSigHashType(SigHashAll).
		Then().
		AddCoin(bohuCoin1).
		AddKeys(bohuPrv1).
		SetSigHashType(SigHashAll).
		Then().
		To(satoshi, 4000).
		Then().
		SetChange(bohu1).
		SendFees(1000).
		Then().
		Sign().
		BuildTransaction()
	assert.Nil(t, err)

	// Verify.
	err = tx.Verify()
	assert.Nil(t, err)

	assert.Equal(t, "4713098984b6d982290f577660577b83df3614aefc68c3fc68b4a5d3b500a480", tx.ID())

	t.Logf("%v", tx.ToString())
	t.Logf("txid:%v", tx.ID())

	signedTx := tx.Serialize()
	t.Logf("signed.tx:%x", signedTx)
}

func TestTransactionBuilderMultisigP2SH(t *testing.T) {
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := bip32.NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohuPub := bohuHDKey.PublicKey()
	bohu := NewPayToPubKeyHashAddress(bohuPub.Hash160())
	t.Logf("bohu.addr:%v", bohu.ToString(network.TestNet))

	// A.
	seed = []byte("this.is.a.seed.")
	aHDKey := bip32.NewHDKey(seed)
	aPrv := aHDKey.PrivateKey()
	aPub := aHDKey.PublicKey().Serialize()

	// B.
	seed = []byte("this.is.b.seed.")
	bHDKey := bip32.NewHDKey(seed)
	bPub := bHDKey.PublicKey().Serialize()

	// C.
	seed = []byte("this.is.c.seed.")
	cHDKey := bip32.NewHDKey(seed)
	cPrv := cHDKey.PrivateKey()
	cPub := cHDKey.PublicKey().Serialize()

	// Redeem script.
	redeemScript := NewPayToMultiSigScript(2, aPub, bPub, cPub)
	redeem, _ := redeemScript.GetLockingScriptBytes()
	t.Logf("redeem.hex:%x", redeem)
	multi := NewPayToScriptHashAddress(xcrypto.Hash160(redeem))
	t.Logf("multi.addr:%v", multi.ToString(network.TestNet))

	// Funding.
	{
		bohuCoin := NewCoinBuilder().AddOutput(
			"092ddeb0fa8205a06494f2cf83afda0377479c86065e60dea5ae347468b27361",
			1,
			6758017,
			"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
		).ToCoins()[0]

		tx, err := NewTransactionBuilder().
			AddCoin(bohuCoin).
			AddKeys(bohuPrv).
			To(multi, 4000).
			Then().
			SetChange(bohu).
			Then().
			Sign().
			BuildTransaction()
		assert.Nil(t, err)

		// Verify.
		err = tx.Verify()
		assert.Nil(t, err)

		t.Logf("%v", tx.ToString())
		t.Logf("txid:%v", tx.ID())
		signedTx := tx.Serialize()
		t.Logf("funding.signed.tx:%x", signedTx)
		assert.Equal(t, "b2e955c95a6ee5752df1477a5936443ead0297ec697475ce6f356cdc6e2301a9", tx.ID())
	}

	// Spending.
	{
		multiCoin := NewCoinBuilder().AddOutput(
			"b2e955c95a6ee5752df1477a5936443ead0297ec697475ce6f356cdc6e2301a9",
			0,
			4000,
			"a914210a461ced66d7540ad2f4649b49dbed7c9fcc2887",
		).ToCoins()[0]

		tx, err := NewTransactionBuilder().
			AddCoin(multiCoin).
			AddKeys(aPrv, cPrv).
			SetRedeemScript(redeem).
			To(bohu, 1000).
			Then().
			SetChange(bohu).
			Then().
			Sign().
			BuildTransaction()
		assert.Nil(t, err)

		// Verify.
		err = tx.Verify()
		assert.Nil(t, err)

		t.Logf("%v", tx.ToString())
		signedTx := tx.Serialize()
		t.Logf("txid:%v", tx.ID())
		t.Logf("spending.signed.tx:%x", signedTx)
		assert.Equal(t, "a28312ed5f5b5d164044f08f3a62e412aeb396043a1ec531c18994ff145ea793", tx.ID())
	}
}

func TestTransactionBuilderP2WPKH(t *testing.T) {
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := bip32.NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohuPub := bohuHDKey.PublicKey()
	bohu := NewPayToPubKeyHashAddress(bohuPub.Hash160())
	t.Logf("bohu.addr:%v", bohu.ToString(network.TestNet))

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := bip32.NewHDKey(seed)
	satoshiPrv := satoshiHDKey.PrivateKey()
	satoshiPubKey := satoshiHDKey.PublicKey()
	satoshi := NewPayToWitnessV0PubKeyHashAddress(satoshiPubKey.Hash160())
	t.Logf("satoshi.p2wpkh.addr:%v", satoshi.ToString(network.TestNet))

	// Funding.
	{
		bohuCoin := NewCoinBuilder().AddOutput(
			"f519a75190312039ddf885231205006b14f2e69f6e5b02314cb0e367b027fa86",
			1,
			127297408,
			"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
		).ToCoins()[0]

		tx, err := NewTransactionBuilder().
			AddCoin(bohuCoin).
			AddKeys(bohuPrv).
			To(satoshi, 666666).
			Then().
			SetChange(bohu).
			SetRelayFeePerKb(20000).
			Then().
			Sign().
			BuildTransaction()
		assert.Nil(t, err)

		// Verify.
		err = tx.Verify()
		assert.Nil(t, err)
		assert.Equal(t, "c37c3154ae611cfd9a57e684f0c12d51491d09060c643adc292565884e947b2b", tx.ID())

		t.Logf("fund:%v", tx.ToString())
		t.Logf("fund.txid:%v", tx.ID())
		t.Logf("fund.tx:%x", tx.Serialize())
		t.Logf("actualsize:%v", len(tx.Serialize()))
	}

	// Spending.
	{
		satoshiCoin := NewCoinBuilder().AddOutput(
			"c37c3154ae611cfd9a57e684f0c12d51491d09060c643adc292565884e947b2b",
			0,
			666666,
			"00148b7f2212ecc4384abcf1df3fc5783e9c2a24d5a5",
		).ToCoins()[0]

		tx, err := NewTransactionBuilder().
			AddCoin(satoshiCoin).
			AddKeys(satoshiPrv).
			To(bohu, 66666).
			Then().
			SetChange(satoshi).
			SetRelayFeePerKb(20000).
			Then().
			Sign().
			BuildTransaction()
		assert.Nil(t, err)

		// Verify.
		err = tx.Verify()
		assert.Nil(t, err)

		t.Logf("spend:%v", tx.ToString())
		t.Logf("spend.txid:%v", tx.ID())
		t.Logf("spend.witnessid:%v", tx.WitnessID())
		t.Logf("spend.tx:%x", tx.Serialize())
		t.Logf("actualsize:%v", len(tx.Serialize()))
		assert.Equal(t, "80cd5fca2589cd97d3da1119214ed339d5284ce068e22f1eb9f32ee99a17d4bf", tx.ID())
	}
}

func TestTransactionBuilderP2WSH(t *testing.T) {
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := bip32.NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohuPub := bohuHDKey.PublicKey()
	bohu := NewPayToPubKeyHashAddress(bohuPub.Hash160())
	t.Logf("bohu.addr:%v", bohu.ToString(network.TestNet))

	// A.
	seed = []byte("this.is.a.seed.")
	aHDKey := bip32.NewHDKey(seed)
	aPrv := aHDKey.PrivateKey()
	aPub := aHDKey.PublicKey().Serialize()

	// B.
	seed = []byte("this.is.b.seed.")
	bHDKey := bip32.NewHDKey(seed)
	bPub := bHDKey.PublicKey().Serialize()

	// C.
	seed = []byte("this.is.c.seed.")
	cHDKey := bip32.NewHDKey(seed)
	cPrv := cHDKey.PrivateKey()
	cPub := cHDKey.PublicKey().Serialize()

	// Redeem script.
	redeemScript := NewPayToMultiSigScript(2, aPub, bPub, cPub)
	redeem, _ := redeemScript.GetLockingScriptBytes()
	t.Logf("redeem.hex:%x", redeem)
	multi := NewPayToWitnessV0ScriptHashAddress(xcrypto.Sha256(redeem))
	t.Logf("multi.addr:%v", multi.ToString(network.TestNet))
	assert.Equal(t, "tb1qrrf2qzw8stxkwhurtamy7wkl3a24vhgu0l3gcf66a8hl5dk9napqtap6rf", multi.ToString(network.TestNet))

	// Funding.
	{
		bohuCoin := NewCoinBuilder().AddOutput(
			"b2e955c95a6ee5752df1477a5936443ead0297ec697475ce6f356cdc6e2301a9",
			1,
			6753017,
			"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
		).ToCoins()[0]

		tx, err := NewTransactionBuilder().
			AddCoin(bohuCoin).
			AddKeys(bohuPrv).
			To(multi, 4000).
			Then().
			SetChange(bohu).
			Then().
			Sign().
			BuildTransaction()
		assert.Nil(t, err)

		// Verify.
		err = tx.Verify()
		assert.Nil(t, err)

		assert.Equal(t, "02f96826dbd8bfec2e88603d110dfef1872809debfd84c12188ab94097da3998", tx.ID())

		t.Logf("%v", tx.ToString())
		t.Logf("txid:%v", tx.ID())
		signedTx := tx.Serialize()
		t.Logf("funding.signed.tx:%x", signedTx)
	}

	// Spending.
	{
		multiCoin := NewCoinBuilder().AddOutput(
			"02f96826dbd8bfec2e88603d110dfef1872809debfd84c12188ab94097da3998",
			0,
			4000,
			"002018d2a009c782cd675f835f764f3adf8f55565d1c7fe28c275ae9effa36c59f42",
		).ToCoins()[0]

		tx, err := NewTransactionBuilder().
			AddCoin(multiCoin).
			AddKeys(aPrv, cPrv).
			SetRedeemScript(redeem).
			To(bohu, 1000).
			Then().
			SetChange(bohu).
			Then().
			Sign().
			BuildTransaction()
		assert.Nil(t, err)

		// Verify.
		err = tx.Verify()
		assert.Nil(t, err)

		t.Logf("%v", tx.ToString())
		signedTx := tx.Serialize()
		t.Logf("txid:%v", tx.ID())
		t.Logf("spending.signed.tx:%x", signedTx)
		assert.Equal(t, "70eaf6275e59e780b933d88ea87b0d1f3135ea2ecb6add971f975155ec80d918", tx.ID())
	}
}

func TestTransactionBuilderWithUncompressedPubKey(t *testing.T) {
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := bip32.NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohuPub := bohuPrv.PubKey()

	// Uncompressed pubkey.
	pubHash := xcrypto.Hash160(bohuPub.SerializeUncompressed())
	script := NewPayToPubKeyHashScript(pubHash)
	bohu := script.GetAddress()
	locking, err := script.GetRawLockingScriptBytes()
	assert.Nil(t, err)

	t.Logf("bohu.addr:%v", bohu.ToString(network.TestNet))

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := bip32.NewHDKey(seed)
	satoshiPub := satoshiHDKey.PublicKey()
	satoshi := NewPayToPubKeyHashAddress(satoshiPub.Hash160())

	// Prepare the UTXOs.
	bohuCoin := NewCoinBuilder().AddOutput(
		"5af1520f1d3e818fca695c2a903baa4a7eec4954f0b35aa01be1f2c1d2cfd802",
		0,
		129990000,
		fmt.Sprintf("%x", locking),
	).ToCoins()[0]

	tx, err := NewTransactionBuilder().
		AddCoin(bohuCoin).
		AddKeys(bohuPrv).
		To(satoshi, 666666).
		SetPubKeyUncompressed().
		Then().
		SetChange(bohu).
		SendFees(10000).
		Then().
		Sign().
		BuildTransaction()
	assert.Nil(t, err)

	// Verify.
	err = tx.Verify()
	assert.Nil(t, err)
}

func TestTransactionBuilderHybrid(t *testing.T) {
	// Alice.
	seed := []byte("this.is.alice.seed.")
	aliceHDKey := bip32.NewHDKey(seed)
	alicePrv := aliceHDKey.PrivateKey()
	alicePub := aliceHDKey.PublicKey()
	alice := NewPayToPubKeyHashAddress(alicePub.Hash160())
	aliceCoin := MockP2PKHCoin(aliceHDKey)

	// Bob.
	seed = []byte("this.is.bob.seed.")
	bobHDKey := bip32.NewHDKey(seed)
	bobPrv := bobHDKey.PrivateKey()
	bobPub := bobHDKey.PublicKey()
	bobCoin := MockP2PKHCoin(bobHDKey)

	// Alice and bob.
	redeem, _ := NewPayToMultiSigScript(2, alicePub.Serialize(), bobPub.Serialize()).GetLockingScriptBytes()
	aliceBobCoin := MockP2SHCoin(aliceHDKey, bobHDKey, redeem)

	// AD.
	pushData := []byte("this.is.pushdata")

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := bip32.NewHDKey(seed)
	satoshiPub := satoshiHDKey.PublicKey()
	satoshi := NewPayToPubKeyHashAddress(satoshiPub.Hash160())

	tx, err := NewTransactionBuilder().
		AddCoin(aliceCoin).
		AddKeys(alicePrv).
		To(satoshi, 10000).
		Then().
		AddCoin(bobCoin).
		AddKeys(bobPrv).
		To(satoshi, 9000).
		Then().
		AddCoin(aliceBobCoin).
		AddKeys(alicePrv, bobPrv).
		SetRedeemScript(redeem).
		To(satoshi, 20000).
		Then().
		SetChange(alice).
		SendFees(1000).
		Then().
		AddPushData(pushData).
		Sign().
		BuildTransaction()
	assert.Nil(t, err)
	signedTx := tx.Serialize()
	err = tx.Verify()
	assert.Nil(t, err)
	t.Logf("signed.hex:%x", signedTx)
	t.Logf("signed.string:%v", tx.ToString())
}

func TestTransactionBuilderFees(t *testing.T) {
	// Alice.
	seed := []byte("this.is.alice.seed.")
	aliceHDKey := bip32.NewHDKey(seed)
	alicePrv := aliceHDKey.PrivateKey()
	alicePub := aliceHDKey.PublicKey()
	alice := NewPayToPubKeyHashAddress(alicePub.Hash160())
	aliceCoin := MockP2PKHCoin(aliceHDKey)

	// Bob.
	seed = []byte("this.is.bob.seed.")
	bobHDKey := bip32.NewHDKey(seed)
	bobPrv := bobHDKey.PrivateKey()
	bobPub := bobHDKey.PublicKey()
	bobCoin := MockP2PKHCoin(bobHDKey)

	// Alice and bob.
	redeem, _ := NewPayToMultiSigScript(2, alicePub.Serialize(), bobPub.Serialize()).GetLockingScriptBytes()
	aliceBobCoin := MockP2SHCoin(aliceHDKey, bobHDKey, redeem)

	// AD.
	pushData := []byte("this.is.pushdata")

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := bip32.NewHDKey(seed)
	satoshiPub := satoshiHDKey.PublicKey()
	satoshi := NewPayToPubKeyHashAddress(satoshiPub.Hash160())

	tx, err := NewTransactionBuilder().
		AddCoin(aliceCoin).
		AddKeys(alicePrv).
		To(satoshi, 10000).
		Then().
		AddCoin(bobCoin).
		AddKeys(bobPrv).
		To(satoshi, 9000).
		Then().
		AddCoin(aliceBobCoin).
		AddKeys(alicePrv, bobPrv).
		SetRedeemScript(redeem).
		To(satoshi, 20000).
		Then().
		SetChange(alice).
		SetRelayFeePerKb(100).
		Then().
		AddPushData(pushData).
		Sign().
		BuildTransaction()
	assert.Nil(t, err)
	signedTx := tx.Serialize()
	t.Logf("actual.size:%v", len(signedTx))
}

func TestTransactionBuilderTSSP2PKH(t *testing.T) {
	// Bohu.
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := bip32.NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohuPub := bohuHDKey.PublicKey()
	bohu := NewPayToPubKeyHashAddress(bohuPub.Hash160())
	t.Logf("bohu.addr:%v", bohu.ToString(network.TestNet))

	// Alice Party.
	aliceSeed := []byte("this.is.alice.seed.")
	aliceHDKey := bip32.NewHDKey(aliceSeed)
	alicePrv := aliceHDKey.PrivateKey()
	alicePub := aliceHDKey.PublicKey()
	aliceParty := xcrypto.NewEcdsaParty(alicePrv)

	// Bob Party.
	bobSeed := []byte("this.is.bob.seed.")
	bobHDKey := bip32.NewHDKey(bobSeed)
	bobPrv := bobHDKey.PrivateKey()
	bobPub := bobHDKey.PublicKey()
	bobParty := xcrypto.NewEcdsaParty(bobPrv)

	// Phase 1.
	sharepub1 := aliceParty.Phase1(bobPub)
	sharepub2 := bobParty.Phase1(alicePub)
	sharepub := sharepub1
	assert.Equal(t, sharepub1, sharepub2)

	// Shared address.
	shared := NewPayToPubKeyHashAddress(sharepub.Hash160())
	t.Logf("shared.addr:%v", shared.ToString(network.TestNet))

	// Funding.
	{
		bohuCoin := NewCoinBuilder().AddOutput(
			"12a2e64f5975c29a5f3b3fc79e59461cfa411e3ed0459a3872873265ad9f9979",
			1,
			6718017,
			"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
		).ToCoins()[0]

		tx, err := NewTransactionBuilder().
			AddCoin(bohuCoin).
			AddKeys(bohuPrv).
			To(shared, 4000).
			Then().
			SetChange(bohu).
			Then().
			Sign().
			BuildTransaction()
		assert.Nil(t, err)

		t.Logf("tx:%v", tx.ToString())
		// Verify.
		err = tx.Verify()
		assert.Nil(t, err)
		assert.Equal(t, "baf5ac09c6047e0f5e082afe37d4a368e40a40a4025ea60057877ea8998f0c60", tx.ID())

		t.Logf("txid:%v", tx.ID())
		signedTx := tx.Serialize()
		t.Logf("funding.signed.tx:%x", signedTx)
	}

	// Spending.
	{
		shareCoin := NewCoinBuilder().AddOutput(
			"baf5ac09c6047e0f5e082afe37d4a368e40a40a4025ea60057877ea8998f0c60",
			0,
			4000,
			"76a914349209f102f196e4a1e4a1eeba1329825ed80e0b88ac",
		).ToCoins()[0]

		tx, err := NewTransactionBuilder().
			AddCoin(shareCoin).
			To(bohu, 3000).
			Then().
			BuildTransaction()
		assert.Nil(t, err)

		idx0sighash := tx.RawSignatureHash(0, SigHashAll)
		t.Logf("idx0.sighash:%x", idx0sighash)

		// Phase 2.
		encpk1, encpub1, scalarR1 := aliceParty.Phase2(idx0sighash)
		encpk2, encpub2, scalarR2 := bobParty.Phase2(idx0sighash)

		// Phase 3.
		shareR1 := aliceParty.Phase3(scalarR2)
		shareR2 := bobParty.Phase3(scalarR1)
		assert.Equal(t, shareR1, shareR2)

		// Phase 4.
		sig1, err := aliceParty.Phase4(encpk2, encpub2, shareR1)
		assert.Nil(t, err)
		sig2, err := bobParty.Phase4(encpk1, encpub1, shareR2)
		assert.Nil(t, err)

		// Phase 5.
		fs1, err := aliceParty.Phase5(shareR1, sig2)
		assert.Nil(t, err)
		fs2, err := bobParty.Phase5(shareR2, sig1)
		assert.Nil(t, err)
		assert.Equal(t, fs1, fs2)
		sharesig := fs1

		// EmbedIdxSignature.
		tx.EmbedIdxEcdsaSignature(0, sharepub, sharesig, SigHashAll)
		t.Logf("tx:%v", tx.ToString())

		// Verify.
		err = tx.Verify()
		assert.Nil(t, err)
		assert.Equal(t, "5784c155eccfbf0ec636698b99b65c5d32417134443d789edb82bbf877f33f90", tx.ID())

		t.Logf("txid:%v", tx.ID())
		signedTx := tx.Serialize()
		t.Logf("spending.signed.tx:%x", signedTx)
	}
}

func TestTransactionBuilderTSSP2WPKH(t *testing.T) {
	// Bohu.
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := bip32.NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohuPub := bohuHDKey.PublicKey()
	bohu := NewPayToPubKeyHashAddress(bohuPub.Hash160())
	t.Logf("bohu.addr:%v", bohu.ToString(network.TestNet))

	// Alice Party.
	aliceSeed := []byte("this.is.alice.seed.")
	aliceHDKey := bip32.NewHDKey(aliceSeed)
	alicePrv := aliceHDKey.PrivateKey()
	alicePub := aliceHDKey.PublicKey()
	aliceParty := xcrypto.NewEcdsaParty(alicePrv)

	// Bob Party.
	bobSeed := []byte("this.is.bob.seed.")
	bobHDKey := bip32.NewHDKey(bobSeed)
	bobPrv := bobHDKey.PrivateKey()
	bobPub := bobHDKey.PublicKey()
	bobParty := xcrypto.NewEcdsaParty(bobPrv)

	// Phase 1.
	sharepub1 := aliceParty.Phase1(bobPub)
	sharepub2 := bobParty.Phase1(alicePub)
	sharepub := sharepub1
	assert.Equal(t, sharepub1, sharepub2)

	// Shared address.
	shared := NewPayToWitnessV0PubKeyHashAddress(sharepub.Hash160())
	t.Logf("shared.addr:%v", shared.ToString(network.TestNet))

	// Funding.
	{
		bohuCoin := NewCoinBuilder().AddOutput(
			"baf5ac09c6047e0f5e082afe37d4a368e40a40a4025ea60057877ea8998f0c60",
			1,
			6713017,
			"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
		).ToCoins()[0]

		tx, err := NewTransactionBuilder().
			AddCoin(bohuCoin).
			AddKeys(bohuPrv).
			To(shared, 4000).
			Then().
			SetChange(bohu).
			Then().
			Sign().
			BuildTransaction()
		assert.Nil(t, err)

		t.Logf("tx:%v", tx.ToString())
		// Verify.
		err = tx.Verify()
		assert.Nil(t, err)
		assert.Equal(t, "16b26d2c4c7bf020641bb29eb89aa7c6a7aef7ed078c161c9d8d39453e31057c", tx.ID())

		t.Logf("txid:%v", tx.ID())
		signedTx := tx.Serialize()
		t.Logf("funding.signed.tx:%x", signedTx)
	}

	// Spending.
	{
		shareCoin := NewCoinBuilder().AddOutput(
			"16b26d2c4c7bf020641bb29eb89aa7c6a7aef7ed078c161c9d8d39453e31057c",
			0,
			4000,
			"0014349209f102f196e4a1e4a1eeba1329825ed80e0b",
		).ToCoins()[0]

		tx, err := NewTransactionBuilder().
			AddCoin(shareCoin).
			To(bohu, 3000).
			Then().
			BuildTransaction()
		assert.Nil(t, err)

		idx0sighash := tx.WitnessV0SignatureHash(0, SigHashAll)
		t.Logf("idx0.sighash:%x", idx0sighash)

		// Phase 2.
		encpk1, encpub1, scalarR1 := aliceParty.Phase2(idx0sighash)
		encpk2, encpub2, scalarR2 := bobParty.Phase2(idx0sighash)

		// Phase 3.
		shareR1 := aliceParty.Phase3(scalarR2)
		shareR2 := bobParty.Phase3(scalarR1)
		assert.Equal(t, shareR1, shareR2)

		// Phase 4.
		sig1, err := aliceParty.Phase4(encpk2, encpub2, shareR1)
		assert.Nil(t, err)
		sig2, err := bobParty.Phase4(encpk1, encpub1, shareR2)
		assert.Nil(t, err)

		// Phase 5.
		fs1, err := aliceParty.Phase5(shareR1, sig2)
		assert.Nil(t, err)
		fs2, err := bobParty.Phase5(shareR2, sig1)
		assert.Nil(t, err)
		assert.Equal(t, fs1, fs2)
		sharesig := fs1

		// EmbedIdxSignature.
		tx.EmbedIdxEcdsaSignature(0, sharepub, sharesig, SigHashAll)
		t.Logf("tx:%v", tx.ToString())

		// Verify.
		err = tx.Verify()
		assert.Nil(t, err)
		assert.Equal(t, "07fa69fe67f3aabe54a5a853b9b051eb3282b6ff673a5a98d277588de63d0fee", tx.ID())

		t.Logf("txid:%v", tx.ID())
		signedTx := tx.Serialize()
		t.Logf("spending.signed.tx:%x", signedTx)
	}
}

func TestTransactionBuilderError(t *testing.T) {
	// Alice.
	seed := []byte("this.is.alice.seed.")
	aliceHDKey := bip32.NewHDKey(seed)
	alicePrv := aliceHDKey.PrivateKey()
	alicePub := aliceHDKey.PublicKey()
	alice := NewPayToPubKeyHashAddress(alicePub.Hash160())
	aliceCoin := MockP2PKHCoin(aliceHDKey)

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := bip32.NewHDKey(seed)
	satoshiPub := satoshiHDKey.PublicKey()
	satoshi := NewPayToPubKeyHashAddress(satoshiPub.Hash160())

	tests := []struct {
		name string
		fn   func() error
		err  error
	}{
		{
			name: "builder.from.nil",
			fn: func() error {
				_, err := NewTransactionBuilder().
					AddKeys(alicePrv).
					To(satoshi, 10000).
					Then().
					SetChange(alice).
					SendFees(1000).
					Then().
					Sign().
					BuildTransaction()
				return err
			},
			err: xerror.NewError(Errors, ER_TRANSACTION_BUILDER_FROM_EMPTY),
		},
		{
			name: "builder.sendto.nil",
			fn: func() error {
				_, err := NewTransactionBuilder().
					AddCoin(aliceCoin).
					AddKeys(alicePrv).
					Then().
					SetChange(alice).
					SendFees(1000).
					Then().
					Sign().
					BuildTransaction()
				return err
			},
			err: xerror.NewError(Errors, ER_TRANSACTION_BUILDER_SENDTO_EMPTY),
		},
		{
			name: "builder.change.nil",
			fn: func() error {
				_, err := NewTransactionBuilder().
					AddCoin(aliceCoin).
					AddKeys(alicePrv).
					To(satoshi, 1000).
					Then().
					SendFees(1000).
					Then().
					Sign().
					BuildTransaction()
				return err
			},
			err: xerror.NewError(Errors, ER_TRANSACTION_BUILDER_CHANGETO_EMPTY),
		},
		{
			name: "builder.fee.not.enough",
			fn: func() error {
				_, err := NewTransactionBuilder().
					AddCoin(aliceCoin).
					AddKeys(alicePrv).
					To(satoshi, 10000).
					Then().
					SetChange(alice).
					SendFees(1000).
					Then().
					Sign().
					BuildTransaction()
				return err
			},
			err: xerror.NewError(Errors, ER_TRANSACTION_BUILDER_MIN_FEE_NOT_ENOUGH, 1000, 0),
		},
		{
			name: "builder.totalout.more.than.totalin",
			fn: func() error {
				_, err := NewTransactionBuilder().
					AddCoin(aliceCoin).
					AddKeys(alicePrv).
					To(satoshi, 1000000).
					Then().
					SetChange(alice).
					SendFees(1000).
					Then().
					Sign().
					BuildTransaction()
				return err
			},
			err: xerror.NewError(Errors, ER_TRANSACTION_BUILDER_AMOUNT_NOT_ENOUGH_ERROR, 1000000, 10000),
		},
		{
			name: "builder.keys.nil",
			fn: func() error {
				_, err := NewTransactionBuilder().
					AddCoin(aliceCoin).
					To(satoshi, 1000).
					Then().
					SetChange(alice).
					SendFees(1000).
					Then().
					Sign().
					BuildTransaction()
				return err
			},
			err: xerror.NewError(Errors, ER_TRANSACTION_BUILDER_SIGN_KEY_EMPTY, 0),
		},
		{
			name: "builder.fee.high",
			fn: func() error {
				_, err := NewTransactionBuilder().
					AddCoin(aliceCoin).
					To(satoshi, 1000).
					Then().
					SetChange(alice).
					SetMaxFees(10).
					SetRelayFeePerKb(1000).
					Then().
					Sign().
					BuildTransaction()
				return err
			},
			err: xerror.NewError(Errors, ER_TRANSACTION_BUILDER_FEE_TOO_HIGH, 192, 10),
		},
	}
	for _, test := range tests {
		err := test.fn()
		assert.Equal(t, test.err.Error(), err.Error())
	}
}

func BenchmarkTransactionBuilder(b *testing.B) {
	// Alice.
	seed := []byte("this.is.alice.seed.")
	aliceHDKey := bip32.NewHDKey(seed)
	alicePrv := aliceHDKey.PrivateKey()
	alicePub := aliceHDKey.PublicKey()
	alice := NewPayToPubKeyHashAddress(alicePub.Hash160())
	aliceCoin := MockP2PKHCoin(aliceHDKey)

	// Bob.
	seed = []byte("this.is.bob.seed.")
	bobHDKey := bip32.NewHDKey(seed)
	bobPrv := bobHDKey.PrivateKey()
	bobCoin := MockP2PKHCoin(bobHDKey)

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := bip32.NewHDKey(seed)
	satoshiPub := satoshiHDKey.PublicKey()
	satoshi := NewPayToPubKeyHashAddress(satoshiPub.Hash160())

	for n := 0; n < b.N; n++ {
		_, err := NewTransactionBuilder().
			AddCoin(aliceCoin).
			AddKeys(alicePrv).
			To(satoshi, 5000).
			Then().
			AddCoin(bobCoin).
			AddKeys(bobPrv).
			To(satoshi, 5000).
			Then().
			SetChange(alice).
			Then().
			BuildTransaction()
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkTransactionBuilderSigned(b *testing.B) {
	// Alice.
	seed := []byte("this.is.alice.seed.")
	aliceHDKey := bip32.NewHDKey(seed)
	alicePrv := aliceHDKey.PrivateKey()
	alicePub := aliceHDKey.PublicKey()
	alice := NewPayToPubKeyHashAddress(alicePub.Hash160())
	aliceCoin := MockP2PKHCoin(aliceHDKey)

	// Bob.
	seed = []byte("this.is.bob.seed.")
	bobHDKey := bip32.NewHDKey(seed)
	bobPrv := bobHDKey.PrivateKey()
	bobCoin := MockP2PKHCoin(bobHDKey)

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := bip32.NewHDKey(seed)
	satoshiPub := satoshiHDKey.PublicKey()
	satoshi := NewPayToPubKeyHashAddress(satoshiPub.Hash160())

	for n := 0; n < b.N; n++ {
		_, err := NewTransactionBuilder().
			AddCoin(aliceCoin).
			AddKeys(alicePrv).
			To(satoshi, 5000).
			Then().
			AddCoin(bobCoin).
			AddKeys(bobPrv).
			To(satoshi, 5000).
			Then().
			SetChange(alice).
			Then().
			Sign().
			BuildTransaction()
		if err != nil {
			panic(err)
		}
	}
}
