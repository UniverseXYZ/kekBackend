package taskmanager

import (
	"time"

	"github.com/go-redis/redis"
)

// Todo inserts a block into the redis sorted set used for queue management using a ZADD command
func (m *Manager) Todo(block int64) error {
	log.WithField("block", block).Trace("adding block to todo")
	return m.redis.ZAdd(m.config.TodoList, redis.Z{
		Score:  float64(block),
		Member: block,
	}).Err()
}

func (m *Manager) Reset() error {
	err := m.redis.Del(m.config.TodoList).Err()
	if err != nil {
		return err
	}

	// set the lastBlockAdded to -1 in order to backfill the whole chain after a reset if the backfill feature is enabled
	m.lastBlockAdded = -1

	return nil
}

func (m *Manager) BatchTodo(list string, from int64, to int64) error {
	start := time.Now()

	var members []redis.Z
	for i := from; i <= to; i++ {
		members = append(members, redis.Z{
			Score:  float64(i),
			Member: i,
		})
	}

	const batchSize = 500000

	batches := int(to-from+1)/batchSize + 1

	for i := 0; i < batches; i++ {
		end := batchSize * (i + 1)
		if end > len(members) {
			end = len(members)
		}
		log.Tracef("queueing batch [%d, %d]", members[batchSize*i].Member, members[end-1].Member)
		err := m.redis.ZAdd(list, members[batchSize*i:end]...).Err()
		if err != nil && err != redis.Nil {
			return err
		}
	}

	log.WithField("duration", time.Since(start)).Trace("queued all blocks")

	return nil
}
