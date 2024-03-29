package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/kekDAO/kekBackend/migrations"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the database to an empty state by truncating all the tables",
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViperToDBFlags(cmd)
		bindViperToRedisFlags(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.GetBool("force") {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("This will reset the database. Are you sure? [y/N]: ")
			text, _ := reader.ReadString('\n')
			if strings.TrimSuffix(strings.ToLower(text), "\n") != "y" {
				fmt.Println("Nobody was harmed.")
				return
			}
		}

		fmt.Print("Deleting todo queue from redis ... ")

		r := redis.NewClient(&redis.Options{
			Addr:        viper.GetString("redis.server"),
			Password:    viper.GetString("REDIS_PASSWORD"),
			DB:          0,
			ReadTimeout: time.Second * 1,
		})

		err := r.Ping().Err()
		if err != nil {
			log.Fatal(err)
			return
		}

		err = r.Del(viper.GetString("redis.list")).Err()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("[done]")

		fmt.Print("Truncating database ... ")

		buildDBConnectionString()

		db, err := sql.Open("postgres", viper.GetString("db.connection-string"))
		if err != nil {
			log.Fatal(err)
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}

		_, err = tx.Exec(`
		drop schema public cascade;
		create schema public;
		`)
		if err != nil {
			log.Fatal(err)
		}

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("[done]")

		fmt.Println("Database was reset. Have fun!")
	},
}

func init() {
	addDBFlags(resetCmd)
	addRedisFlags(resetCmd)

	resetCmd.Flags().Bool("force", false, "Run this command without entering interactive mode")
	viper.BindPFlag("force", resetCmd.Flag("force"))
}
