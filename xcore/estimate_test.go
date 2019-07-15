// tokucore
//
// Copyright (c) 2019 TokuBlock
// BSD License

package xcore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEstimateNormalSize(t *testing.T) {
	size := EstimateNormalSize(2, 2)
	assert.Equal(t, int64(374), size)
}
