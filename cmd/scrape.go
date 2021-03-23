package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/barnbridge/barnbridge-backend/processor"
	"github.com/barnbridge/barnbridge-backend/processor/storable/barn"
	"github.com/barnbridge/barnbridge-backend/processor/storable/bond"
	"github.com/barnbridge/barnbridge-backend/processor/storable/governance"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYield"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYieldPrices"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYieldRewards"
	"github.com/barnbridge/barnbridge-backend/processor/storable/smartYieldState"
	"github.com/barnbridge/barnbridge-backend/processor/storable/yieldFarming"

	"github.com/barnbridge/barnbridge-backend/scraper"

	"github.com/barnbridge/barnbridge-backend/taskmanager"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/barnbridge/barnbridge-backend/core"
	"github.com/barnbridge/barnbridge-backend/eth/bestblock"
)

var scrapeCmd = &cobra.Command{
	Use:    "scrape",
	Short:  "Track new blocks and index them",
	PreRun: scrapeCmdPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		buildDBConnectionString()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT)
		signal.Notify(stopChan, syscall.SIGTERM)

		c := core.New(core.Config{
			BestBlockTracker: bestblock.Config{
				NodeURL:      viper.GetString("eth.client.http"),
				NodeURLWS:    viper.GetString("eth.client.ws"),
				PollInterval: viper.GetDuration("eth.client.poll-interval"),
			},
			TaskManager: taskmanager.Config{
				RedisServer:     viper.GetString("redis.server"),
				RedisPassword:   viper.GetString("REDIS_PASSWORD"),
				TodoList:        viper.GetString("redis.list"),
				BackfillEnabled: viper.GetBool("feature.backfill.enabled"),
			},
			Scraper: scraper.Config{
				NodeURL:      viper.GetString("eth.client.http"),
				EnableUncles: false,
			},
			PostgresConnectionString: viper.GetString("db.connection-string"),
			Features: core.Features{
				Backfill: viper.GetBool("feature.backfill.enabled"),
				Lag: core.FeatureLag{
					Enabled: viper.GetBool("feature.lag.enabled"),
					Value:   viper.GetInt64("feature.lag.value"),
				},
				Automigrate: viper.GetBool("feature.automigrate.enabled"),
				Uncles:      viper.GetBool("feature.uncles.enabled"),
			},
			AbiPath: viper.GetString("abi-path"),
			Processor: processor.Config{
				Bond: bond.Config{
					BondAddress: viper.GetString("storable.bond.address"),
				},
				Barn: barn.Config{
					BarnAddress: viper.GetString("storable.barn.address"),
				},
				Governance: governance.Config{
					GovernanceAddress: viper.GetString("storable.governance.address"),
				},
				YieldFarming: yieldFarming.Config{
					Address: viper.GetString("storable.yieldFarming.address"),
				},
				SmartYield: smartYield.Config{},
				SmartYieldState: smartYieldState.Config{
					ComptrollerAddress: viper.GetString("storable.smartYieldState.compound-comptroller"),
					BlocksPerMinute:    viper.GetInt64("storable.smartYieldState.blocks-per-minute"),
				},
				SmartYieldPrice: smartYieldPrices.Config{
					ComptrollerAddress: viper.GetString("storable.smartYieldState.compound-comptroller"),
				},
				SmartYieldRewards: smartYieldRewards.Config{
					PoolFactoryAddress: viper.GetString("storable.smartYieldRewards.pool-factory-address"),
				},
			},
		})
		c.Run()

		select {
		case <-stopChan:
			log.Info("Got stop signal. Finishing work.")
			err := c.Close()
			if err != nil {
				log.Fatal(err)
			}
			log.Info("Work done. Goodbye!")
		}
	},
}

func scrapeCmdPreRun(cmd *cobra.Command, args []string) {
	bindViperToDBFlags(cmd)
	bindViperToRedisFlags(cmd)
	bindViperToFeatureFlags(cmd)
	bindViperToEthFlags(cmd)
	bindViperToStorableFlags(cmd)

	viper.BindPFlag("abi-path", cmd.Flag("abi-path"))
}

func init() {
	RootCmd.AddCommand(scrapeCmd)

	addDBFlags(scrapeCmd)
	addRedisFlags(scrapeCmd)
	addFeatureFlags(scrapeCmd)
	addEthFlags(scrapeCmd)
	addStorableFlags(scrapeCmd)

	scrapeCmd.Flags().String("abi-path", "./abis", "Path of directory from which to read contract ABIs")
}
