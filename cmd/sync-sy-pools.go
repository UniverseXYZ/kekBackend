package cmd

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"

	"github.com/barnbridge/barnbridge-backend/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/barnbridge/barnbridge-backend/types"
)

var syncSyPoolsCmd = &cobra.Command{
	Use:   "sync-sy-pools",
	Short: "Sync SmartYield pools in the database with the ones in the json file",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperToDBFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		buildDBConnectionString()

		db, err := sql.Open("postgres", viper.GetString("db.connection-string"))
		if err != nil {
			log.Fatal(err)
		}

		data, err := ioutil.ReadFile(viper.GetString("file"))
		if err != nil {
			log.Fatal(err)
		}

		var pools struct {
			SmartYield []types.SYPool       `json:"smartYield"`
			Rewards    []types.SYRewardPool `json:"rewardPools"`
		}

		err = json.Unmarshal(data, &pools)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("removing current pools from database")

		_, err = db.Exec(`delete from smart_yield_pools;`)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(`delete from smart_yield_reward_pools;`)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("done removing pools")

		log.WithField("count", len(pools.SmartYield)).Info("adding SmartYield pools from file to database")

		for _, p := range pools.SmartYield {
			_, err = db.Exec("insert into smart_yield_pools (protocol_id, controller_address, model_address, provider_address, sy_address, oracle_address, junior_bond_address, senior_bond_address, receipt_token_address, underlying_address, underlying_symbol, underlying_decimals, start_at_block) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);",
				p.ProtocolId,
				utils.NormalizeAddress(p.ControllerAddress),
				utils.NormalizeAddress(p.ModelAddress),
				utils.NormalizeAddress(p.ProviderAddress),
				utils.NormalizeAddress(p.SmartYieldAddress),
				utils.NormalizeAddress(p.OracleAddress),
				utils.NormalizeAddress(p.JuniorBondAddress),
				utils.NormalizeAddress(p.SeniorBondAddress),
				utils.NormalizeAddress(p.ReceiptTokenAddress),
				utils.NormalizeAddress(p.UnderlyingAddress),
				p.UnderlyingSymbol,
				p.UnderlyingDecimals,
				p.StartAtBlock,
			)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.WithField("count", len(pools.Rewards)).Info("adding Reward pools from file to database")

		for _, p := range pools.Rewards {
			_, err = db.Exec("insert into smart_yield_reward_pools (pool_address, pool_token_address, reward_token_address, start_at_block) values ($1,$2,$3, $4)",
				utils.NormalizeAddress(p.PoolAddress),
				utils.NormalizeAddress(p.PoolTokenAddress),
				utils.NormalizeAddress(p.RewardTokenAddress),
				p.StartAtBlock,
			)
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Println("done")
	},
}

func init() {
	RootCmd.AddCommand(syncSyPoolsCmd)
	addDBFlags(syncSyPoolsCmd)

	syncSyPoolsCmd.Flags().String("file", "./pools.kovan.json", "Path to list of pools in json format")
	viper.BindPFlag("file", syncSyPoolsCmd.Flag("file"))
}
