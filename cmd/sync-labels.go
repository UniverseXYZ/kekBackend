package cmd

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

var syncLabelsCmd = &cobra.Command{
	Use:   "sync-labels",
	Short: "Sync labels in the database with the ones in the json file",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperToDBFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		buildDBConnectionString()

		db, err := sql.Open("postgres", viper.GetString("db.connection-string"))
		if err != nil {
			log.Fatal(err)
		}

		data, err := ioutil.ReadFile(viper.GetString("labels-file"))
		if err != nil {
			log.Fatal(err)
		}

		var jsonLabels struct {
			Labels []types.LabelStruct `json:"labels"`
		}

		err = json.Unmarshal(data, &jsonLabels)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("removing current labels from database")

		_, err = db.Exec(`delete from labels;`)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("done removing labels")

		log.WithField("count", len(jsonLabels.Labels)).Info("adding labels from file to database")

		for _, l := range jsonLabels.Labels {
			_, err = db.Exec(`insert into labels (address,label) values ($1,$2)`,
				utils.NormalizeAddress(l.Address),
				l.Label)
			if err != nil {
				log.Fatal(err)
			}
		}
		log.Println("done")
	},
}

func init() {
	RootCmd.AddCommand(syncLabelsCmd)
	addDBFlags(syncLabelsCmd)

	syncLabelsCmd.Flags().String("labels-file", "./labels.kovan.json", "Path to list of labels in json format")
	viper.BindPFlag("labels-file", syncLabelsCmd.Flag("labels-file"))
}
