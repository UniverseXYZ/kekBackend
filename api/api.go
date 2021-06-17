package api

import (
	"database/sql"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"

	"github.com/kekDAO/kekBackend/contracts"
	"github.com/kekDAO/kekBackend/state"
)

var log = logrus.WithField("module", "api")

type Config struct {
	Port           string
	DevCorsEnabled bool
	DevCorsHost    string
	XYZ            string
	RPCUrl         string
}

type API struct {
	config Config
	engine *gin.Engine

	db  *sql.DB
	eth *ethclient.Client
	xyz *contracts.ERC20

	CirculatingSupply decimal.Decimal
}

func New(db *sql.DB, config Config) *API {
	err := state.Init(db)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := ethclient.Dial(config.RPCUrl)
	if err != nil {
		log.Fatalf("failed to connect to the Ethereum client: %v", err)
	}

	xyz, err := contracts.NewERC20(common.HexToAddress(config.XYZ), conn)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not initialize BOND erc20 contract"))
	}

	return &API{
		config: config,
		db:     db,
		eth:    conn,
		xyz:    xyz,
	}
}

func (a *API) Run() {
	a.engine = gin.Default()

	t := time.NewTicker(1 * time.Minute)

	go a.updateCirculatingSupply()

	go func() {
		log.Info("setting up ticker to refresh state")

		for {
			select {
			case <-t.C:
				err := state.Refresh()
				if err != nil {
					log.Error(err)
				}
			}
		}
	}()

	if a.config.DevCorsEnabled {
		a.engine.Use(cors.New(cors.Config{
			AllowOrigins:     []string{a.config.DevCorsHost},
			AllowMethods:     []string{"PUT", "PATCH", "GET", "POST"},
			AllowHeaders:     []string{"Origin"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
		}))
	}

	a.setRoutes()

	err := a.engine.Run(":" + a.config.Port)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *API) Close() {
}
