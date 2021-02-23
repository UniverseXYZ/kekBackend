package cmd

import (
	"database/sql"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		_, err = db.Exec("insert into smart_yield_pools (protocol_id, controller_address, model_address, provider_address, sy_address, oracle_address, junior_bond_address, senior_bond_address, receipt_token_address, underlying_address, underlying_symbol, underlying_decimals) values ('compound/v2', '0xE76ECAcDe3F1B6B2f63Ef803f3eDf1bB7988839b', '0x6837Ed0f69ab1a8946cFB69f92D7963DCdB99493', '0x5e454d980F49AbE778a8cA5d0156f998f2d510CC', '0xD165c8CAE4D824E75588282821C57fB3b74c7f33', '0x7837aa1D4177708407006189574412B59866Fc30', '0xAad4380ED94C7372cbEac0f5AdA627B57b3D5C38', '0xfcc5D3AEF0a7476B48CF46B6d5476bF933074C03', '0x4a92e71227d294f041bd82dd8f78591b75140d63', '0xb7a4f3e9097c08da09517b5ab877f7a917224ede', 'USDC', 6);")
		if err != nil {
			log.Fatal(err)
		}

		log.Println("done")
	},
}

func init() {
	RootCmd.AddCommand(devSetupCmd)
	addDBFlags(devSetupCmd)
}
