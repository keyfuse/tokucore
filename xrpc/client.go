// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Request -- request for RPC.
type Request struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     uint64        `json:"id"`
}

// Error --
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Response -- response for RPC.
type Response struct {
	ID     uint64           `json:"id"`
	Result *json.RawMessage `json:"result"`
	Error  *Error           `json:"error"`
}

// Client -- RPC client.
type Client struct {
	rpcHost    string
	rpcUser    string
	rpcPass    string
	httpClient *http.Client
}

// NewClient -- creates Client.
func NewClient(host string, user string, pass string) *Client {
	return &Client{
		rpcHost:    host,
		rpcUser:    user,
		rpcPass:    pass,
		httpClient: &http.Client{},
	}
}

func (c *Client) call(method string, params ...interface{}) (*Response, error) {
	url := fmt.Sprintf("http://%v/", c.rpcHost)
	param := &Request{
		Method: method,
		Params: params,
	}
	enc, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(enc))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.rpcUser, c.rpcPass)

	// 5s timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	req = req.WithContext(ctx)
	defer cancel()

	rawResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rawResp.Body.Close()

	resp := &Response{}
	if err := json.NewDecoder(rawResp.Body).Decode(resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("%v", resp.Error)
	}
	return resp, nil
}
