package core

import (
	"database/sql"
	"sync"
	"time"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/alethio/web3-go/ethrpc/provider/httprpc"
	"github.com/alethio/web3-go/validator"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
	"github.com/sirupsen/logrus"

	"github.com/kekDAO/kekBackend/eth/bestblock"
	"github.com/kekDAO/kekBackend/integrity"
	"github.com/kekDAO/kekBackend/metrics"
	"github.com/kekDAO/kekBackend/processor"
	"github.com/kekDAO/kekBackend/scraper"
	"github.com/kekDAO/kekBackend/state"
	"github.com/kekDAO/kekBackend/taskmanager"
)

var log = logrus.WithField("module", "core")

type Core struct {
	config Config

	metrics          *metrics.Provider
	bbtracker        *bestblock.Tracker
	taskmanager      *taskmanager.Manager
	scraper          *scraper.Scraper
	db               *sql.DB
	integrityChecker *integrity.Checker

	abis     map[string]abi.ABI
	ethConn  *ethclient.Client
	ethBatch *ethrpc.ETH

	stopMu sync.Mutex
}

func New(config Config) *Core {
	bbtracker, err := bestblock.NewTracker(config.BestBlockTracker)
	if err != nil {
		log.Fatal("could not start best block tracker")
		return nil
	}

	go bbtracker.Run()
	err = <-bbtracker.Err()
	if err != nil {
		log.Fatal("could not start best block tracker")
	}

	go func() {
		// todo: can we handle these errors?
		for err := range bbtracker.Err() {
			log.Error(err)
		}
	}()

	m := metrics.New()
	m.RecordLatestBlock(bbtracker.BestBlock())

	var lag int64
	if config.Features.Lag.Enabled {
		lag = config.Features.Lag.Value
	}

	tm, err := taskmanager.New(bbtracker, lag, m, config.TaskManager)
	if err != nil {
		log.Fatal("could not start task manager")
	}

	s, err := scraper.New(config.Scraper)
	if err != nil {
		log.Fatal("could not start scraper")
	}

	log.Info("connecting to postgres")
	db, err := sql.Open("postgres", config.PostgresConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	if config.Features.Automigrate {
		log.Info("attempting automatic execution of migrations")
		err = goose.Up(db, "/")
		if err != nil && err != goose.ErrNoNextVersion {
			log.Fatal(err)
		}
		log.Info("database version is up to date")
	}

	log.Info("connected to postgres successfuly")

	c := Core{
		config:      config,
		metrics:     m,
		bbtracker:   bbtracker,
		taskmanager: tm,
		scraper:     s,
		db:          db,
		abis:        make(map[string]abi.ABI),
	}

	log.Info("loading ABIs from contracts by given path")

	err = c.loadABIs()
	if err != nil {
		log.Fatal(err)
	}
	log.Info("done getting ABIs")

	// Create an IPC based RPC connection to a remote node
	conn, err := ethclient.Dial(config.Scraper.NodeURL)
	if err != nil {
		log.Fatalf("failed to connect to the Ethereum client: %v", err)
	}

	c.ethConn = conn

	batchLoader, err := httprpc.NewBatchLoader(0, 4*time.Millisecond)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not init batch loader"))
	}

	provider, err := httprpc.NewWithLoader(config.Scraper.NodeURL, batchLoader)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not init httprpc provider"))
	}

	c.ethBatch, err = ethrpc.New(provider)
	if err != nil {
		log.Fatal(err)
	}

	c.integrityChecker = integrity.NewChecker(c.db, c.bbtracker, c.taskmanager, lag)

	err = state.Init(c.db)
	if err != nil {
		log.Fatal(err)
	}

	return &c
}

func (c *Core) Run() {
	blockChan := make(chan int64)

	go func() {
		for b := range c.bbtracker.Subscribe() {
			c.Metrics().RecordLatestBlock(b)
		}
	}()

	go func() {
		max, err := c.getHighestBlock()
		if err != nil {
			log.Fatal("could not get highest block from db:", err)
		}

		log.WithField("block", max).Info("got highest block from db")

		best := c.bbtracker.BestBlock()

		log.WithField("block", best).Info("got highest block from network")

		if c.config.Features.Backfill {
			var lag int64
			if c.config.Features.Lag.Enabled {
				lag = c.config.Features.Lag.Value
			}

			backfillTarget := best - lag

			if max+1 < backfillTarget {
				log.Infof("adding tasks for %d blocks to be backfilled", backfillTarget-max+1)
				for i := max; i <= backfillTarget; i++ {
					err := c.taskmanager.Todo(i)
					if err != nil {
						log.Fatal("could not add task:", err)
					}
				}
			}
		} else {
			log.Info("skipping backfilling since feature is disabled")
		}
	}()

	go c.taskmanager.FeedToChan(blockChan)
	go c.integrityChecker.Run()

	go func() {
		for b := range blockChan {
			c.stopMu.Lock()
			log := log.WithField("block", b)
			log.Info("processing block")

			start := time.Now()
			blk, err := c.scraper.Exec(b)
			if err != nil {
				log.Error(err)

				c.stopMu.Unlock()
				err1 := c.taskmanager.Todo(b)
				if err1 != nil {
					log.Fatal(err1)
				}
				time.Sleep(2 * time.Second)
				continue
			}

			c.metrics.RecordScrapingTime(time.Since(start))

			log.Debug("validating block")
			v := validator.New()
			v.LoadBlock(blk.Block)
			if c.config.Features.Uncles {
				v.LoadUncles(blk.Uncles)
			}
			v.LoadReceipts(blk.Receipts)

			_, err = v.Run()
			if err != nil {
				c.stopMu.Unlock()
				c.metrics.RecordInvalidBlock()
				log.Error("error validating block: ", err)
				err1 := c.taskmanager.Todo(b)
				if err1 != nil {
					log.Fatal(err1)
				}
				continue
			}
			log.Debug("block is valid")

			log.Debug("storing block into the database")

			indexingStart := time.Now()

			p, err := processor.New(c.config.Processor, blk, c.abis, c.ethConn, c.ethBatch)
			if err != nil {
				c.stopMu.Unlock()
				log.Error("error storing block: ", err)
				err1 := c.taskmanager.Todo(b)
				if err1 != nil {
					log.Fatal(err1)
				}
				continue
			}

			err = p.Store(c.db, c.metrics)
			if err != nil {
				c.stopMu.Unlock()
				log.Error("error storing block: ", err)
				err1 := c.taskmanager.Todo(b)
				if err1 != nil {
					log.Fatal(err1)
				}
				continue
			}
			c.metrics.RecordIndexingTime(time.Since(indexingStart))
			c.metrics.RecordProcessingTime(time.Since(start))
			log.WithField("duration", time.Since(start)).Info("done processing block")
			c.stopMu.Unlock()
		}
	}()
}

func (c *Core) Close() error {
	c.stopMu.Lock()
	defer c.stopMu.Unlock()

	c.bbtracker.Close()
	log.Info("closed best block tracker")

	err := c.db.Close()
	if err != nil {
		return err
	}
	log.Info("closed db connection")

	errChan := make(chan error)
	go func() {
		err = c.taskmanager.Close()
		if err != nil {
			errChan <- err
		}
		log.Info("closed task manager")
		errChan <- nil
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(5 * time.Second):
		log.Warn("could not close task manager, exiting uncleanly")
	}

	return nil
}
