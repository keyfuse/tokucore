// tokucore
//
// Copyright (c) 2018 TokuBlock
// BSD License

package xvm

import (
	"encoding/json"
	"io/ioutil"
)

func readTests(testfile string) ([][]interface{}, error) {
	file, err := ioutil.ReadFile(testfile)
	if err != nil {
		return nil, err
	}

	var tests [][]interface{}
	if err := json.Unmarshal(file, &tests); err != nil {
		return nil, err
	}
	return tests, nil
}
