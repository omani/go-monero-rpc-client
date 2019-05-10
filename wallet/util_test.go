package wallet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXMRToDecimal(t *testing.T) {
	assert.Equal(t, "0.034000200000", XMRToDecimal(34000200000))
	assert.Equal(t, "15.000000000000", XMRToDecimal(15e12))
}

func TestXMRToFloat64(t *testing.T) {
	assert.Equal(t, float64(0.02), XMRToFloat64(20000000000))
	assert.Equal(t, float64(3.14), XMRToFloat64(314e10))
}
