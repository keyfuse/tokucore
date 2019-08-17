// tokucore
//
// Copyright (c) 2019 Stefano Mozart
// Copyright 2019 by KeyFuse Labs
// BSD License

package paillier

import (
	"fmt"
	"math/big"

	"crypto/rand"
)

var (
	one  = big.NewInt(1)
	zero = big.NewInt(0)
)

// PubKey -- used to perform encryption and homomorphic operations.
type PubKey struct {
	G  *big.Int `json:"G"`
	N  *big.Int `json:"N"`
	NN *big.Int `json:"NN"`
}

// PrvKey -- used to perform decryption.
type PrvKey struct {
	mu     *big.Int
	pk     *PubKey
	lambda *big.Int
}

// GenerateKeyPair -- returns a Paillier key pair.
func GenerateKeyPair(bitlen int) (*PubKey, *PrvKey, error) {
	p, err := getPrimer(bitlen / 2)
	if err != nil {
		return nil, nil, err
	}
	q, err := getPrimer(bitlen / 2)
	if err != nil {
		return nil, nil, err
	}

	n := new(big.Int).Mul(p, q)
	nn := new(big.Int).Mul(n, n)
	g := new(big.Int).Add(n, one)

	lambda := phi(p, q)
	mu := new(big.Int).ModInverse(lambda, n)

	pk := &PubKey{
		G:  g,
		N:  n,
		NN: nn,
	}

	sk := &PrvKey{
		mu:     mu,
		pk:     pk,
		lambda: lambda,
	}
	return pk, sk, nil
}

// Encrypt -- returns a IND-CPA secure ciphertext for the message `msg`.
func (pk *PubKey) Encrypt(msg *big.Int) (*big.Int, error) {
	m := new(big.Int).Set(msg)
	if m.Cmp(zero) == -1 || m.Cmp(pk.N) != -1 {
		return nil, fmt.Errorf("plaintext.invalid")
	}

	r, err := getRandom(pk.N)
	if err != nil {
		return nil, err
	}
	// c=g^m*r^n (mod n^2)
	r.Exp(r, pk.N, pk.NN)
	m.Exp(pk.G, m, pk.NN)

	c := new(big.Int).Mul(m, r)
	return c.Mod(c, pk.NN), nil
}

// Decrypt -- returns the plaintext corresponding to the ciphertext (ct).
func (sk *PrvKey) Decrypt(ct *big.Int) (*big.Int, error) {
	if ct == nil || ct.Cmp(zero) != 1 {
		return nil, fmt.Errorf("ciphertext.invalid")
	}

	// m = l(c^lambda mod n^2)*mu mod n where L(x) = (x-1)/n
	clambda := ctlambda(ct, sk.lambda, sk.pk.NN)
	m := l(clambda, sk.pk.N)
	m.Mul(m, sk.mu)
	m.Mod(m, sk.pk.N)
	return m, nil
}

// phi -- computes Euler's totient function `φ(p,q) = (p-1)*(q-1)`.
func phi(x, y *big.Int) *big.Int {
	p1 := new(big.Int).Sub(x, one)
	q1 := new(big.Int).Sub(y, one)
	return new(big.Int).Mul(p1, q1)
}

// l -- (x,n) = (x-1)/n is the largest integer quocient `q` to satisfy (x-1) >= q*n.
func l(x, n *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Sub(x, one), n)
}

// ctlambda -- returns c^lambda mod n^2.
func ctlambda(ct, lambda *big.Int, nn *big.Int) *big.Int {
	return new(big.Int).Exp(ct, lambda, nn)
}

func getPrimer(bits int) (*big.Int, error) {
	return rand.Prime(rand.Reader, bits)
}

// TODO(BohuTANG): ensure {gcd(r,n)=1}
// https://en.wikipedia.org/wiki/Paillier_cryptosystem#Encryption
func getRandom(n *big.Int) (*big.Int, error) {
	return rand.Int(rand.Reader, n)
}
