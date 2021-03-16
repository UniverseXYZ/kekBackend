package notifications

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

const (
	ProposalCreated           = "proposal-created"
	ProposalActivating        = "proposal-activated"
	ProposalVotingOpen        = "proposal-voting-open"
	ProposalVotingEnding      = "proposal-voting-ending"
	ProposalOutcome           = "proposal-outcome"
	ProposalAccepted          = "proposal-accepted"
	ProposalFailed            = "proposal-failed"
	ProposalGracePeriod       = "proposal-grace"
	ProposalFinalState        = "proposal-final-state"
	AbrogationProposalCreated = "abrogation-proposal-created"
)

const (
	ProposalStateWarmUp      = "WARMUP"
	ProposalStateActive      = "ACTIVE"
	ProposalStateAccepted    = "ACCEPTED"
	ProposalStateFailed      = "FAILED"
	ProposalStateQueued      = "QUEUED"
	ProposalStateGracePeriod = "GRACE"
	ProposalStateExecuted    = "EXECUTED"
	ProposalStateExpired     = "EXPIRED"
	ProposalStateCanceled    = "CANCELED"
)

// new proposal
type ProposalCreatedJobData ProposalJobData
type ProposalActivatingJobData ProposalJobData
type ProposalVotingOpenJobData ProposalJobData
type ProposalVotingEndingJobData ProposalJobData
type ProposalOutcomeJobData ProposalJobData
type ProposalGracePeriodJobData ProposalJobData
type ProposalFinalStateJobData ProposalJobData

// canceled proposal

// queued proposal

// abrogated proposal
type AbrogationProposalCreatedJobData AbrogationProposalJobData

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

type AbrogationProposalJobData struct {
	Id                    int64
	Proposer              string
	CreateTime            int64
	IncludedInBlockNumber int64
}

// proposal created
func NewProposalCreatedJob(data *ProposalCreatedJobData) (*Job, error) {
	return NewJob(ProposalCreated, 0, data.IncludedInBlockNumber, data)
}

func (jd *ProposalCreatedJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) (*Job, error) {
	log.Tracef("executing proposal created job for PID-%d", jd.Id)

	// send created notification
	err := saveNotification(
		ctx, tx,
		"system",
		ProposalCreated,
		jd.CreateTime,
		jd.CreateTime+jd.WarmUpDuration-300,
		fmt.Sprintf("Proposal PID-%d created by %s", jd.Id, jd.Proposer),
		nil,
		jd.IncludedInBlockNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "save create proposal notification to db")
	}

	// schedule job for next notification
	njd := ProposalActivatingJobData(*jd)
	next, err := NewProposalActivatingJob(&njd)
	if err != nil {
		return nil, errors.Wrap(err, "create create proposal next job")
	}

	return next, nil
}

// proposal voting starting soon
func NewProposalActivatingJob(data *ProposalActivatingJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration - 300
	return NewJob(ProposalActivating, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalActivatingJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) (*Job, error) {
	log.Tracef("executing proposal activated job for PID-%d", jd.Id)

	// check if proposal is still in warm up phase
	ps, err := proposalState(ctx, tx, jd.Id)
	if err != nil {
		return nil, err
	}
	if ps != ProposalStateWarmUp {
		log.Tracef("proposal PID-%d was not in WARMUP state but %s", jd.Id, ps)
		return nil, nil
	}

	// send voting starts notification
	err = saveNotification(
		ctx, tx,
		"system",
		ProposalActivating,
		jd.CreateTime+jd.WarmUpDuration-300,
		jd.CreateTime+jd.WarmUpDuration,
		fmt.Sprintf("Proposal PID-%d voting starts in 5 minutes", jd.Id),
		nil,
		jd.IncludedInBlockNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "save activating proposal notification to db")
	}

	// schedule job for next notification
	njd := ProposalVotingOpenJobData(*jd)
	next, err := NewProposalVotingOpenJob(&njd)
	if err != nil {
		return nil, errors.Wrap(err, "create create proposal next job")
	}

	return next, nil
}

// proposal started - voting open
func NewProposalVotingOpenJob(data *ProposalVotingOpenJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration
	return NewJob(ProposalVotingOpen, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalVotingOpenJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) (*Job, error) {
	log.Tracef("executing proposal voting open job for PID-%d", jd.Id)

	// check if proposal in active phase
	ps, err := proposalState(ctx, tx, jd.Id)
	if err != nil {
		return nil, err
	}
	if ps != ProposalStateActive {
		log.Tracef("proposal PID-%d was not in ACTIVE state but %s", jd.Id, ps)
		return nil, nil
	}

	// send voting started notification
	err = saveNotification(
		ctx, tx,
		"system",
		ProposalVotingOpen,
		jd.CreateTime+jd.WarmUpDuration,
		jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration-300,
		fmt.Sprintf("Proposal PID-%d voting period started, cast your vote now", jd.Id),
		nil,
		jd.IncludedInBlockNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "save proposal voting opened notification to db")
	}

	// schedule job for next notification
	njd := ProposalVotingEndingJobData(*jd)
	next, err := NewProposalVotingEndingJob(&njd)
	if err != nil {
		return nil, errors.Wrap(err, "create proposal voting open next job")
	}

	return next, nil
}

// voting ending soon
func NewProposalVotingEndingJob(data *ProposalVotingEndingJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration + data.ActiveDuration - 300
	return NewJob(ProposalVotingEnding, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalVotingEndingJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) (*Job, error) {
	log.Tracef("executing proposal voting ending job for PID-%d", jd.Id)

	// check if proposal in active phase
	ps, err := proposalState(ctx, tx, jd.Id)
	if err != nil {
		return nil, err
	}
	if ps != ProposalStateActive {
		log.Tracef("proposal PID-%d was not in ACTIVE state but %s", jd.Id, ps)
		return nil, nil
	}

	// send voting ending soon notification
	err = saveNotification(
		ctx, tx,
		"system",
		ProposalVotingEnding,
		jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration-300,
		jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration,
		fmt.Sprintf("Voting period for proposal PID-%d is ending soon, cast your vote now", jd.Id),
		nil,
		jd.IncludedInBlockNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "save proposal voting ending soon notification to db")
	}

	// schedule job for next notification
	njd := ProposalOutcomeJobData(*jd)
	next, err := NewProposalOutcomeJob(&njd)
	if err != nil {
		return nil, errors.Wrap(err, "create proposal voting ending next job")
	}

	return next, nil
}

// outcome of proposal voting period
func NewProposalOutcomeJob(data *ProposalOutcomeJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration + data.ActiveDuration
	return NewJob(ProposalOutcome, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalOutcomeJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) (*Job, error) {
	log.Tracef("executing proposal voting outcome job for PID-%d", jd.Id)

	// check if proposal in active phase
	ps, err := proposalState(ctx, tx, jd.Id)
	if err != nil {
		return nil, err
	}

	if ps == ProposalStateAccepted {
		// send proposal accepted notification
		err = saveNotification(
			ctx, tx,
			"system",
			ProposalAccepted,
			jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration,
			// TODO ? decide timings
			jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+jd.QueueDuration,
			fmt.Sprintf("Proposal PID-%d has been accepted and is queued for execution", jd.Id),
			nil,
			jd.IncludedInBlockNumber,
		)
		if err != nil {
			return nil, errors.Wrap(err, "save proposal accepted notification to db")
		}
	} else if ps == ProposalStateFailed {
		// send proposal failed notification
		err = saveNotification(
			ctx, tx,
			"system",
			ProposalFailed,
			jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration,
			// TODO ? decide timings
			jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+60*60*24,
			fmt.Sprintf("Proposal PID-%d failed to achive voting quorum", jd.Id),
			nil,
			jd.IncludedInBlockNumber,
		)
		if err != nil {
			return nil, errors.Wrap(err, "save proposal failed notification to db")
		}
	} else {
		log.Errorf("unknown proposal state after ending: PID-%d: %s", jd.Id, ps)
		return nil, nil
	}

	// schedule job for next notification
	njd := ProposalGracePeriodJobData(*jd)
	next, err := NewProposalGracePeriodJob(&njd)
	if err != nil {
		return nil, errors.Wrap(err, "create proposal voting outcome next job")
	}

	return next, nil
}

// proposal entering the grace period
func NewProposalGracePeriodJob(data *ProposalGracePeriodJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration + data.ActiveDuration + data.QueueDuration
	return NewJob(ProposalGracePeriod, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalGracePeriodJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) (*Job, error) {
	log.Tracef("executing proposal entering grace period job for PID-%d", jd.Id)

	// check if proposal in grace period
	ps, err := proposalState(ctx, tx, jd.Id)
	if err != nil {
		return nil, err
	}
	if ps != ProposalStateGracePeriod {
		log.Tracef("proposal PID-%d was not in GRACE state but %s", jd.Id, ps)
		return nil, nil
	}

	// send proposal in grace period notification
	err = saveNotification(
		ctx, tx,
		"system",
		ProposalGracePeriod,
		jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+jd.QueueDuration,
		jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+jd.QueueDuration+jd.GraceDuration,
		fmt.Sprintf("Proposal PID-%d can now be executed", jd.Id),
		nil,
		jd.IncludedInBlockNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "save proposal in grace period notification to db")
	}

	// TODO ? maybe we should schedule this from the get go
	// schedule job for next notification
	njd := ProposalFinalStateJobData(*jd)
	next, err := NewProposalFinalStateJob(&njd)
	if err != nil {
		return nil, errors.Wrap(err, "create proposal voting open next job")
	}

	return next, nil
}

// proposal execution result
func NewProposalFinalStateJob(data *ProposalFinalStateJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration + data.ActiveDuration + data.QueueDuration + data.GraceDuration
	return NewJob(ProposalFinalState, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalFinalStateJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) (*Job, error) {
	log.Tracef("executing proposal entering grace period job for PID-%d", jd.Id)

	// check if proposal in grace period
	ps, err := proposalState(ctx, tx, jd.Id)
	if err != nil {
		return nil, err
	}

	var msg string
	switch ps {
	case ProposalStateExecuted:
		msg = fmt.Sprintf("Proposal PID-%d has been executed sucessfully", jd.Id)
	case ProposalStateExpired:
		msg = fmt.Sprintf("Proposal PID-%d has not been executed and it expired", jd.Id)
	default:
		log.Tracef("proposal PID-%d was not in EXECUTED or EXPIRED state but %s", jd.Id, ps)
		return nil, nil
	}

	// send proposal in grace period notification
	err = saveNotification(
		ctx, tx,
		"system",
		ProposalFinalState,
		// TODO ? decide timings
		jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+jd.QueueDuration+jd.GraceDuration,
		jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+jd.QueueDuration+jd.GraceDuration+60*60*24,
		msg,
		nil,
		jd.IncludedInBlockNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "save proposal in grace period notification to db")
	}

	return nil, nil
}

// new abrogation proposal
func NewAbrogationProposalCreatedJob(data *AbrogationProposalCreatedJobData) (*Job, error) {
	return NewJob(AbrogationProposalCreated, 0, data.IncludedInBlockNumber, data)
}

func (jd *AbrogationProposalCreatedJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) (*Job, error) {
	log.Tracef("executing abrogation proposal created job for PID-%d", jd.Id)
	//
	// // send created notification
	// err := saveNotification(
	// 	ctx, tx,
	// 	"system",
	// 	ProposalCreated,
	// 	jd.CreateTime,
	// 	jd.CreateTime+jd.WarmUpDuration-300,
	// 	fmt.Sprintf("Proposal PID-%d created by %s", jd.Id, jd.Proposer),
	// 	nil,
	// 	jd.IncludedInBlockNumber,
	// )
	// if err != nil {
	// 	return nil, errors.Wrap(err, "save create proposal notification to db")
	// }
	//
	// // schedule job for next notification
	// njd := ProposalActivatingJobData(*jd)
	// next, err := NewProposalActivatingJob(&njd)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "create create proposal next job")
	// }
	//
	// return next, nil
	return nil, nil
}

func proposalState(ctx context.Context, tx *sql.Tx, Id int64) (string, error) {
	var ps string
	err := tx.QueryRowContext(ctx, "select * from proposal_state($1)", Id).Scan(&ps)
	if err != nil && err != sql.ErrNoRows {
		return ps, errors.Wrap(err, "get proposal state")
	}

	return ps, nil
}

func saveNotification(ctx context.Context, tx *sql.Tx, target string, typ string, starts int64, expires int64, msg string, metadata map[string]interface{}, block int64) error {
	notif := NewNotification(
		target,
		typ,
		starts,
		expires,
		msg,
		metadata,
		block,
	)

	return notif.ToDBWithTx(ctx, tx)
}
