// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xrpc

// BlockHeaderResult -- response for `GetBlockHeader` RPC call.
type BlockHeaderResult struct {
	Hash              string  `json:"hash"`
	Confirmations     int64   `json:"confirmations"`
	Height            int64   `json:"height"`
	Version           int64   `json:"version"`
	VersionHex        string  `json:"versionHex"`
	MerkleRoot        string  `json:"merkleroot"`
	Time              uint64  `json:"time"`
	MedianTime        int64   `json:"mediantime"`
	Nonce             int64   `json:"nonce"`
	Bits              string  `json:"bits"`
	Difficulty        float64 `json:"difficulty"`
	ChainWork         string  `json:"chainwork"`
	PreviousBlockHash []byte  `json:"previousblockhash"`
	NextBlockHash     []byte  `json:"nextblockhash"`
}

// BlockResult --
type BlockResult struct {
	Hash          string              `json:"hash"`
	Confirmations uint64              `json:"confirmations"`
	StrippedSize  int32               `json:"strippedsize"`
	Size          int32               `json:"size"`
	Weight        int32               `json:"weight"`
	Height        int64               `json:"height"`
	Version       int32               `json:"version"`
	VersionHex    string              `json:"versionHex"`
	MerkleRoot    string              `json:"merkleroot"`
	Tx            []TransactionResult `json:"tx,omitempty"`
	NTx           uint32              `json:"nTx"`
	Time          int64               `json:"time"`
	Nonce         uint32              `json:"nonce"`
	Bits          string              `json:"bits"`
	Difficulty    float64             `json:"difficulty"`
	ChainWork     string              `json:"chainwork"`
	PreviousHash  string              `json:"previousblockhash"`
	NextHash      string              `json:"nextblockhash,omitempty"`
}

// ScriptSig --
type ScriptSig struct {
	Asm string `json:"asm"`
	Hex string `json:"hex"`
}

// Vin -- vin.
type Vin struct {
	Coinbase  string    `json:"coinbase"`
	Txid      string    `json:"txid"`
	Vout      int       `json:"vout"`
	ScriptSig ScriptSig `json:"scriptSig"`
	Sequence  uint32    `json:"sequence"`
}

// ScriptPubKey --
type ScriptPubKey struct {
	Asm       string   `json:"asm"`
	Hex       string   `json:"hex"`
	ReqSigs   int      `json:"reqSigs,omitempty"`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses,omitempty"`
}

// Vout -- vout.
type Vout struct {
	Value        float64      `json:"value"`
	N            int          `json:"n"`
	ScriptPubKey ScriptPubKey `json:"scriptPubKey"`
}

// TransactionResult -- a raw transaction.
type TransactionResult struct {
	Hex      string `json:"hex"`
	Txid     string `json:"txid"`
	Version  int32  `json:"version"`
	Size     uint32 `json:"size"`
	Vsize    uint32 `json:"vsize"`
	Weight   uint32 `json:"weight"`
	LockTime uint32 `json:"locktime"`
	Vin      []Vin  `json:"vin"`
	Vout     []Vout `json:"vout"`
}
