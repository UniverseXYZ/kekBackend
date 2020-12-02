package core

import (
	"io/ioutil"
	"strings"
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

			key := strings.TrimSuffix(file.Name(), ".json")
			c.abis[key] = byteValue
		}
	}

	return nil
}
