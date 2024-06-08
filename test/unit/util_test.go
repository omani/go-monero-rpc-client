package test

import (
	"testing"

	"github.com/monero-ecosystem/go-monero-rpc-client/util"
	"github.com/stretchr/testify/assert"
)

func TestXMRToDecimalTest(t *testing.T) {
	assert.Equal(t, "0.034000200000", util.XMRToDecimal(34000200000))
	assert.Equal(t, "15.000000000000", util.XMRToDecimal(15e12))
}

func TestXMRToFloat64Test(t *testing.T) {
	assert.Equal(t, float64(0.02), util.XMRToFloat64(20000000000))
	assert.Equal(t, float64(3.14), util.XMRToFloat64(314e10))
}
