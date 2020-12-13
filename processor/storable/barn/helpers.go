package barn

import (
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/utils"
)

func (b *BarnStorable) getBaseLog(log web3types.Log) (*BaseLog, error) {
	txIndex, err := strconv.ParseInt(log.TransactionIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert transactionIndex from barn contract to int64")
	}

	logIndex, err := strconv.ParseInt(log.LogIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert logIndex from  barn contract to int64")
	}

	return &BaseLog{
		LoggedBy:         utils.CleanUpHex(log.Address),
		TransactionHash:  utils.CleanUpHex(log.TransactionHash),
		TransactionIndex: txIndex,
		LogIndex:         logIndex,
	}, nil
}
