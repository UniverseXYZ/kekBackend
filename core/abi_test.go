package core

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := &Core{
		config: Config{AbiPath: "../abis"},
		abis:   make(map[string]abi.ABI),
	}
	err := c.loadABIs()

	assert.NoError(t, err)
	for _, v := range c.abis {
		assert.NotEmpty(t, v)
	}
}
