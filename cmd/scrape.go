package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var scrapeCmd = &cobra.Command{
	Use:    "scrape",
	Short:  "Track new blocks and index them",
	PreRun: scrapeCmdPreRun,
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
