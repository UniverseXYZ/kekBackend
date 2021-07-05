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
		if utils.CleanUpHex(a) == utils.CleanUpHex(log.Address) {
			return true
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
