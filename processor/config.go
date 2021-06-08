package processor

import (
	"github.com/kekDAO/kekBackend/processor/storable/accountERC20Transfers"
	"github.com/kekDAO/kekBackend/processor/storable/auction"
	"github.com/kekDAO/kekBackend/processor/storable/governance"
	"github.com/kekDAO/kekBackend/processor/storable/kek"
	"github.com/kekDAO/kekBackend/processor/storable/supernova"
	"github.com/kekDAO/kekBackend/processor/storable/universe"
	"github.com/kekDAO/kekBackend/processor/storable/yieldFarming"
)

type Config struct {
	Kek                   kek.Config
	Supernova             supernova.Config
	Governance            governance.Config
	YieldFarming          yieldFarming.Config
	Auction               auction.Config
	Universe              universe.Config
	AccountErc20Transfers accountERC20Transfers.Config
}
