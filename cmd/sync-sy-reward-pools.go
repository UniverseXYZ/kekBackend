package cmd

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"

	"github.com/kekDAO/kekBackend/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kekDAO/kekBackend/types"
)

var syncSyRewardPoolsCmd = &cobra.Command{
	Use:   "sync-sy-reward-pools",
	Short: "Sync SmartYield reward pools in the database with the ones in the json file",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperToDBFlags(cmd)
		viper.BindPFlag("file", cmd.Flag("file"))
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
			Rewards []types.SYRewardPool `json:"rewardPools"`
		}

		err = json.Unmarshal(data, &pools)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("removing current pools from database")

		_, err = db.Exec(`delete from smart_yield_reward_pools;`)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("done removing pools")

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
	RootCmd.AddCommand(syncSyRewardPoolsCmd)
	addDBFlags(syncSyRewardPoolsCmd)

	syncSyRewardPoolsCmd.Flags().String("file", "./pools.kovan.json", "Path to list of pools in json format")
}
