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

		_, err = db.Exec("insert into smart_yield_pools (protocol_id, controller_address, model_address, provider_address, sy_address, oracle_address, junior_bond_address, senior_bond_address, ctoken_address, underlying_address, underlying_symbol, underlying_decimals) values ('compound/v2', '0x017d2260487F09387AFC40A26687AD14d4050243', '0xf2276DCbda583cad56C6eaAc29D4Fe19A6d15633', '0x851D638Cd5bBC4e879aFcEe6b882CF4e75abaf92', '0x17f040Ac9B1947e985F4532044299bCF8f9f859c', '0x97331B6be7e8Add530B58395a5da831d5aFeC0E6', '0x7b54D06Bb26a5d754E47AdB4007791919BDF5589', '0x54aC12aE1459B4001D0d13B0e5D18bF96680C1E1', '0x4a92e71227d294f041bd82dd8f78591b75140d63', '0xb7a4f3e9097c08da09517b5ab877f7a917224ede', 'USDC', 6);")
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
