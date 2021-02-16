package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/barnbridge/barnbridge-backend/dashboard"
	"github.com/barnbridge/barnbridge-backend/processor"
	"github.com/barnbridge/barnbridge-backend/processor/storable/barn"
	"github.com/barnbridge/barnbridge-backend/processor/storable/bond"
	"github.com/barnbridge/barnbridge-backend/processor/storable/governance"
	"github.com/barnbridge/barnbridge-backend/processor/storable/yieldFarming"
	"github.com/barnbridge/barnbridge-backend/types"

	"github.com/barnbridge/barnbridge-backend/api"

	"github.com/barnbridge/barnbridge-backend/scraper"

	"github.com/barnbridge/barnbridge-backend/taskmanager"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/barnbridge/barnbridge-backend/core"
	"github.com/barnbridge/barnbridge-backend/eth/bestblock"
)

var runCmd = &cobra.Command{
	Use:    "run",
	Short:  "Track new blocks and index them",
	PreRun: runCmdPreRun,
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
				SlackNotify: types.SlackNotif{
					Enabled: viper.GetBool("feature.slack.enabled"),
					Webhook: viper.GetString("feature.slack.webhook"),
				},
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
					Address: viper.GetString("storable.yield.address"),
				},
			},
		})
		c.Run()

		a := api.New(c.DB(), api.Config{
			Port:           viper.GetString("api.port"),
			DevCorsEnabled: viper.GetBool("api.dev-cors"),
			DevCorsHost:    viper.GetString("api.dev-cors-host"),
		})
		go a.Run()

		d := dashboard.New(c, dashboard.Config{
			Port:          viper.GetString("dashboard.port"),
			ConfigEnabled: viper.GetBool("dashboard.config-management.enabled"),
		})
		go d.Run()

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

func runCmdPreRun(cmd *cobra.Command, args []string) {
	bindViperToDBFlags(cmd)
	bindViperToRedisFlags(cmd)
	bindViperToAPIFlags(cmd)
	bindViperToFeatureFlags(cmd)
	bindViperToEthFlags(cmd)
	bindViperToStorableFlags(cmd)

	viper.BindPFlag("abi-path", cmd.Flag("abi-path"))
}

func init() {
	addDBFlags(runCmd)
	addRedisFlags(runCmd)
	addAPIFlags(runCmd)
	addFeatureFlags(runCmd)
	addEthFlags(runCmd)
	addStorableFlags(runCmd)

	// dashboard
	runCmd.Flags().String("dashboard.port", "3000", "Dashboard port")
	viper.BindPFlag("dashboard.port", runCmd.Flag("dashboard.port"))

	runCmd.Flags().Bool("dashboard.config-management.enabled", true, "Enable/disable the config management option from dashboard")
	viper.BindPFlag("dashboard.config-management.enabled", runCmd.Flag("dashboard.config-management.enabled"))

	runCmd.Flags().String("abi-path", "./abis", "Path of directory from which to read contract ABIs")
}
