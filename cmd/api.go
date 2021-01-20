package cmd

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/barnbridge/barnbridge-backend/api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Run the API server exposing data from the database",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperToDBFlags(cmd)
		bindViperToAPIFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		buildDBConnectionString()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT)
		signal.Notify(stopChan, syscall.SIGTERM)

		log.Info("connecting to postgres")
		db, err := sql.Open("postgres", viper.GetString("db.connection-string"))
		if err != nil {
			log.Fatal(err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}

		log.Info("connected to postgres successfuly")

		a := api.New(db, api.Config{
			Port:           viper.GetString("api.port"),
			DevCorsEnabled: viper.GetBool("api.dev-cors"),
			DevCorsHost:    viper.GetString("api.dev-cors-host"),
			EthClientURL:   viper.GetString("eth.client.http"),
		})
		go a.Run()

		select {
		case <-stopChan:
			log.Info("Got stop signal. Finishing work.")

			log.Info("Work done. Goodbye!")
		}
	},
}

func init() {
	addDBFlags(apiCmd)
	addAPIFlags(apiCmd)
}
