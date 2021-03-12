package notifications

import (
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	ProposalCreated   = "proposal-created"
	ProposalActivated = "proposal-activated"
)

type JobExecuter interface {
	ExecuteWithTx(tx *sql.Tx) (*Job, error)
}

type Job struct {
	JobType         string
	ExecuteOn       int64
	JobData         json.RawMessage
	IncludedInBlock int64
}

func NewJob(typ string, executeOn int64, block int64, data interface{}) (*Job, error) {
	d, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "new job marshal")
	}

	return &Job{
		JobType:         typ,
		ExecuteOn:       executeOn,
		JobData:         d,
		IncludedInBlock: block,
	}, nil
}

func ExecuteJobsWithTx(tx *sql.Tx, jobs ...*Job) error {
	var nextJobs []*Job
	for _, j := range jobs {
		var je JobExecuter
		switch j.JobType {
		case ProposalCreated:
			var jd ProposalCreatedJobData
			err := json.Unmarshal(j.JobData, &jd)
			if err != nil {
				return errors.Wrap(err, "unmarshal job data")
			}
			je = &jd
		default:
			return errors.Errorf("unknown job type %s", j.JobType)
		}
		n, err := je.ExecuteWithTx(tx)
		if err != nil {
			return errors.Wrap(err, "execute job")
		}
		nextJobs = append(nextJobs, n)
	}

	if len(nextJobs) > 0 {
		err := ScheduleJobsWithTx(tx, nextJobs...)
		if err != nil {
			return errors.Wrap(err, "scheduling next jobs")
		}
	}
	return nil
}

func ScheduleJobsWithTx(tx *sql.Tx, jobs ...*Job) error {
	stmt, err := tx.Prepare(pq.CopyIn("notification_jobs", "type", "execute_on", "metadata", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "prepare notification job statement")
	}
	for _, j := range jobs {
		_, err := stmt.Exec(j.JobType, j.ExecuteOn, j.JobData, j.IncludedInBlock)
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
