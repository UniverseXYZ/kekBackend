package core

import (
	"github.com/kekDAO/kekBackend/eth/bestblock"
	"github.com/kekDAO/kekBackend/processor"
	"github.com/kekDAO/kekBackend/scraper"
	"github.com/kekDAO/kekBackend/taskmanager"
)

type Features struct {
	Backfill    bool
	Lag         FeatureLag
	Automigrate bool
	Uncles      bool
}

type FeatureLag struct {
	Enabled bool
	Value   int64
}

type Config struct {
	BestBlockTracker         bestblock.Config
	TaskManager              taskmanager.Config
	Scraper                  scraper.Config
	PostgresConnectionString string
	Features                 Features
	AbiPath                  string
	Processor                processor.Config
}
