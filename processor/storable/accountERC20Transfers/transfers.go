package accountERC20Transfers

import (
	web3types "github.com/alethio/web3-go/types"
)

func (s *Storable) decodeLogs(logs []web3types.Log) error {
	for _, log := range logs {
		t, err := s.decodeTransfer(log)
		if err != nil {
			return err
		}

		if t != nil {
			s.processed.transfers = append(s.processed.transfers, *t)
		}
	}
	return nil
}
