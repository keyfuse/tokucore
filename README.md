# tokucore – A Simple, Powerful Library for Bitcoin Apps.

[![Build Status](https://travis-ci.org/keyfuse/tokucore.png)](https://travis-ci.org/keyfuse/tokucore) [![Go Report Card](https://goreportcard.com/badge/github.com/keyfuse/tokucore)](https://goreportcard.com/report/github.com/keyfuse/tokucore) [![codecov.io](https://codecov.io/gh/keyfuse/tokucore/graphs/badge.svg)](https://codecov.io/gh/keyfuse/tokucore/branch/master) [![BSD License](http://img.shields.io/badge/license-BSD-blue.svg?style=flat)](LICENSE) <img src="http://segwit.co/static/public/images/logo.png" width="100">


## tokucore

*tokucore* is a simple Go (golang) library for creating and manipulating bitcoin data structures like creating keys and addresses (HD/BIP32/BIP39/SegWit) or parsing, creating and signing transactions.

## Overview

* Base58 encoding/decoding
* Block headers, block and transaction parsing
* Transaction creation, signature and verification
* Script parsing and execution
* BIP 32 (deterministic wallets)
* BIP 39 (mnemonic code for generating deterministic keys)
* BIP 173 (Base32 address format for native v0-16 witness outputs)
* Two-Party ECDSA Threshold Signature Scheme (TSS)
* Mult-Party Schnorr Threshold Signature Scheme (TSS)
* Scriptless Adaptor Signature

## Focus

* Simple and easy to use
* No external dependencies
* Full test coverage

## Tests

```
$ export GOPATH=`pwd`
$ go get -u github.com/keyfuse/tokucore/xcore
$ cd src/github.com/keyfuse/tokucore/
$ make test
```

## Examples

- [Generate a Random Address](examples/address_rand.go)
- [Generate a P2PKH Address](examples/address_p2pkh.go)
- [Generate a P2SH Address](examples/address_p2sh.go)
- [Generate a 2-of-3 P2SH MultiSig Address](examples/address_multisig.go)
- [Generate a P2WPKH SegWit Address](examples/address_p2wpkh_v0.go)
- [Generate a P2WSH  SegWit Address](examples/address_p2wsh_v0.go)
- [Create a P2PKH Transaction](examples/transaction_p2pkh.go)
- [Create a 2-to-3 P2SH MultiSig Transaction](examples/transaction_multisig.go)
- [Create a Transaction with an OP_RETURN Output](examples/transaction_opreturn.go)
- [Create a Transaction with Verify](examples/transaction_p2pkh.go)
- [Create a Transaction with P2WPKH Segwit Output](examples/transaction_p2wpkh_v0.go)
- [Create a Transaction with P2WPKH SegWit Input](examples/transaction_p2wpkh_v0.go)
- [Create a Transaction with P2WSH  SegWit Output](examples/transaction_p2wsh_v0.go)
- [Create a Transaction with P2WSH  SegWit Input](examples/transaction_p2wsh_v0.go)
- [Create a Two-Party-Threshold ECDSA Transaction with P2PKH Output](examples/two_party_ecdsa_transaction_p2pkh.go)
- [Create a Two-Party-Threshold ECDSA Transaction with P2PKH Input](examples/two_party_ecdsa_transaction_p2pkh.go)
- [Create a Two-Party-Threshold ECDSA Transaction with P2WPKH SegWit Output](examples/two_party_ecdsa_transaction_p2wpkh.go)
- [Create a Two-Party-Threshold ECDSA Transaction with P2WPKH SegWit Input](examples/two_party_ecdsa_transaction_p2wpkh.go)
- [Scriptless ECDSA adaptor signature](examples/scriptless_ecdsa.go)
- [HDWallet](examples/hdwallet.go)
- [Mnemonic](examples/bip39.go)

## Performance
The performance was done on an Intel i5 (2.50GHz) with golang benchmark.

|    Implementation   |                       |
|---------------------|-----------------------|
| 2P-ECDSA KeyGen     |      2.13 ms          |
| 2P-ECDSA Signing    |      426.05 ms        |
| 2P-Schnorr KeyGen   |      0.01 ms          |
| 2P-Schnorr Signing  |      4.31 ms          |

```
$ go test -bench=BenchmarkMpc*  ./xcrypto

goos: linux
goarch: amd64
pkg: github.com/keyfuse/tokucore/xcrypto
BenchmarkMpcEcdsaKeyGen-4      	    1000	   2133566 ns/op
BenchmarkMpcEcdsaSigning-4     	       3	 426058258 ns/op
BenchmarkMpcSchnorrKeyGen-4    	  100000	     11310 ns/op
BenchmarkMpcSchnorrSigning-4   	     300	   4315631 ns/op
```

## Can I trust this code?
> Don't trust. Verify.

## License

tokucore is released under the BSD License.
