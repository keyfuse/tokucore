// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"bytes"
	"fmt"

	"github.com/tokublock/tokucore/xbase"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xerror"
)

type mpay struct {
	fee      uint64
	amount   uint64
	locktime uint32
	redeem   []byte
	bondCoin *Coin
	addr     Address
	prv      *xcrypto.PrivateKey
	payerPub *xcrypto.PublicKey
	payeePub *xcrypto.PublicKey
}

// MicroPayer -- payer of micropayment.
type MicroPayer mpay

// MicroPayee -- payee of micropayment.
type MicroPayee mpay

// NewMicroPayer -- creates new payer.
func NewMicroPayer(payer *xcrypto.PrivateKey, payee *xcrypto.PublicKey, amount uint64, fee uint64, locktime uint32) *MicroPayer {
	return &MicroPayer{
		fee:      fee,
		amount:   amount,
		locktime: locktime,
		prv:      payer,
		payerPub: payer.PubKey(),
		payeePub: payee,
		addr:     NewPayToPubKeyHashAddress(payer.PubKey().Hash160()),
	}
}

// Address -- returns the address of MicroPayer.
func (m *MicroPayer) Address() Address {
	return m.addr
}

// CreateBond -- create bond transaction and sign.
func (m *MicroPayer) CreateBond(coin *Coin) (*Transaction, error) {
	bondScript := NewPayToMultiSigScript(2, m.payerPub.Serialize(), m.payeePub.Serialize())
	redeem, err := bondScript.GetLockingScriptBytes()
	if err != nil {
		return nil, err
	}
	m.redeem = redeem

	bond := bondScript.GetAddress()
	bondTx, err := NewTransactionBuilder().
		AddCoins(coin).
		AddKeys(m.prv).
		To(bond, m.amount).
		Then().
		SetChange(m.addr).
		SendFees(m.fee).
		Then().
		Sign().
		BuildTransaction()
	if err != nil {
		return nil, err
	}
	if err := bondTx.Verify(); err != nil {
		return nil, err
	}
	return bondTx, nil
}

// CreateRefund -- create refund transaction without sign.
func (m *MicroPayer) CreateRefund(bondTx *Transaction) (*Transaction, error) {
	bondCoin := NewCoinBuilder().AddOutput(
		bondTx.ID(),
		0,
		m.amount,
		fmt.Sprintf("%x", bondTx.Outputs()[0].Script),
	).ToCoins()[0]
	m.bondCoin = bondCoin

	refundTx, err := NewTransactionBuilder().
		AddCoins(bondCoin).
		SetRedeemScript(m.redeem).
		To(m.addr, m.amount-m.fee).
		Then().
		SendFees(m.fee).
		Then().
		SetLockTime(m.locktime).
		Then().
		BuildTransaction()
	if err != nil {
		return nil, err
	}
	return refundTx, nil
}

// SignRefund -- sign the refund trasaction.
func (m *MicroPayer) SignRefund(refund *Transaction, sign []byte) (*Transaction, error) {
	pubkeys := make([]PubKeySign, 2)

	mysign, err := refund.RawSignature(0, m.prv)
	if err != nil {
		return nil, err
	}
	pubkeys[0] = PubKeySign{PubKey: m.payerPub.Serialize(), Signature: mysign}
	pubkeys[1] = PubKeySign{PubKey: m.payeePub.Serialize(), Signature: sign}
	if err := refund.EmbedIdxSignature(refund.SignIdx(), pubkeys); err != nil {
		return nil, err
	}
	if err := refund.Verify(); err != nil {
		return nil, err
	}
	return refund, nil
}

// SignPayment -- sign the payment transaction and return signature.
func (m *MicroPayer) SignPayment(amount uint64) ([]byte, error) {
	payee := NewPayToPubKeyHashAddress(m.payeePub.Hash160())
	tx, err := NewTransactionBuilder().
		AddCoins(m.bondCoin).
		SetRedeemScript(m.redeem).
		To(payee, amount).
		Then().
		SetChange(m.addr).
		SendFees(m.fee).
		Then().
		BuildTransaction()
	if err != nil {
		return nil, err
	}
	return tx.RawSignature(0, m.prv)
}

// NewMicroPayee -- creates new payee.
func NewMicroPayee(payee *xcrypto.PrivateKey, payer *xcrypto.PublicKey, amount uint64, fee uint64, locktime uint32) *MicroPayee {
	return &MicroPayee{
		fee:      fee,
		amount:   amount,
		locktime: locktime,
		prv:      payee,
		payerPub: payer,
		payeePub: payee.PubKey(),
		addr:     NewPayToPubKeyHashAddress(payee.PubKey().Hash160()),
	}
}

// Address -- returns the address of MicroPayee.
func (m *MicroPayee) Address() Address {
	return m.addr
}

// SignRefund -- sign the refund transaction and return the signature.
func (m *MicroPayee) SignRefund(refund *Transaction) ([]byte, error) {
	// Sanity check.
	if m.locktime != refund.LockTime() {
		return nil, xerror.NewError(Errors, ER_MICROPAYMENT_LOCKTIME_MISMATCH, m.locktime, refund.LockTime())
	}

	sig, err := refund.RawSignature(refund.SignIdx(), m.prv)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

// CheckBond -- check/set bond and refund tx.
func (m *MicroPayee) CheckBond(refund *Transaction, bond *Transaction) error {
	bondScript := NewPayToMultiSigScript(2, m.payerPub.Serialize(), m.payeePub.Serialize())
	redeem, err := bondScript.GetLockingScriptBytes()
	if err != nil {
		return err
	}
	m.redeem = redeem

	input := refund.Inputs()[0]
	bondCoin := NewCoinBuilder().AddOutput(
		xbase.NewIDToString(input.Hash),
		0,
		m.amount,
		fmt.Sprintf("%x", input.PrevLockingScript),
	).ToCoins()[0]
	m.bondCoin = bondCoin

	// Check bond's output is pay to redeem script.
	script, err := NewPayToScriptHashScript(xcrypto.Hash160(redeem)).GetLockingScriptBytes()
	if err != nil || !bytes.Equal(script, bond.Outputs()[0].Script) {
		return xerror.NewError(Errors, ER_MICROPAYMENT_REFUND_BOND_MISMATCH)
	}

	// Check refund input[0]'s hash is bond's.
	if !bytes.Equal(refund.Inputs()[0].Hash, bond.Hash()) {
		return xerror.NewError(Errors, ER_MICROPAYMENT_REFUND_BOND_MISMATCH)
	}
	return nil
}

// SignPayment -- sign the payment transaction with payer signature.
func (m *MicroPayee) SignPayment(amount uint64, sign []byte) (*Transaction, error) {
	payer := NewPayToPubKeyHashAddress(m.payerPub.Hash160())
	tx, err := NewTransactionBuilder().
		AddCoins(m.bondCoin).
		SetRedeemScript(m.redeem).
		To(m.addr, amount).
		Then().
		SetChange(payer).
		SendFees(m.fee).
		Then().
		BuildTransaction()
	if err != nil {
		return nil, err
	}

	pubkeys := make([]PubKeySign, 2)
	sig, err := tx.RawSignature(0, m.prv)
	if err != nil {
		return nil, err
	}
	pubkeys[0] = PubKeySign{PubKey: m.payerPub.Serialize(), Signature: sign}
	pubkeys[1] = PubKeySign{PubKey: m.payeePub.Serialize(), Signature: sig}
	if err := tx.EmbedIdxSignature(0, pubkeys); err != nil {
		return nil, err
	}
	if err := tx.Verify(); err != nil {
		return nil, err
	}
	return tx, nil
}
