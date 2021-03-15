package notifications

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type JobExecuter interface {
	ExecuteWithTx(ctx context.Context, tx *sql.Tx) (*Job, error)
}

type Job struct {
	Id              int64
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

func ExecuteJobsWithTx(ctx context.Context, tx *sql.Tx, jobs ...*Job) error {
	var nextJobs []*Job
	for _, j := range jobs {
		var je JobExecuter
		switch j.JobType {
		case ProposalCreated:
			var jd ProposalCreatedJobData
			err := json.Unmarshal(j.JobData, &jd)
			if err != nil {
				return errors.Wrap(err, "unmarshal proposal created job data")
			}
			je = &jd
		case ProposalActivating:
			var jd ProposalActivatingJobData
			err := json.Unmarshal(j.JobData, &jd)
			if err != nil {
				return errors.Wrap(err, "unmarshal proposal activated job data")
			}
			je = &jd
		default:
			return errors.Errorf("unknown job type %s", j.JobType)
		}

		n, err := je.ExecuteWithTx(ctx, tx)
		if err != nil {
			return errors.Wrap(err, "execute job")
		}
		if n != nil {
			nextJobs = append(nextJobs, n)
		}
	}

	if len(nextJobs) > 0 {
		err := ScheduleJobsWithTx(ctx, tx, nextJobs...)
		if err != nil {
			return errors.Wrap(err, "scheduling next jobs")
		}
	}
	return nil
}

func ScheduleJobsWithTx(ctx context.Context, tx *sql.Tx, jobs ...*Job) error {
	stmt, err := tx.PrepareContext(ctx, pq.CopyIn("notification_jobs", "type", "execute_on", "metadata", "included_in_block"))
	if err != nil {
		return errors.Wrap(err, "prepare notification job statement")
	}
	for _, j := range jobs {
		_, err := stmt.ExecContext(ctx, j.JobType, j.ExecuteOn, j.JobData, j.IncludedInBlock)
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

func DeleteJobsWithTx(ctx context.Context, tx *sql.Tx, jobs ...*Job) error {
	var ids []int64
	for _, j := range jobs {
		ids = append(ids, j.Id)
	}

	del := `
		DELETE FROM
			"notification_jobs"
		WHERE
			id = ANY($1)
		;
	`
	_, err := tx.ExecContext(ctx, del, pq.Array(ids))
	if err != nil {
		return errors.Wrap(err, "delete notification jobs")
	}
	return nil
}
