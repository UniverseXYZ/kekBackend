package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := &Core{
		config: Config{AbiPath: "../abis"},
		abis:   make(map[string][]byte),
	}
	err := c.loadABIs()

	assert.NoError(t, err)
	assert.Len(t, c.abis, 3)
	for _, v := range c.abis {
		assert.NotEmpty(t, v)
	}
}
