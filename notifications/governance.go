package notifications

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

type ProposalCreatedJobData ProposalJobData
type ProposalActivatingJobData ProposalJobData

type ProposalJobData struct {
	Id                    int64
	Proposer              string
	Title                 string
	CreateTime            int64
	WarmUpDuration        int64
	ActiveDuration        int64
	QueueDuration         int64
	GraceDuration         int64
	IncludedInBlockNumber int64
}

func (jd *ProposalCreatedJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) (*Job, error) {
	log.Tracef("executing proposal created job for PID-%d", jd.Id)

	// send created notification
	notif := NewNotification(
		"system",
		ProposalCreated,
		jd.CreateTime,
		jd.CreateTime+jd.WarmUpDuration-300,
		fmt.Sprintf("Proposal PID-%d created by %s", jd.Id, jd.Proposer),
		nil,
		jd.IncludedInBlockNumber,
	)

	err := notif.ToDBWithTx(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "save create proposal notification to db")
	}

	// schedule job for next notification
	njd := ProposalActivatingJobData(*jd)
	next, err := NewProposalActivatedJob(&njd)
	if err != nil {
		return nil, errors.Wrap(err, "create create proposal next job")
	}

	return next, nil
}

func (jd *ProposalActivatingJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) (*Job, error) {
	log.Tracef("executing proposal activated job for PID-%d", jd.Id)
	// check if proposal is still in warm up phase

	// send activated notification
	notif := NewNotification(
		"system",
		ProposalActivating,
		jd.CreateTime,
		jd.CreateTime+jd.WarmUpDuration-300,
		fmt.Sprintf("Proposal PID-%d activating in 5 minutes", jd.Id),
		nil,
		jd.IncludedInBlockNumber,
	)

	err := notif.ToDBWithTx(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "save activated proposal notification to db")
	}

	// // schedule job for next notification
	// jd := ProposalActivatingJobData(*jd)
	// next, err := NewProposalActivatedJob(&jd)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "create create proposal next job")
	// }

	//return next, nil
	return nil, nil
}

func NewProposalCreatedJob(data *ProposalCreatedJobData) (*Job, error) {
	return NewJob(ProposalCreated, 0, data.IncludedInBlockNumber, data)
}

func NewProposalActivatedJob(data *ProposalActivatingJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration - 300
	return NewJob(ProposalActivating, x, data.IncludedInBlockNumber, data)
}

// 		_, err = stmt.Exec(p.Id.Int64(), p.Proposer.String(), p.Title, p.CreateTime.Int64(), p.WarmUpDuration.Int64(), p.ActiveDuration.Int64(), p.QueueDuration.Int64(), p.GracePeriodDuration.Int64(), g.Preprocessed.BlockTimestamp)
// func FromGovernanceProposal(id int64, proposer string, title string, createTime int64, warmUpDuration int64, activeDuration int64, queueDuration int64, graceDuration int64, blockNumber int64, blockTime int64) []Notification {
// 	// TODO starts at blockTime -1 or creation time?
// 	startTime := blockTime - 1
//
// 	createNotif := NewNotification(
// 		"system",
// 		"proposal-created",
// 		blockNumber,
// 		startTime,
// 		startTime+warmUpDuration-300,
// 		fmt.Sprintf("Proposal PID-%d created by %s", id, proposer),
// 		nil,
// 	)
// 	activatingNotif := NewNotification(
// 		"system",
// 		"proposal-activating",
// 		blockNumber,
// 		startTime+warmUpDuration-300,
// 		startTime+warmUpDuration,
// 		fmt.Sprintf(fmt.Sprintf("Voting period for PID-%d starting in 5 minutes"), id),
// 		nil,
// 	)
// 	activeNotif := NewNotification(
// 		"system",
// 		"proposal-active",
// 		blockNumber,
// 		startTime+warmUpDuration,
// 		startTime+warmUpDuration+activeDuration-300,
// 		fmt.Sprintf(fmt.Sprintf("Governace proposal PID-%d is now active"), id),
// 		nil,
// 	)
//
// 	return []Notification{
// 		createNotif,
// 		activatingNotif,
// 		activeNotif,
// 	}
// }
