package cmd

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/barnbridge/barnbridge-backend/dashboard"
	"github.com/barnbridge/barnbridge-backend/processor"
	"github.com/barnbridge/barnbridge-backend/processor/storable/barn"
	"github.com/barnbridge/barnbridge-backend/processor/storable/bond"

	"github.com/barnbridge/barnbridge-backend/api"

	"github.com/barnbridge/barnbridge-backend/scraper"

	"github.com/barnbridge/barnbridge-backend/taskmanager"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/barnbridge/barnbridge-backend/core"
	"github.com/barnbridge/barnbridge-backend/eth/bestblock"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Track new blocks and index them",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperToDBFlags(cmd)
		bindViperToRedisFlags(cmd)
	},
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
				EnableUncles: viper.GetBool("feature.uncles.enabled"),
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
			},
		})
		c.Run()

		a := api.New(c, api.Config{
			Port:           viper.GetString("api.port"),
			DevCorsEnabled: viper.GetBool("api.dev-cors"),
			DevCorsHost:    viper.GetString("api.dev-cors-host"),
			EthClientURL:   viper.GetString("eth.client.http"),
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

func init() {
	addDBFlags(runCmd)
	addRedisFlags(runCmd)

	// feature flags
	runCmd.Flags().Bool("feature.backfill.enabled", true, "Enable/disable the automatic backfilling of data")
	viper.BindPFlag("feature.backfill.enabled", runCmd.Flag("feature.backfill.enabled"))

	runCmd.Flags().Bool("feature.lag.enabled", false, "Enable/disable the lag behind feature (used to avoid reorgs)")
	viper.BindPFlag("feature.lag.enabled", runCmd.Flag("feature.lag.enabled"))

	runCmd.Flags().Int64("feature.lag.value", 10, "The amount of blocks to lag behind the tip of the chain")
	viper.BindPFlag("feature.lag.value", runCmd.Flag("feature.lag.value"))

	runCmd.Flags().Bool("feature.automigrate.enabled", true, "Enable/disable the automatic migrations feature")
	viper.BindPFlag("feature.automigrate.enabled", runCmd.Flag("feature.automigrate.enabled"))

	runCmd.Flags().Bool("feature.uncles.enabled", true, "Enable/disable uncles scraping")
	viper.BindPFlag("feature.uncles.enabled", runCmd.Flag("feature.uncles.enabled"))

	// eth
	runCmd.Flags().String("eth.client.http", "", "HTTP endpoint of JSON-RPC enabled Ethereum node")
	viper.BindPFlag("eth.client.http", runCmd.Flag("eth.client.http"))

	runCmd.Flags().String("eth.client.ws", "", "WS endpoint of JSON-RPC enabled Ethereum node (provide this only if you want to use websocket subscription for tracking best block)")
	viper.BindPFlag("eth.client.ws", runCmd.Flag("eth.client.ws"))

	runCmd.Flags().Duration("eth.client.poll-interval", 15*time.Second, "Interval to be used for polling the Ethereum node for best block")
	viper.BindPFlag("eth.client.poll-interval", runCmd.Flag("eth.client.poll-interval"))

	// api
	runCmd.Flags().String("api.port", "3001", "HTTP API port")
	viper.BindPFlag("api.port", runCmd.Flag("api.port"))

	runCmd.Flags().Bool("api.dev-cors", false, "Enable development cors for HTTP API")
	viper.BindPFlag("api.dev-cors", runCmd.Flag("api.dev-cors"))

	runCmd.Flags().String("api.dev-cors-host", "", "Allowed host for HTTP API dev cors")
	viper.BindPFlag("api.dev-cors-host", runCmd.Flag("api.dev-cors-host"))

	// dashboard
	runCmd.Flags().String("dashboard.port", "3000", "Dashboard port")
	viper.BindPFlag("dashboard.port", runCmd.Flag("dashboard.port"))

	runCmd.Flags().Bool("dashboard.config-management.enabled", true, "Enable/disable the config management option from dashboard")
	viper.BindPFlag("dashboard.config-management.enabled", runCmd.Flag("dashboard.config-management.enabled"))

	//abi
	runCmd.Flags().String("abi-path", "./abis", "Path of directory from which to read contract ABIs")
	viper.BindPFlag("abi-path", runCmd.Flag("abi-path"))

	//bond
	runCmd.Flags().String("storable.bond.address", "0x0391D2021f89DC339F60Fff84546EA23E337750f", "Address of the bond token")
	viper.BindPFlag("storable.bond.address", runCmd.Flag("storable.bond.address"))

	//barn
	runCmd.Flags().String("storable.barn.address", "0x9b24CA2bd6C44490F0529bCa48E837E2ad26BaFa", "Address of the barn token")
	viper.BindPFlag("storable.barn.address", runCmd.Flag("storable.barn.address"))
}
