package processor

import (
	"github.com/barnbridge/barnbridge-backend/processor/storable/accountERC20Transfers"
	"github.com/barnbridge/barnbridge-backend/processor/storable/barn"
	"github.com/barnbridge/barnbridge-backend/processor/storable/bond"
	"github.com/barnbridge/barnbridge-backend/processor/storable/governance"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYield"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYieldPrices"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYieldRewards"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYieldState"
	"github.com/barnbridge/barnbridge-backend/processor/storable/yieldFarming"
)

type Config struct {
	Bond                  bond.Config
	Barn                  barn.Config
	Governance            governance.Config
	YieldFarming          yieldFarming.Config
	SmartYield            smartYield.Config
	SmartYieldState       smartYieldState.Config
	SmartYieldPrice       smartYieldPrices.Config
	SmartYieldRewards     smartYieldRewards.Config
	AccountErc20Transfers accountERC20Transfers.Config
}
