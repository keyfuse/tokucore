// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xcore

import (
	"bytes"
	"math/rand"
	"strconv"
	"strings"

	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"

	"github.com/tokublock/tokucore/network"
	"github.com/tokublock/tokucore/xbase"
	"github.com/tokublock/tokucore/xcrypto"
	"github.com/tokublock/tokucore/xerror"
)

// Formulas:
// pubkey1 = privkey1·G
// n = hash(nonce|i), i=2,3,4...
// pubkey2 = pubkey1 + n2·G = (privkey1 + n2)·G
// privkey2 = (privkey1 + n2)
// https://www.cs.cornell.edu/~iddo/detwal.pdf
//
// BIP32:
// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
const (
	// HardenedKeyStart is the index at which a hardended key starts.  Each
	// extended key has 2^31 normal child keys and 2^31 hardned child keys.
	// Thus the range for normal child keys is [0, 2^31 - 1] and the range
	// for hardened child keys is [2^31, 2^32 - 1].
	HardenedKeyStart = 0x80000000 // 2^31

	// PublicKeyCompressedLength is the byte count of a compressed public key
	PublicKeyCompressedLength = 33

	// serializedKeyLen is the length of a serialized public or private
	// extended key.  It consists of 4 bytes version, 1 byte depth, 4 bytes
	// fingerprint, 4 bytes child number, 32 bytes chain code, and 33 bytes
	// public/private key data.
	serializedKeyLen = 4 + 1 + 4 + 4 + 32 + 33 // 78 bytes
)

// HDKey -- a BIP32 Hierarchically Derived key.
type HDKey struct {
	childNum  []byte // 4 bytes
	parentFP  []byte // 4bytes
	chainCode []byte // 32 bytes
	depth     byte   // 1 bytes
	isPrivate bool   // Unserialized
	prvkey    *xcrypto.PrivateKey
	pubkey    *xcrypto.PublicKey
}

// NewHDKey -- creates a new master HDKey from a seed.
func NewHDKey(seed []byte) *HDKey {
	hmac := hmac.New(sha512.New, []byte("Bitcoin seed"))
	if _, err := hmac.Write(seed); err != nil {
		return nil
	}
	intermediary := hmac.Sum(nil)

	// Spit it into our key and chain code.
	keyBytes := intermediary[:32]
	chainCode := intermediary[32:]
	prvkey := xcrypto.PrvKeyFromBytes(keyBytes)
	return &HDKey{
		childNum:  []byte{0x00, 0x00, 0x00, 0x00},
		parentFP:  []byte{0x00, 0x00, 0x00, 0x00},
		chainCode: chainCode,
		depth:     0x00,
		isPrivate: true,
		prvkey:    prvkey,
	}
}

// NewHDKeyRand -- returns the HDKey with random seed.
func NewHDKeyRand() (*HDKey, error) {
	seed := make([]byte, 256)
	if _, err := rand.Read(seed); err != nil {
		return nil, err
	}
	return NewHDKey(seed), nil
}

// Derive -- returns a derived child extended key at the given index.
// When this extended key is a private extended key (as determined by the IsPrivate
// function), a private extended key will be derived.  Otherwise, the derived
// extended key will be also be a public extended key.
//
// When the index is greater to or equal than the HardenedKeyStart constant, the
// derived extended key will be a hardened extended key.  It is only possible to
// derive a hardended extended key from a private extended key.  Consequently,
// this function will return ErrDeriveHardFromPublic if a hardened child
// extended key is requested from a public extended key.
//
// A hardened extended key is useful since, as previously mentioned, it requires
// a parent private extended key to derive.  In other words, normal child
// extended public keys can be derived from a parent public extended key (no
// knowledge of the parent private key) whereas hardened extended keys may not
// be.
//
// NOTE: There is an extremely small chance (< 1 in 2^127) the specific child
// index does not derive to a usable child.  The ErrInvalidChild error will be
// returned if this should occur, and the caller is expected to ignore the
// invalid child and simply increment to the next index.
func (k *HDKey) Derive(childIdx uint32) (*HDKey, error) {
	// There are four scenarios that could happen here:
	// 1) Private extended key -> Hardened child private extended key
	// 2) Private extended key -> Non-hardened child private extended key
	// 3) Public extended key -> Non-hardened child public extended key
	// 4) Public extended key -> Hardened child public extended key (INVALID!)

	// Case #4 is invalid, so error out early.
	// A hardened child extended key may not be created from a public
	// extended key.
	isChildHardened := childIdx >= HardenedKeyStart
	if !k.isPrivate && isChildHardened {
		return nil, xerror.NewError(Errors, ER_HDKEY_DERIVE_HARD_FROM_PUBLIC)
	}

	// Split "I" into two 32-byte sequences Il and Ir where:
	//   Il = intermediate key used to derive the child
	//   Ir = child chain code
	intermediary, err := k.getIntermediary(childIdx)
	if err != nil {
		return nil, err
	}

	// Create child Key with data common to all both scenarios.
	childKey := &HDKey{
		childNum:  uint32ToBytes(childIdx),
		chainCode: intermediary[32:],
		depth:     k.depth + 1,
		isPrivate: k.isPrivate,
	}

	prvkeyBytes := intermediary[:32]
	if k.isPrivate {
		pubkeyBytes := k.prvkey.PubKey().Serialize()
		fingerprint := xcrypto.Hash160(pubkeyBytes)
		childKey.parentFP = fingerprint[:4]
		childKey.prvkey = k.prvkey.Add(prvkeyBytes)
	} else {
		pubkeyBytes := k.pubkey.Serialize()
		fingerprint := xcrypto.Hash160(pubkeyBytes)
		childKey.parentFP = fingerprint[:4]
		childKey.pubkey = k.pubkey.Add(prvkeyBytes)
	}
	return childKey, nil
}

// DeriveByPath -- derive by path.
// m/1/2/3
func (k *HDKey) DeriveByPath(path string) (*HDKey, error) {
	if path == "" {
		return k, nil
	}

	steps := strings.Split(path, "/")
	if steps[0] != "m" {
		return nil, xerror.NewError(Errors, ER_HDKEY_DERIVE_PATH_INVALID, path)
	}

	hd := k
	for i := 1; i < len(steps); i++ {
		var isHardened bool
		step := steps[i]
		if step[len(step)-1] == '\'' {
			isHardened = true
			step = step[:len(step)-1]
		}

		idx, err := strconv.ParseUint(step, 10, 32)
		if err != nil {
			return nil, xerror.NewError(Errors, ER_HDKEY_DERIVE_PATH_INVALID, path)
		}
		if isHardened {
			idx += HardenedKeyStart
		}
		hd, err = hd.Derive(uint32(idx))
		if err != nil {
			return nil, err
		}
	}
	return hd, nil
}

// getIntermediary --
// get intermediary to create key and chaincode.
// Hardened children are based on the private key.
// NonHardened children are based on public key.
func (k *HDKey) getIntermediary(childIdx uint32) ([]byte, error) {
	var data []byte
	childIndexBytes := uint32ToBytes(childIdx)
	if childIdx >= HardenedKeyStart {
		data = append([]byte{0x0}, k.prvkey.Serialize()...)
	} else {
		if k.isPrivate {
			data = k.prvkey.PubKey().Serialize()
		} else {
			data = k.pubkey.Serialize()
		}
	}
	data = append(data, childIndexBytes...)
	hmac := hmac.New(sha512.New, k.chainCode)
	_, err := hmac.Write(data)
	if err != nil {
		return nil, err
	}
	return hmac.Sum(nil), nil
}

// HDPublicKey -- the public HDkey.
func (k *HDKey) HDPublicKey() *HDKey {
	pubkey := k.pubkey
	if pubkey == nil {
		pubkey = k.prvkey.PubKey()
	}
	return &HDKey{
		depth:     k.depth,
		childNum:  k.childNum,
		parentFP:  k.parentFP,
		chainCode: k.chainCode,
		isPrivate: false,
		pubkey:    pubkey,
	}
}

// PublicKey -- the ecdsa public key.
func (k *HDKey) PublicKey() *xcrypto.PublicKey {
	if k.isPrivate {
		return k.prvkey.PubKey()
	}
	return k.pubkey
}

// PrivateKey -- ecdsa private key.
// If HDKey is public key, it will returns nil.
func (k *HDKey) PrivateKey() *xcrypto.PrivateKey {
	return k.prvkey
}

// ToString -- the HDkey as a human-readable base58-encoded string, WIF, 78bytes.
func (k *HDKey) ToString(net *network.Network) string {
	var keyBytes []byte
	var version []byte

	if k.isPrivate {
		version = net.HDPrivateKeyID
		keyBytes = append([]byte{0x0}, k.prvkey.Serialize()...)
	} else {
		version = net.HDPublicKeyID
		keyBytes = k.pubkey.Serialize()
	}

	// The serialized format is:
	//   version (4) || depth (1) || parent fingerprint (4)) ||
	//   child num (4) || chain code (32) || key data (33) || checksum (4)
	buffer := new(bytes.Buffer)
	buffer.Write(version)
	buffer.WriteByte(k.depth)
	buffer.Write(k.parentFP)
	buffer.Write(k.childNum)
	buffer.Write(k.chainCode)
	buffer.Write(keyBytes)
	data := buffer.Bytes()
	checksum := xcrypto.DoubleSha256(data)[:4]
	datas := append(data, checksum...)
	return xbase.Base58Encode(datas)
}

// GetAddress -- returns the P2PKH address.
func (k *HDKey) GetAddress() Address {
	pubkeyHash := k.PublicKey().Hash160()
	return NewPayToPubKeyHashAddress(pubkeyHash)
}

// NewHDKeyFromString -- import WIF string to HDKey.
func NewHDKeyFromString(wif string) (*HDKey, error) {
	data := xbase.Base58Decode(wif)
	if len(data) != (serializedKeyLen + 4) {
		return nil, xerror.NewError(Errors, ER_HDKEY_SERIALIZED_KEY_WRONG_SIZE)
	}

	// The serialized format is:
	//   version (4) || depth (1) || parent fingerprint (4)) ||
	//   child num (4) || chain code (32) || key data (33) || checksum (4)
	var k = &HDKey{}
	k.depth = data[4]
	k.parentFP = data[5:9]
	k.childNum = data[9:13]
	k.chainCode = data[13:45]

	if data[45] == byte(0) {
		k.isPrivate = true
		k.prvkey = xcrypto.PrvKeyFromBytes(data[46:78])
	} else {
		k.isPrivate = false
		pubkey, err := xcrypto.PubKeyFromBytes(data[45:78])
		if err != nil {
			return nil, err
		}
		k.pubkey = pubkey
	}

	// validate checksum
	cs1 := xcrypto.DoubleSha256(data[0 : len(data)-4])[:4]
	cs2 := data[len(data)-4:]
	for i := range cs1 {
		if cs1[i] != cs2[i] {
			return nil, xerror.NewError(Errors, ER_HDKEY_CHECKSUM_MISMATCH)
		}
	}
	return k, nil
}

func uint32ToBytes(i uint32) []byte {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, i)
	return bytes
}
