package notifications

import (
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	ProposalCreated = "proposal-created"
)

type Job struct {
	JobType         string
	ExecuteOn       int64
	Metadata        json.RawMessage
	IncludedInBlock int64
}

func NewJob(typ string, executeOn int64, data interface{}) (*Job, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "new job marshal")
	}

	return &Job{
		JobType:         typ,
		ExecuteOn:       executeOn,
		Metadata:        d,
		IncludedInBlock: 0,
	}, nil
}

func ExecuteJobsWithTx(tx *sql.Tx, jobs ...*Job) error {
	_ = tx
	for _, j := range jobs {
		switch j.JobType {
		case ProposalCreated:
			log.Trace("executing proposal created job")
		}
	}
	return nil
}

func ScheduleJobsWithTx(tx *sql.Tx, jobs ...*Job) error {
	stmt, err := tx.Prepare(pq.CopyIn("notification_job", "type", "execute_on", "metadata", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "prepare notification job statement")
	}
	for _, j := range jobs {
		_, err := stmt.Exec(j.JobType, j.ExecuteOn, j.Metadata, j.IncludedInBlock)
		if err != nil {
			return errors.Wrap(err, "could not exec statement")
		}
	}

	err = stmt.Close()
	if err != nil {
		return errors.Wrap(err, "could not close exec statement")
	}

	return nil
}
