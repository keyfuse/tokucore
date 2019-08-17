// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xrpc

import (
	"encoding/json"
)

// GetBlockHash -- get block hash by height.
func (c *Client) GetBlockHash(height int) (string, error) {
	rawResp, err := c.call("getblockhash", height)
	if err != nil {
		return "", err
	}

	resp := ""
	if err = json.Unmarshal(*rawResp.Result, &resp); err != nil {
		return "", err
	}
	return resp, nil
}

// GetBlockCount -- gets block number of blocks in the longest blockchain.
func (c *Client) GetBlockCount() (int, error) {
	rawResp, err := c.call("getblockcount")
	if err != nil {
		return 0, err
	}

	resp := 0
	if err = json.Unmarshal(*rawResp.Result, &resp); err != nil {
		return 0, err
	}
	return resp, nil
}

// GetBlockHeader -- gets block header by block hash.
func (c *Client) GetBlockHeader(hash string) (*BlockHeaderResult, error) {
	rawResp, err := c.call("getblockheader", hash)
	if err != nil {
		return nil, err
	}

	resp := &BlockHeaderResult{}
	if err = json.Unmarshal(*rawResp.Result, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetBlock -- gets verbose block by hash.
func (c *Client) GetBlock(hash string) (*BlockResult, error) {
	// Verbose > 2 is more.
	rawResp, err := c.call("getblock", hash, 3)
	if err != nil {
		return nil, err
	}

	resp := &BlockResult{}
	if err = json.Unmarshal(*rawResp.Result, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
