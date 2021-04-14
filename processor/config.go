package processor

import (
	"github.com/kekDAO/kekBackend/processor/storable/accountERC20Transfers"
	"github.com/kekDAO/kekBackend/processor/storable/barn"
	"github.com/kekDAO/kekBackend/processor/storable/bond"
	"github.com/kekDAO/kekBackend/processor/storable/governance"
	"github.com/kekDAO/kekBackend/processor/storable/yieldFarming"
)

type Config struct {
	Bond                  bond.Config
	Barn                  barn.Config
	Governance            governance.Config
	YieldFarming          yieldFarming.Config
	AccountErc20Transfers accountERC20Transfers.Config
}
