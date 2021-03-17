package cmd

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	formatter "github.com/kwix/logrus-module-formatter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initLogging() {
	logging := viper.GetString("logging")

	if verbose {
		logging = "*=debug"
	}

	if vverbose {
		logging = "*=trace"
	}

	if logging == "" {
		logging = "*=info"
	}

	gin.SetMode(gin.DebugMode)

	modules := formatter.NewModulesMap(logging)
	if level, exists := modules["gin"]; exists {
		if level < logrus.DebugLevel {
			gin.SetMode(gin.ReleaseMode)
		}
	} else {
		level := modules["*"]
		if level < logrus.DebugLevel {
			gin.SetMode(gin.ReleaseMode)
		}
	}

	f, err := formatter.New(modules)
	if err != nil {
		panic(err)
	}

	logrus.SetFormatter(f)

	log.Debug("Debug mode")
}

func addDBFlags(cmd *cobra.Command) {
	cmd.Flags().String("db.connection-string", "", "Postgres connection string.")
	cmd.Flags().String("db.host", "localhost", "Database host")
	cmd.Flags().String("db.port", "5432", "Database port")
	cmd.Flags().String("db.sslmode", "disable", "Database sslmode")
	cmd.Flags().String("db.dbname", "name", "Database name")
	cmd.Flags().String("db.user", "", "Database user (also allowed via PG_USER env)")
}

func bindViperToDBFlags(cmd *cobra.Command) {
	viper.BindPFlag("db.connection-string", cmd.Flag("db.connection-string"))
	viper.BindPFlag("db.host", cmd.Flag("db.host"))
	viper.BindPFlag("db.port", cmd.Flag("db.port"))
	viper.BindPFlag("db.sslmode", cmd.Flag("db.sslmode"))
	viper.BindPFlag("db.dbname", cmd.Flag("db.dbname"))
	viper.BindPFlag("db.user", cmd.Flag("db.user"))
}

func addAPIFlags(cmd *cobra.Command) {
	cmd.Flags().String("api.port", "3001", "HTTP API port")
	cmd.Flags().Bool("api.dev-cors", false, "Enable development cors for HTTP API")
	cmd.Flags().String("api.dev-cors-host", "", "Allowed host for HTTP API dev cors")
}

func bindViperToAPIFlags(cmd *cobra.Command) {
	viper.BindPFlag("api.port", cmd.Flag("api.port"))
	viper.BindPFlag("api.dev-cors", cmd.Flag("api.dev-cors"))
	viper.BindPFlag("api.dev-cors-host", cmd.Flag("api.dev-cors-host"))
}

func addRedisFlags(cmd *cobra.Command) {
	cmd.Flags().String("redis.server", "localhost:6379", "Redis server URL")
	cmd.Flags().String("redis.list", "todo", "The name of the list to be used for task management")
}

func bindViperToRedisFlags(cmd *cobra.Command) {
	viper.BindPFlag("redis.server", cmd.Flag("redis.server"))
	viper.BindPFlag("redis.list", cmd.Flag("redis.list"))
}

func buildDBConnectionString() {
	if viper.GetString("db.connection-string") == "" {
		var user, pass string
		if !viper.IsSet("db.user") {
			user = viper.GetString("PG_USER")
		} else {
			user = viper.GetString("db.user")
		}

		if !viper.IsSet("db.password") {
			pass = viper.GetString("PG_PASSWORD")
		} else {
			pass = viper.GetString("db.password")
		}

		p := fmt.Sprintf("host=%s port=%s sslmode=%s dbname=%s user=%s password=%s", viper.GetString("db.host"), viper.GetString("db.port"), viper.GetString("db.sslmode"), viper.GetString("db.dbname"), user, pass)
		viper.Set("db.connection-string", p)
	}
}

func addFeatureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("feature.backfill.enabled", true, "Enable/disable the automatic backfilling of data")
	cmd.Flags().Bool("feature.lag.enabled", false, "Enable/disable the lag behind feature (used to avoid reorgs)")
	cmd.Flags().Int64("feature.lag.value", 10, "The amount of blocks to lag behind the tip of the chain")
	cmd.Flags().Bool("feature.automigrate.enabled", true, "Enable/disable the automatic migrations feature")
}

func bindViperToFeatureFlags(cmd *cobra.Command) {
	viper.BindPFlag("feature.backfill.enabled", cmd.Flag("feature.backfill.enabled"))
	viper.BindPFlag("feature.lag.enabled", cmd.Flag("feature.lag.enabled"))
	viper.BindPFlag("feature.lag.value", cmd.Flag("feature.lag.value"))
	viper.BindPFlag("feature.automigrate.enabled", cmd.Flag("feature.automigrate.enabled"))
}

func addEthFlags(cmd *cobra.Command) {
	cmd.Flags().String("eth.client.http", "", "HTTP endpoint of JSON-RPC enabled Ethereum node")
	cmd.Flags().String("eth.client.ws", "", "WS endpoint of JSON-RPC enabled Ethereum node (provide this only if you want to use websocket subscription for tracking best block)")
	cmd.Flags().Duration("eth.client.poll-interval", 15*time.Second, "Interval to be used for polling the Ethereum node for best block")

}

func bindViperToEthFlags(cmd *cobra.Command) {
	viper.BindPFlag("eth.client.http", cmd.Flag("eth.client.http"))
	viper.BindPFlag("eth.client.ws", cmd.Flag("eth.client.ws"))
	viper.BindPFlag("eth.client.poll-interval", cmd.Flag("eth.client.poll-interval"))
}

func addStorableFlags(cmd *cobra.Command) {
	cmd.Flags().String("storable.bond.address", "0x0391D2021f89DC339F60Fff84546EA23E337750f", "BuyerAddress of the bond token")
	cmd.Flags().String("storable.barn.address", "0x19cFBFd65021af353aB8A7126Caf51920163f0D2", "BuyerAddress of the barn contract")
	cmd.Flags().String("storable.governance.address", "0x8EAcaEdD6D3BaCBC8A09C0787c5567f86eE96d02", "BuyerAddress of the governance contract")
	cmd.Flags().String("storable.yieldFarming.address", "0x2e93403C675Ccb9C564edf2dC6001233d0650582", "BuyerAddress of the yield farming contract")
	cmd.Flags().String("storable.smartYieldState.compound-comptroller", "0x3d9819210a31b4961b30ef54be2aed79b9c9cd3b", "Address of compound comptroller")
	cmd.Flags().String("storable.smartYieldState.compound-oracle-override", "", "Address to use instead of comptroller.oracle()")
	cmd.Flags().Int64("storable.smartYieldState.blocks-per-minute", 4, "How many blocks per minute on the blockchain we're scraping")
	cmd.Flags().Int64("storable.smartYieldPrice.startAt", 0, "How many blocks per minute on the blockchain we're scraping")
	cmd.Flags().String("storable.smartYieldRewards.pool-factory-address", "", "Address of rewards pool factory")
}

func bindViperToStorableFlags(cmd *cobra.Command) {
	viper.BindPFlag("storable.bond.address", cmd.Flag("storable.bond.address"))
	viper.BindPFlag("storable.barn.address", cmd.Flag("storable.barn.address"))
	viper.BindPFlag("storable.governance.address", cmd.Flag("storable.governance.address"))
	viper.BindPFlag("storable.yieldFarming.address", cmd.Flag("storable.yieldFarming.address"))
	viper.BindPFlag("storable.smartYieldState.compound-comptroller", cmd.Flag("storable.smartYieldState.compound-comptroller"))
	viper.BindPFlag("storable.smartYieldState.blocks-per-minute", cmd.Flag("storable.smartYieldState.blocks-per-minute"))
	viper.BindPFlag("storable.smartYieldPrice.startAt", cmd.Flag("storable.smartYieldPrice.startAt"))
	viper.BindPFlag("storable.smartYieldRewards.pool-factory-address", cmd.Flag("storable.smartYieldRewards.pool-factory-address"))
}
