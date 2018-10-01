// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xerror"
)

func TestTransactionBuilderP2PKH(t *testing.T) {
	msg, _ := DataOutput([]byte("666...satoshi"))

	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohu := bohuHDKey.GetAddress()
	t.Logf("bohu.addr:%v", bohu.ToString(network.TestNet))

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := NewHDKey(seed)
	satoshi := satoshiHDKey.GetAddress()

	// Prepare the UTXOs.
	bohuCoin := NewCoinBuilder().AddOutput(
		"5af1520f1d3e818fca695c2a903baa4a7eec4954f0b35aa01be1f2c1d2cfd802",
		0,
		129990000,
		"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
	).ToCoins()[0]

	tx, err := NewTransactionBuilder().
		AddCoins(bohuCoin).
		AddKeys(bohuPrv).
		To(satoshi, 666666).
		Then().
		SetChange(bohu).
		SendFees(10000).
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
	signedTx := tx.Serialize()
	t.Logf("signed.tx:%x", signedTx)
}

func TestTransactionBuilderWithUncompressedPubKey(t *testing.T) {
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohuPub := bohuPrv.PubKey()

	// Uncompressed pubkey.
	pubHash := xcrypto.Hash160(bohuPub.SerializeUncompressed())
	script := NewPayToPubKeyHashScript(pubHash)
	bohu := script.GetAddress()
	locking, err := script.GetLockingScriptBytes()
	assert.Nil(t, err)

	t.Logf("bohu.addr:%v", bohu.ToString(network.TestNet))

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := NewHDKey(seed)
	satoshi := satoshiHDKey.GetAddress()

	// Prepare the UTXOs.
	bohuCoin := NewCoinBuilder().AddOutput(
		"5af1520f1d3e818fca695c2a903baa4a7eec4954f0b35aa01be1f2c1d2cfd802",
		0,
		129990000,
		fmt.Sprintf("%x", locking),
	).ToCoins()[0]

	tx, err := NewTransactionBuilder().
		AddCoins(bohuCoin).
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

func TestTransactionBuilderMultisig(t *testing.T) {
	seed := []byte("this.is.bohu.seed.")
	bohuHDKey := NewHDKey(seed)
	bohuPrv := bohuHDKey.PrivateKey()
	bohu := bohuHDKey.GetAddress()
	t.Logf("bohu.addr:%v", bohu.ToString(network.TestNet))

	// A.
	seed = []byte("this.is.a.seed.")
	aHDKey := NewHDKey(seed)
	aPrv := aHDKey.PrivateKey()
	aPub := aHDKey.PublicKey().Serialize()

	// B.
	seed = []byte("this.is.b.seed.")
	bHDKey := NewHDKey(seed)
	bPub := bHDKey.PublicKey().Serialize()

	// C.
	seed = []byte("this.is.c.seed.")
	cHDKey := NewHDKey(seed)
	cPrv := cHDKey.PrivateKey()
	cPub := cHDKey.PublicKey().Serialize()

	// Redeem script.
	redeemScript := NewPayToMultiSigScript(2, aPub, bPub, cPub)
	multi := redeemScript.GetAddress()
	t.Logf("multi.addr:%v", multi.ToString(network.TestNet))
	redeem, _ := redeemScript.GetLockingScriptBytes()
	t.Logf("redeem.hex:%x", redeem)

	// Funding.
	{
		bohuCoin := NewCoinBuilder().AddOutput(
			"b470aab9f18259b71fc7cb930929bedb6f6a15f7447219e7216db9a42c782984",
			0,
			129995000,
			"76a9145a927ddadc0ef3ae4501d0d9872b57c9584b9d8888ac",
		).ToCoins()[0]

		tx, err := NewTransactionBuilder().
			AddCoins(bohuCoin).
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
	}

	// Spending.
	{
		multiCoin := NewCoinBuilder().AddOutput(
			"5af1520f1d3e818fca695c2a903baa4a7eec4954f0b35aa01be1f2c1d2cfd802",
			1,
			4000,
			"a914210a461ced66d7540ad2f4649b49dbed7c9fcc2887",
		).ToCoins()[0]

		tx, err := NewTransactionBuilder().
			AddCoins(multiCoin).
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
		err = tx.VerifyDebug()
		assert.Nil(t, err)

		t.Logf("%v", tx.ToString())
		signedTx := tx.Serialize()
		t.Logf("txid:%v", tx.ID())
		t.Logf("spending.signed.tx:%x", signedTx)
	}
}

func TestTransactionBuilderHybrid(t *testing.T) {
	// Alice.
	seed := []byte("this.is.alice.seed.")
	aliceHDKey := NewHDKey(seed)
	alice := aliceHDKey.GetAddress()
	aliceKey := aliceHDKey.PrivateKey()
	aliceCoin := MockP2PKHCoin(aliceHDKey)

	// Bob.
	seed = []byte("this.is.bob.seed.")
	bobHDKey := NewHDKey(seed)
	bobKey := bobHDKey.PrivateKey()
	bobCoin := MockP2PKHCoin(bobHDKey)

	// Alice and bob.
	redeem, _ := NewPayToMultiSigScript(2, aliceHDKey.PublicKey().Serialize(), bobHDKey.PublicKey().Serialize()).GetLockingScriptBytes()
	aliceBobCoin := MockP2SHCoin(aliceHDKey, bobHDKey, redeem)

	// AD.
	pushData, _ := DataOutput([]byte("this.is.pushdata"))

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := NewHDKey(seed)
	satoshi := satoshiHDKey.GetAddress()

	tx, err := NewTransactionBuilder().
		AddCoins(aliceCoin).
		AddKeys(aliceKey).
		To(satoshi, 10000).
		Then().
		AddCoins(bobCoin).
		AddKeys(bobKey).
		To(satoshi, 9000).
		Then().
		AddCoins(aliceBobCoin).
		AddKeys(aliceKey, bobKey).
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
	t.Logf("builder.stats:%#v", tx.Stats())
	t.Logf("signed.string:%v", tx.ToString())
}

func TestTransactionBuilderFees(t *testing.T) {
	// Alice.
	seed := []byte("this.is.alice.seed.")
	aliceHDKey := NewHDKey(seed)
	alice := aliceHDKey.GetAddress()
	aliceKey := aliceHDKey.PrivateKey()
	aliceCoin := MockP2PKHCoin(aliceHDKey)

	// Bob.
	seed = []byte("this.is.bob.seed.")
	bobHDKey := NewHDKey(seed)
	bobKey := bobHDKey.PrivateKey()
	bobCoin := MockP2PKHCoin(bobHDKey)

	// Alice and bob.
	redeem, _ := NewPayToMultiSigScript(2, aliceHDKey.PublicKey().Serialize(), bobHDKey.PublicKey().Serialize()).GetLockingScriptBytes()
	aliceBobCoin := MockP2SHCoin(aliceHDKey, bobHDKey, redeem)

	// AD.
	pushData, _ := DataOutput([]byte("this.is.pushdata"))

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := NewHDKey(seed)
	satoshi := satoshiHDKey.GetAddress()

	tx, err := NewTransactionBuilder().
		AddCoins(aliceCoin).
		AddKeys(aliceKey).
		To(satoshi, 10000).
		Then().
		AddCoins(bobCoin).
		AddKeys(bobKey).
		To(satoshi, 9000).
		Then().
		AddCoins(aliceBobCoin).
		AddKeys(aliceKey, bobKey).
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
	t.Logf("builder.stats:%#v", tx.Stats())
}

func TestTransactionBuilderError(t *testing.T) {
	// Alice.
	seed := []byte("this.is.alice.seed.")
	aliceHDKey := NewHDKey(seed)
	alice := aliceHDKey.GetAddress()
	aliceKey := aliceHDKey.PrivateKey()
	aliceCoin := MockP2PKHCoin(aliceHDKey)

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := NewHDKey(seed)
	satoshi := satoshiHDKey.GetAddress()

	tests := []struct {
		name string
		fn   func() error
		err  error
	}{
		{
			name: "builder.from.nil",
			fn: func() error {
				_, err := NewTransactionBuilder().
					AddKeys(aliceKey).
					To(satoshi, 10000).
					Then().
					SetChange(alice).
					SendFees(1000).
					Then().
					Sign().
					BuildTransaction()
				return err
			},
			err: xerror.NewError(Errors, ER_TRANSACTION_BUILDER_FROM_EMPTY, 0),
		},
		{
			name: "builder.sendto.nil",
			fn: func() error {
				_, err := NewTransactionBuilder().
					AddCoins(aliceCoin).
					AddKeys(aliceKey).
					Then().
					SetChange(alice).
					SendFees(1000).
					Then().
					Sign().
					BuildTransaction()
				return err
			},
			err: xerror.NewError(Errors, ER_TRANSACTION_BUILDER_SENDTO_EMPTY, 0),
		},
		{
			name: "builder.change.nil",
			fn: func() error {
				_, err := NewTransactionBuilder().
					AddCoins(aliceCoin).
					AddKeys(aliceKey).
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
					AddCoins(aliceCoin).
					AddKeys(aliceKey).
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
					AddCoins(aliceCoin).
					AddKeys(aliceKey).
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
					AddCoins(aliceCoin).
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
	}
	for _, test := range tests {
		err := test.fn()
		assert.Equal(t, test.err, err)
	}
}

func BenchmarkTransactionBuilder(b *testing.B) {
	// Alice.
	seed := []byte("this.is.alice.seed.")
	aliceHDKey := NewHDKey(seed)
	alice := aliceHDKey.GetAddress()
	aliceKey := aliceHDKey.PrivateKey()
	aliceCoin := MockP2PKHCoin(aliceHDKey)

	// Bob.
	seed = []byte("this.is.bob.seed.")
	bobHDKey := NewHDKey(seed)
	bobKey := bobHDKey.PrivateKey()
	bobCoin := MockP2PKHCoin(bobHDKey)

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := NewHDKey(seed)
	satoshi := satoshiHDKey.GetAddress()

	for n := 0; n < b.N; n++ {
		_, err := NewTransactionBuilder().
			AddCoins(aliceCoin).
			AddKeys(aliceKey).
			To(satoshi, 5000).
			Then().
			AddCoins(bobCoin).
			AddKeys(bobKey).
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
	aliceHDKey := NewHDKey(seed)
	alice := aliceHDKey.GetAddress()
	aliceKey := aliceHDKey.PrivateKey()
	aliceCoin := MockP2PKHCoin(aliceHDKey)

	// Bob.
	seed = []byte("this.is.bob.seed.")
	bobHDKey := NewHDKey(seed)
	bobKey := bobHDKey.PrivateKey()
	bobCoin := MockP2PKHCoin(bobHDKey)

	// Satoshi.
	seed = []byte("this.is.satoshi.seed.")
	satoshiHDKey := NewHDKey(seed)
	satoshi := satoshiHDKey.GetAddress()

	for n := 0; n < b.N; n++ {
		_, err := NewTransactionBuilder().
			AddCoins(aliceCoin).
			AddKeys(aliceKey).
			To(satoshi, 5000).
			Then().
			AddCoins(bobCoin).
			AddKeys(bobKey).
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
