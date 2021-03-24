package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/barnbridge/barnbridge-backend/api"
	"github.com/barnbridge/barnbridge-backend/dashboard"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:    "run",
	Short:  "Track new blocks and index them",
	PreRun: runCmdPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		buildDBConnectionString()
		requireNotEmptyFlags([]string{
			"storable.bond.address",
			"storable.barn.address",
			"storable.governance.address",
			"storable.yieldFarming.address",
			"storable.smartYieldState.compound-comptroller",
			"storable.smartYieldRewards.pool-factory-address",
		})

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT)
		signal.Notify(stopChan, syscall.SIGTERM)

		c := initCore()
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
