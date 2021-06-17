package core

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

func (c *Core) loadABIs() error {
	files, err := ioutil.ReadDir(c.config.AbiPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".json") {
			byteValue, err := ioutil.ReadFile(c.config.AbiPath + "/" + file.Name())
			if err != nil {
				return err
			}

			a, err := abi.JSON(bytes.NewReader(byteValue))
			if err != nil {
				return errors.Wrap(err, "could not decode abi for: "+file.Name())
			}

			key := strings.ToLower(strings.TrimSuffix(file.Name(), ".json"))
			c.abis[key] = a
		}
	}

	return nil
}
