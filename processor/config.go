package processor

import (
	"github.com/barnbridge/barnbridge-backend/processor/storable/barn"
	"github.com/barnbridge/barnbridge-backend/processor/storable/bond"
	"github.com/barnbridge/barnbridge-backend/processor/storable/governance"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYield"
	"github.com/barnbridge/barnbridge-backend/processor/storable/yieldFarming"
)

type Config struct {
	Bond         bond.Config
	Barn         barn.Config
	Governance   governance.Config
	YieldFarming yieldFarming.Config
	SmartYield   smartYield.Config
}
