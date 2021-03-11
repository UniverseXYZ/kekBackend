package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/barnbridge/barnbridge-backend/notifications"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var notifsCmd = &cobra.Command{
	Use:    "notifications",
	Short:  "generate and store notifications",
	PreRun: notifsCmdPreRun,
	Run: func(cmd *cobra.Command, args []string) {
		buildDBConnectionString()

		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT)
		signal.Notify(stopChan, syscall.SIGTERM)

		n, err := notifications.New(notifications.Config{
			PostgresConnectionString: viper.GetString("db.connection-string"),
		})
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer func() {
			signal.Stop(stopChan)
			cancel()
		}()

		n.Run(ctx)

		select {
		case <-stopChan:
			log.Info("Got stop signal. Finishing work.")
			cancel()
			time.Sleep(time.Second * 2) // give it a few seconds so everything cleans up
			log.Info("Work done. Goodbye!")
		}
	},
}

func notifsCmdPreRun(cmd *cobra.Command, args []string) {
	bindViperToDBFlags(cmd)
}

func init() {
	RootCmd.AddCommand(notifsCmd)

	addDBFlags(notifsCmd)
	// addRedisFlags(scrapeCmd)
	// addFeatureFlags(scrapeCmd)
	// addEthFlags(scrapeCmd)
	// addStorableFlags(scrapeCmd)
}
