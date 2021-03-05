package smartYieldRewards

import (
	"database/sql"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/sirupsen/logrus"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/types"
)

//event Claim(address indexed user, uint256 amount);
//event Deposit(address indexed user, uint256 amount, uint256 balanceAfter);
//event Withdraw(address indexed user, uint256 amount, uint256 balanceAfter);

var log = logrus.WithField("module", "storable(smart yield rewards)")

type Storable struct {
	config      Config
	raw         *types.RawData
	syRewardABI abi.ABI

	processed struct {
		stakingActions []StakingAction
		claims         []ClaimEvent
		blockNumber    int64
		blockTimestamp int64
	}
}

func NewStorable(config Config, raw *types.RawData, syRewardABI abi.ABI) *Storable {
	return &Storable{
		config:      config,
		raw:         raw,
		syRewardABI: syRewardABI,
	}
}

func (s *Storable) ToDB(tx *sql.Tx) error {
	var rewardLogs []web3types.Log

	for _, data := range s.raw.Receipts {
		for _, log := range data.Logs {
			if state.PoolByAddress(log.Address) != nil {
				rewardLogs = append(rewardLogs, log)
			}
		}
	}

	if len(rewardLogs) == 0 {
		log.WithField("handler", "smart yield rewards").Debug("No events found")
		return nil
	}

	return nil
}
