package cmd

import (
	"github.com/barnbridge/barnbridge-backend/api"
	"github.com/barnbridge/barnbridge-backend/core"
	"github.com/barnbridge/barnbridge-backend/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

var runCmd = &cobra.Command{
	Use:   "run",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperToDBFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		buildDBConnectionString()
		
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT)
		signal.Notify(stopChan, syscall.SIGTERM)
		
		d, err := db.New(db.Config{
			ConnectionString: viper.GetString("db.connection-string"),
			Automigrate:      viper.GetBool("feature.automigrate.enabled"),
		})
		if err != nil {
			log.Fatal("could not connect to db")
			return
		}
		
		c := core.New(core.Config{
		
		})
		go c.Run()
		
		a := api.New(api.Config{
			Port:           viper.GetString("api.port"),
			DevCorsEnabled: viper.GetBool("api.dev-cors"),
			DevCorsHost:    viper.GetString("api.dev-cors-host"),
		}, d)
		go a.Run()
		
		select {
		case <-stopChan:
			log.Info("Got stop signal. Finishing work.")
			// close whatever there is to close
			err := d.Close()
			if err != nil {
				log.Error(err)
			}
			
			c.Close()
			
			log.Info("Work done. Goodbye!")
		}
	},
}

func init() {
	addDBFlags(runCmd)
	
	// api
	runCmd.Flags().String("api.port", "3001", "HTTP API port")
	viper.BindPFlag("api.port", runCmd.Flag("api.port"))
	
	runCmd.Flags().Bool("api.dev-cors", false, "Enable development cors for HTTP API")
	viper.BindPFlag("api.dev-cors", runCmd.Flag("api.dev-cors"))
	
	runCmd.Flags().String("api.dev-cors-host", "", "Allowed host for HTTP API dev cors")
	viper.BindPFlag("api.dev-cors-host", runCmd.Flag("api.dev-cors-host"))
	
	// feature flags
	runCmd.Flags().Bool("feature.automigrate.enabled", true, "Enable/disable the automatic migrations feature")
	viper.BindPFlag("feature.automigrate.enabled", runCmd.Flag("feature.automigrate.enabled"))
}