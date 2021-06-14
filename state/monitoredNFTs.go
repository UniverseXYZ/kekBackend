package state

import (
	web3types "github.com/alethio/web3-go/types"

	"github.com/kekDAO/kekBackend/utils"
)

func MonitoredNFTs() []string {
	return instance.monitoredNTFs
}

func IsMonitoredNFT(log web3types.Log) bool {
	for _, a := range instance.monitoredNTFs {
		if len(log.Topics) >= 3 {
			if utils.NormalizeAddress(a) == utils.Topic2Address(log.Topics[1]) ||
				utils.NormalizeAddress(a) == utils.Topic2Address(log.Topics[2]) {
				return true
			}
		}
	}
	return false
}

func AddMonitoredNFTToState(nft string) {
	instance.monitoredNTFs = append(instance.monitoredNTFs, nft)
}

func AddMonitoredNFTToDB(nft string) error {
	_, err := instance.db.Exec(`insert into monitored_nfts (address) values ($1)`,
		nft)

	if err != nil {
		return err
	}

	return nil
}
