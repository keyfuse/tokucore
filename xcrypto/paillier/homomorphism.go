// tokucore
//
// Copyright (c) 2019 Stefano Mozart
// Copyright 2019 by KeyFuse Labs
// BSD License

package paillier

import (
	"fmt"
	"math/big"
)

// MultPlaintext -- returns the ciphertext the will decipher to multiplication
// of the plaintexts (i.e. if ct = Enc(m1), then Dec(MultPlaintext(ct, m2)) = m1 * m2 mod N).
func (pk *PubKey) MultPlaintext(ct *big.Int, msg *big.Int) (*big.Int, error) {
	if ct == nil || ct.Cmp(zero) != 1 {
		return nil, fmt.Errorf("input.invalid")
	}
	m := new(big.Int).Set(msg)
	return new(big.Int).Exp(ct, m, pk.NN), nil
}

// AddPlaintext -- returns the ciphertext the will decipher to addition
// of the plaintexts (i.e if ct = Enc(m1), then Dec(AddPlaintext(ct, m2)) = m1 + m2 mod N)
func (pk *PubKey) AddPlaintext(ct *big.Int, msg *big.Int) (*big.Int, error) {
	if ct == nil || ct.Cmp(zero) != 1 {
		return nil, fmt.Errorf("input.invalid")
	}
	m := new(big.Int).Set(msg)
	ct2 := new(big.Int).Exp(pk.G, m, pk.NN)
	// ct * g^msg mod N^2
	return new(big.Int).Mod(new(big.Int).Mul(ct, ct2), pk.NN), nil
}

// DivPlaintext -- returns the ciphertext the will decipher to division of the plaintexts
// (i.e if ct = Enc(m1), then Dec(DivPlaintext(ct, m2)) = m1 / m2 mod N)
func (pk *PubKey) DivPlaintext(ct *big.Int, msg *big.Int) (*big.Int, error) {
	if ct == nil || ct.Cmp(zero) != 1 {
		return nil, fmt.Errorf("input.invalid")
	}
	m := new(big.Int).Set(msg)
	return new(big.Int).Exp(ct, m.ModInverse(m, pk.NN), pk.NN), nil
}

// Add -- returns a ciphertext `ct3` that will decipher to the sum of
// the corresponding plaintext messages (`m1`, `m2`) ciphered to (`ct1`, `ct2`)
// (i.e if ct1 = Enc(m1) and ct2 = Enc(m2), then Dec(Add(ct1, ct2)) = m1 + m2 mod N)
func (pk *PubKey) Add(ct1, ct2 *big.Int) (*big.Int, error) {
	if ct1 == nil || ct2 == nil || ct1.Cmp(zero) != 1 || ct2.Cmp(zero) != 1 {
		return nil, fmt.Errorf("input.invalid")
	}
	z := new(big.Int).Mul(ct1, ct2)
	return z.Mod(z, pk.NN), nil
}

// Sub -- executes homomorphic subtraction, which corresponds to the addition
// with the modular inverse. That is, it computes a ciphertext ct3 that will
// decipher to the subtration of the corresponding plaintexts. So, if ct1 = Enc(m1)
// and ct2 = Enc(m2), and m1 > m2, then Dec(Sub(ct1, ct2)) = ct1 - ct2 mod N.
// Note that the ciphertext produced by this operation will only make sense if m1>m2.
func (pk *PubKey) Sub(ct1, ct2 *big.Int) *big.Int {
	neg := new(big.Int).ModInverse(ct2, pk.NN)
	neg.Mul(ct1, neg)
	return neg.Mod(neg, pk.NN)
}
