package processor

import (
	"github.com/barnbridge/barnbridge-backend/processor/storable/barn"
	"github.com/barnbridge/barnbridge-backend/processor/storable/bond"
)

type Config struct {
	Bond bond.Config
	Barn barn.Config
}
