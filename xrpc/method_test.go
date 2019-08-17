// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xrpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpError(t *testing.T) {
	str := ``
	ts := mockServerHTTPError(str)
	defer ts.Close()

	client := NewClient(ts.URL[7:], "", "")
	_, err := client.GetBlockHeader("00000000b873e79784647a6c82962c70d228557d24a747ea4d1b8bbe878e1206")
	assert.NotNil(t, err)
}

func TestClientError(t *testing.T) {
	str := ``
	ts := mockServerWithError(str)
	defer ts.Close()

	client := NewClient(ts.URL[7:], "", "")
	_, err := client.GetBlockHeader("00000000b873e79784647a6c82962c70d228557d24a747ea4d1b8bbe878e1206")
	assert.NotNil(t, err)
}

func TestGetBlockHash(t *testing.T) {
	str := `"00000000b873e79784647a6c82962c70d228557d24a747ea4d1b8bbe878e1206"`
	ts := mockServer(str)
	defer ts.Close()

	client := NewClient(ts.URL[7:], "", "")
	_, err := client.GetBlockHash(1)
	assert.Nil(t, err)
}

func TestGetBlockCount(t *testing.T) {
	str := "1"
	ts := mockServer(str)
	defer ts.Close()

	client := NewClient(ts.URL[7:], "", "")
	resp, err := client.GetBlockCount()
	assert.Nil(t, err)
	want := 1
	assert.Equal(t, want, resp)
}

func TestGetBlockHeader(t *testing.T) {
	headers := mockBlockHeaders()
	ts := mockServer(headers[0])
	defer ts.Close()

	client := NewClient(ts.URL[7:], "", "")
	_, err := client.GetBlockHeader("00000000b873e79784647a6c82962c70d228557d24a747ea4d1b8bbe878e1206")
	assert.Nil(t, err)
}

func TestGetBlock(t *testing.T) {
	blocks := mockBlocks()
	ts := mockServer(blocks[0])
	defer ts.Close()

	client := NewClient(ts.URL[7:], "", "")
	_, err := client.GetBlock("000000000058b74204bb9d59128e7975b683ac73910660b6531e59523fb4a102")
	assert.Nil(t, err)
}

func BenchmarkGetBlock(b *testing.B) {
	blocks := mockBlocks()
	ts := mockServer(blocks[0])
	defer ts.Close()

	client := NewClient(ts.URL[7:], "", "")
	for n := 0; n < b.N; n++ {
		client.GetBlock("000000000058b74204bb9d59128e7975b683ac73910660b6531e59523fb4a102")
	}
}
