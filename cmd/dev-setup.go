package cmd

import (
	"database/sql"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

var devSetupCmd = &cobra.Command{
	Use:   "dev-setup",
	Short: "Setup env for development",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperToDBFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		buildDBConnectionString()

		db, err := sql.Open("postgres", viper.GetString("db.connection-string"))
		if err != nil {
			log.Fatal(err)
		}

		var pools = []types.SYPool{
			{
				ProtocolId:          "compound/v2",
				ControllerAddress:   utils.NormalizeAddress("0xE76ECAcDe3F1B6B2f63Ef803f3eDf1bB7988839b"),
				ModelAddress:        utils.NormalizeAddress("0x6837Ed0f69ab1a8946cFB69f92D7963DCdB99493"),
				ProviderAddress:     utils.NormalizeAddress("0x5e454d980F49AbE778a8cA5d0156f998f2d510CC"),
				SmartYieldAddress:   utils.NormalizeAddress("0xD165c8CAE4D824E75588282821C57fB3b74c7f33"),
				OracleAddress:       utils.NormalizeAddress("0x7837aa1D4177708407006189574412B59866Fc30"),
				JuniorBondAddress:   utils.NormalizeAddress("0xAad4380ED94C7372cbEac0f5AdA627B57b3D5C38"),
				SeniorBondAddress:   utils.NormalizeAddress("0xfcc5D3AEF0a7476B48CF46B6d5476bF933074C03"),
				ReceiptTokenAddress: utils.NormalizeAddress("0x4a92e71227d294f041bd82dd8f78591b75140d63"),
				UnderlyingAddress:   utils.NormalizeAddress("0xb7a4f3e9097c08da09517b5ab877f7a917224ede"),
				UnderlyingSymbol:    "USDC",
				UnderlyingDecimals:  6,
			},
		}

		for _, p := range pools {
			_, err = db.Exec("insert into smart_yield_pools (protocol_id, controller_address, model_address, provider_address, sy_address, oracle_address, junior_bond_address, senior_bond_address, receipt_token_address, underlying_address, underlying_symbol, underlying_decimals) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);", p.ProtocolId, p.ControllerAddress, p.ModelAddress, p.ProviderAddress, p.SmartYieldAddress, p.OracleAddress, p.JuniorBondAddress, p.SeniorBondAddress, p.ReceiptTokenAddress, p.UnderlyingAddress, p.UnderlyingSymbol, p.UnderlyingDecimals)
			if err != nil {
				log.Fatal(err)
			}
		}

		_, err = db.Exec("insert into smart_yield_reward_pools (pool_address, pool_token_address, reward_token_address) values ($1,$2,$3)", "0x31c5A8F6864AEDD4146BBE435A07b3d4d7Ef3595", "0x32", "0x3121")

		log.Println("done")
	},
}

func init() {
	RootCmd.AddCommand(devSetupCmd)
	addDBFlags(devSetupCmd)
}
