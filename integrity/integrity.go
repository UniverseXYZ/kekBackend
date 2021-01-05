package integrity

import (
	"database/sql"
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/barnbridge/barnbridge-backend/eth/bestblock"
	"github.com/barnbridge/barnbridge-backend/taskmanager"
)

type Checker struct {
	db      *sql.DB
	tracker *bestblock.Tracker
	tm      *taskmanager.Manager
	logger  *logrus.Entry
}

func NewChecker(db *sql.DB, tracker *bestblock.Tracker, tm *taskmanager.Manager) *Checker {
	return &Checker{
		db:      db,
		tracker: tracker,
		tm:      tm,
		logger:  logrus.WithField("module", "integrity-checker"),
	}
}

func (c *Checker) Run() {
	t := time.NewTicker(1 * time.Minute)

	for {
		select {
		case <-t.C:
			err := c.lifecycle()
			if err != nil {
				c.logger.Error(err)
			}
		}
	}
}

func (c *Checker) lifecycle() error {
	c.logger.Trace("running")
	start := time.Now()
	defer func() {
		c.logger.WithField("duration", time.Since(start)).Trace("done")
	}()

	best := c.tracker.BestBlock()
	checkpoint, err := c.getLastCheckpoint()
	if err != nil {
		return err
	}
	if checkpoint == -1 {
		return nil
	}

	missing, err := c.checkMissingBlocks(checkpoint, best)
	if err != nil {
		return err
	}

	broken, err := c.checkBrokenHashChain(checkpoint, best)
	if err != nil {
		return err
	}

	all := append(missing, broken...)
	if len(all) == 0 {
		_, err = c.db.Exec("insert into integrity_checkpoints (number) values((select max(number) from blocks))")
		if err != nil {
			return errors.Wrap(err, "could not store new integrity checkpoint")
		}

		c.logger.Info("finished checking integrity; all good!")
		return nil
	}

	var uniqueBlocks = make(map[int64]bool)

	for _, block := range append(missing, broken...) {
		uniqueBlocks[block] = true
	}

	var blocks []int64
	for k := range uniqueBlocks {
		blocks = append(blocks, k)
	}

	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i] < blocks[j]
	})

	for _, block := range blocks {
		err = c.tm.Todo(block)
		if err != nil {
			return errors.Wrap(err, "could not queue block for rescrape")
		}
	}

	_, err = c.db.Exec("insert into integrity_checkpoints (number) values($1)", blocks[0]-1)
	if err != nil {
		return errors.Wrap(err, "could not store new integrity checkpoint")
	}

	c.logger.WithField("count", len(blocks)).Warn("found inconsistent blocks & queued for rescrape")

	return nil
}

func (c *Checker) getLastCheckpoint() (int64, error) {
	var b int64
	err := c.db.QueryRow(`select number from integrity_checkpoints order by created_at desc limit 1`).Scan(&b)
	if err == sql.ErrNoRows {
		err1 := c.db.QueryRow(`select min(number) from blocks`).Scan(&b)
		if err1 == sql.ErrNoRows {
			return -1, nil
		}
		if err1 != nil {
			return 0, errors.Wrap(err, "could not get min block number from db")
		}

		return b, nil
	}
	if err != nil {
		return 0, errors.Wrap(err, "could not get latest integrity checkpoint from db")
	}

	return b, nil
}

func (c *Checker) checkMissingBlocks(start, end int64) ([]int64, error) {
	rows, err := c.db.Query(`
		select x.number
		from generate_series($1::bigint, $2::bigint) as x(number)
				 left join (select number from blocks where number between $1 and $2) b on x.number = b.number
		where b.number is null
		order by number;
	`, start, end)
	if err != nil {
		return nil, errors.Wrap(err, "could not query database for missing blocks")
	}

	var blocks []int64
	for rows.Next() {
		var b int64

		err = rows.Scan(&b)
		if err != nil {
			return nil, errors.Wrap(err, "could not scan missing block from db")
		}

		blocks = append(blocks, b)
	}

	return blocks, nil
}

func (c *Checker) checkBrokenHashChain(start, end int64) ([]int64, error) {
	rows, err := c.db.Query(`
		with a as (
			select number
			from blocks as t1
			where t1.number between $1 and $2
			  and (select block_hash from blocks as t2 where t2.number = t1.number - 1) != t1.parent_block_hash
		)
		select number
		from a
		union all
		select number - 1
		from a
		order by number;
	`, start-100, end)
	if err != nil {
		return nil, errors.Wrap(err, "could not query database for broken hash chain")
	}

	var blocks []int64
	for rows.Next() {
		var b int64

		err = rows.Scan(&b)
		if err != nil {
			return nil, errors.Wrap(err, "could not scan inconsistent block from db")
		}

		blocks = append(blocks, b)
	}

	return blocks, nil
}
