// tokucore
//
// Copyright 2019 by KeyFuse Labs
// BSD License

package xrpc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func mockServer(str string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		data := []byte(str)
		resp := &Response{
			Result: (*json.RawMessage)(&data),
		}
		enc, _ := json.Marshal(resp)
		w.Write(enc)
	}))
}

func mockServerWithError(str string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := &Response{
			Error: &Error{Code: 88, Message: "mock.error"},
		}
		enc, _ := json.Marshal(resp)
		w.Write(enc)
	}))
}

func mockServerHTTPError(str string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, "oops", http.StatusMethodNotAllowed)
	}))
}

func mockBlockHeaders() []string {
	path := "testdata/blockheaders.json"
	datas, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var results []BlockHeaderResult
	err = json.Unmarshal(datas, &results)
	if err != nil {
		panic(err)
	}

	var resp []string
	for _, r := range results {
		json, _ := json.Marshal(r)
		resp = append(resp, string(json))
	}
	return resp
}

func mockBlocks() []string {
	path := "testdata/blocks.json"
	datas, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var results []BlockResult
	err = json.Unmarshal(datas, &results)
	if err != nil {
		panic(err)
	}

	var resp []string
	for _, r := range results {
		json, _ := json.Marshal(r)
		resp = append(resp, string(json))
	}
	return resp
}
