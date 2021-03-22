package notifications

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

const (
	ProposalCreated           = "proposal-created"
	ProposalActivating        = "proposal-activating-soon"
	ProposalCanceled          = "proposal-canceled"
	ProposalVotingOpen        = "proposal-voting-open"
	ProposalVotingEnding      = "proposal-voting-ending-soon"
	ProposalOutcome           = "proposal-outcome"
	ProposalAccepted          = "proposal-accepted"
	ProposalFailed            = "proposal-failed"
	ProposalQueued            = "proposal-queued"
	ProposalGracePeriod       = "proposal-grace"
	ProposalExecuted          = "proposal-executed"
	ProposalExpired           = "proposal-expired"
	AbrogationProposalCreated = "abrogation-proposal-created"
	ProposalAbrogated         = "proposal-abrogated"
)

const (
	ProposalStateWarmUp      = "WARMUP"
	ProposalStateActive      = "ACTIVE"
	ProposalStateCanceled    = "CANCELED"
	ProposalStateAccepted    = "ACCEPTED"
	ProposalStateFailed      = "FAILED"
	ProposalStateQueued      = "QUEUED"
	ProposalStateGracePeriod = "GRACE"
	ProposalStateExecuted    = "EXECUTED"
	ProposalStateAbrogated   = "ABROGATED"
	ProposalStateExpired     = "EXPIRED"
)

// new proposal
type ProposalCreatedJobData ProposalJobData
type ProposalActivatingJobData ProposalJobData
type ProposalCanceledJobData ProposalEventsJobData
type ProposalVotingOpenJobData ProposalJobData
type ProposalVotingEndingJobData ProposalJobData
type ProposalOutcomeJobData ProposalJobData
type ProposalQueuedJobData ProposalEventsJobData
type ProposalGracePeriodJobData ProposalJobData
type ProposalExpiredJobData ProposalJobData
type ProposalExecutedJobData ProposalEventsJobData
type AbrogationProposalCreatedJobData AbrogationProposalJobData
type ProposalAbrogatedJobData ProposalJobData

type ProposalJobData struct {
	Id                    int64  `json:"proposalId"`
	Proposer              string `json:"proposer"`
	Title                 string `json:"title"`
	CreateTime            int64  `json:"createTime"`
	WarmUpDuration        int64  `json:"warmUpDuration"`
	ActiveDuration        int64  `json:"activeDuration"`
	QueueDuration         int64  `json:"queueDuration"`
	GraceDuration         int64  `json:"graceDuration"`
	IncludedInBlockNumber int64  `json:"includedInBlockNumber"`
}

type AbrogationProposalJobData struct {
	Id                    int64  `json:"proposalId"`
	Proposer              string `json:"proposer"`
	CreateTime            int64  `json:"createTime"`
	IncludedInBlockNumber int64  `json:"includedInBlockNumber"`
}

type ProposalEventsJobData struct {
	Id                    int64 `json:"proposalId"`
	CreateTime            int64 `json:"createTime"`
	IncludedInBlockNumber int64 `json:"includedInBlockNumber"`
}

// proposal created
func NewProposalCreatedJob(data *ProposalCreatedJobData) (*Job, error) {
	return NewJob(ProposalCreated, 0, data.IncludedInBlockNumber, data)
}

func (jd *ProposalCreatedJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) ([]*Job, error) {
	log.Tracef("executing proposal created job for PID-%d", jd.Id)

	// send created notification
	err := saveNotification(
		ctx, tx,
		"system",
		ProposalCreated,
		jd.CreateTime,
		jd.CreateTime+jd.WarmUpDuration-300,
		fmt.Sprintf("Proposal PID-%d created by %s", jd.Id, jd.Proposer),
		jobDataMetadata((*ProposalJobData)(jd)),
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

	return []*Job{
		next,
	}, nil
}

// proposal voting starting soon
func NewProposalActivatingJob(data *ProposalActivatingJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration - 300
	return NewJob(ProposalActivating, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalActivatingJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) ([]*Job, error) {
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
		jobDataMetadata((*ProposalJobData)(jd)),
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

	return []*Job{
		next,
	}, nil
}

// proposal started - voting open
func NewProposalVotingOpenJob(data *ProposalVotingOpenJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration
	return NewJob(ProposalVotingOpen, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalVotingOpenJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) ([]*Job, error) {
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
		jobDataMetadata((*ProposalJobData)(jd)),
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

	return []*Job{
		next,
	}, nil
}

// voting ending soon
func NewProposalVotingEndingJob(data *ProposalVotingEndingJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration + data.ActiveDuration - 300
	return NewJob(ProposalVotingEnding, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalVotingEndingJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) ([]*Job, error) {
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
		jobDataMetadata((*ProposalJobData)(jd)),
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

	return []*Job{
		next,
	}, nil
}

// outcome of proposal voting period
func NewProposalOutcomeJob(data *ProposalOutcomeJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration + data.ActiveDuration + 180 // delay to make sure we are free of reorgs
	return NewJob(ProposalOutcome, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalOutcomeJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) ([]*Job, error) {
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
			jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+180,
			// TODO ? decide timings
			jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+jd.QueueDuration,
			fmt.Sprintf("Proposal PID-%d has been accepted and is awaiting queuing for execution", jd.Id),
			jobDataMetadata((*ProposalJobData)(jd)),
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
			jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+180,
			// TODO ? decide timings
			jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+60*60*24,
			fmt.Sprintf("Proposal PID-%d failed to pass", jd.Id),
			jobDataMetadata((*ProposalJobData)(jd)),
			jd.IncludedInBlockNumber,
		)
		if err != nil {
			return nil, errors.Wrap(err, "save proposal failed notification to db")
		}
	} else {
		log.Tracef("proposal PID-%d state after ending: %s", jd.Id, ps)
		return nil, nil
	}

	return nil, nil
}

// proposal entering the grace period
func NewProposalGracePeriodJob(data *ProposalGracePeriodJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration + data.ActiveDuration + data.QueueDuration
	return NewJob(ProposalGracePeriod, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalGracePeriodJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) ([]*Job, error) {
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
		jobDataMetadata((*ProposalJobData)(jd)),
		jd.IncludedInBlockNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "save proposal in grace period notification to db")
	}

	// schedule job for next notification
	njd := ProposalExpiredJobData(*jd)
	next, err := NewProposalExpiredJob(&njd)
	if err != nil {
		return nil, errors.Wrap(err, "create proposal voting open next job")
	}

	return []*Job{
		next,
	}, nil
}

// proposal expired - scheduled from the start as a fallback
func NewProposalExpiredJob(data *ProposalExpiredJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration + data.ActiveDuration + data.QueueDuration + data.GraceDuration + 180 // delay to make sure we are free of reorgs
	return NewJob(ProposalExpired, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalExpiredJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) ([]*Job, error) {
	log.Tracef("executing proposal expired job for PID-%d", jd.Id)

	ps, err := proposalState(ctx, tx, jd.Id)
	if err != nil {
		return nil, err
	}
	if ps != ProposalStateExpired {
		log.Tracef("proposal PID-%d was not in EXPIRED state but %s", jd.Id, ps)
		return nil, nil
	}

	// send proposal expired notification
	err = saveNotification(
		ctx, tx,
		"system",
		ProposalExpired,
		// TODO ? decide timings
		jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+jd.QueueDuration+jd.GraceDuration+180,
		jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+jd.QueueDuration+jd.GraceDuration+60*60*24,
		fmt.Sprintf("Proposal PID-%d has expired", jd.Id),
		jobDataMetadata((*ProposalJobData)(jd)),
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

func (jd *AbrogationProposalCreatedJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) ([]*Job, error) {
	log.Tracef("executing abrogation proposal created job for PID-%d", jd.Id)

	// get the original proposal
	pjd, err := proposalAsJobData(ctx, tx, jd.Id)
	if err != nil {
		return nil, errors.Wrap(err, "proposal as job data")
	}
	// TODO should this be  fatal?
	if pjd == nil {
		log.Errorf("proposal PID-%d was not found but we have an abrogated event", jd.Id)
		return nil, nil
	}
	// setting these to abrogation proposal details in case we need to recall them later
	pjd.Proposer = jd.Proposer
	pjd.IncludedInBlockNumber = jd.IncludedInBlockNumber

	// send abrogation proposal created notification
	err = saveNotification(
		ctx, tx,
		"system",
		AbrogationProposalCreated,
		jd.CreateTime,
		pjd.CreateTime+pjd.WarmUpDuration+pjd.ActiveDuration+pjd.QueueDuration, // TODO see about timings
		fmt.Sprintf("Abrogation proposal for PID-%d created by %s", jd.Id, jd.Proposer),
		jobDataMetadata(pjd),
		jd.IncludedInBlockNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "save create abrogation proposal notification to db")
	}

	// schedule job for next notification
	// NOTE we are replacing the abrogation proposal with the original proposal + some overwritten fields
	njd := ProposalAbrogatedJobData(*pjd)
	next, err := NewProposalAbrogatedJob(&njd)
	if err != nil {
		return nil, errors.Wrap(err, "create abrogation proposal next job")
	}

	return []*Job{
		next,
	}, nil
}

func NewProposalAbrogatedJob(data *ProposalAbrogatedJobData) (*Job, error) {
	x := data.CreateTime + data.WarmUpDuration + data.ActiveDuration + data.QueueDuration + 180 // delay for safety against reorgs
	return NewJob(ProposalAbrogated, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalAbrogatedJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) ([]*Job, error) {
	log.Tracef("executing abrogated proposal job for PID-%d", jd.Id)

	ps, err := proposalState(ctx, tx, jd.Id)
	if err != nil {
		return nil, err
	}
	if ps != ProposalStateAbrogated {
		log.Tracef("proposal PID-%d was not in ABROGATED state but %s", jd.Id, ps)
		return nil, nil
	}

	// send abrogation proposal created notification
	err = saveNotification(
		ctx, tx,
		"system",
		ProposalAbrogated,
		jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+jd.QueueDuration+180,
		jd.CreateTime+jd.WarmUpDuration+jd.ActiveDuration+jd.QueueDuration+60*60*24, // TODO see about timings
		fmt.Sprintf("Proposal PID-%d has been abrogated", jd.Id),
		jobDataMetadata((*ProposalJobData)(jd)),
		jd.IncludedInBlockNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "save abrogated proposal notification to db")
	}

	return nil, nil
}

// events

// proposal canceled
func NewProposalCanceledJob(data *ProposalCanceledJobData) (*Job, error) {
	x := data.CreateTime + 180 // delay for safety against reorgs
	return NewJob(ProposalCanceled, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalCanceledJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) ([]*Job, error) {
	log.Tracef("executing proposal canceled for PID-%d", jd.Id)

	ps, err := proposalState(ctx, tx, jd.Id)
	if err != nil {
		return nil, err
	}
	if ps != ProposalStateCanceled {
		log.Tracef("proposal PID-%d was not in CANCELED state but %s", jd.Id, ps)
		return nil, nil
	}

	// send proposal canceled notification
	err = saveNotification(
		ctx, tx,
		"system",
		ProposalCanceled,
		jd.CreateTime+180,
		jd.CreateTime+60*60*24, // TODO see about timings
		fmt.Sprintf("Proposal PID-%d has been canceled", jd.Id),
		eventJobDataMetadata((*ProposalEventsJobData)(jd)),
		jd.IncludedInBlockNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "save proposal canceled notification to db")
	}

	return nil, nil
}

// proposal queued
func NewProposalQueuedJob(data *ProposalQueuedJobData) (*Job, error) {
	x := data.CreateTime
	return NewJob(ProposalQueued, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalQueuedJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) ([]*Job, error) {
	log.Tracef("executing proposal queued for PID-%d", jd.Id)

	ps, err := proposalState(ctx, tx, jd.Id)
	if err != nil {
		return nil, err
	}
	if ps != ProposalStateQueued {
		log.Tracef("proposal PID-%d was not in QUEUED state but %s", jd.Id, ps)
		return nil, nil
	}

	pjd, err := proposalAsJobData(ctx, tx, jd.Id)
	if err != nil {
		return nil, errors.Wrap(err, "proposal as job data")
	}
	// TODO should this be  fatal?
	if pjd == nil {
		log.Errorf("proposal PID-%d was not found but we have a queued event", jd.Id)
		return nil, nil
	}
	pjd.IncludedInBlockNumber = jd.IncludedInBlockNumber

	// send proposal queued notification
	err = saveNotification(
		ctx, tx,
		"system",
		ProposalQueued,
		jd.CreateTime,
		pjd.CreateTime+pjd.WarmUpDuration+pjd.ActiveDuration+pjd.QueueDuration,
		fmt.Sprintf("Proposal PID-%d has been queued for execution", jd.Id),
		eventJobDataMetadata((*ProposalEventsJobData)(jd)),
		jd.IncludedInBlockNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "save proposal queued notification to db")
	}

	// schedule job for next notification
	njd := ProposalGracePeriodJobData(*pjd)
	next, err := NewProposalGracePeriodJob(&njd)
	if err != nil {
		return nil, errors.Wrap(err, "create proposal queued next job")
	}

	return []*Job{
		next,
	}, nil
}

// proposal executed
func NewProposalExecutedJob(data *ProposalExecutedJobData) (*Job, error) {
	x := data.CreateTime + 180 // delay for safety against reorgs
	return NewJob(ProposalExecuted, x, data.IncludedInBlockNumber, data)
}

func (jd *ProposalExecutedJobData) ExecuteWithTx(ctx context.Context, tx *sql.Tx) ([]*Job, error) {
	log.Tracef("executing proposal executed for PID-%d", jd.Id)

	ps, err := proposalState(ctx, tx, jd.Id)
	if err != nil {
		return nil, err
	}
	if ps != ProposalStateExecuted {
		log.Tracef("proposal PID-%d was not in EXECUTED state but %s", jd.Id, ps)
		return nil, nil
	}

	// send proposal executed notification
	err = saveNotification(
		ctx, tx,
		"system",
		ProposalExecuted,
		jd.CreateTime+180,
		jd.CreateTime+60*60*24, // TODO see about timings
		fmt.Sprintf("Proposal PID-%d has been executed", jd.Id),
		eventJobDataMetadata((*ProposalEventsJobData)(jd)),
		jd.IncludedInBlockNumber,
	)
	if err != nil {
		return nil, errors.Wrap(err, "save proposal executed notification to db")
	}

	return nil, nil
}

func proposalState(ctx context.Context, tx *sql.Tx, Id int64) (string, error) {
	var ps string
	sel := `SELECT * FROM proposal_state($1);`
	err := tx.QueryRowContext(ctx, sel, Id).Scan(&ps)
	if err != nil && err != sql.ErrNoRows {
		return ps, errors.Wrap(err, "get proposal state")
	}

	return ps, nil
}

func proposalAsJobData(ctx context.Context, tx *sql.Tx, id int64) (*ProposalJobData, error) {
	var pjd ProposalJobData

	query := `
		SELECT
			"proposal_id",
			"proposer",
			"title",
			"create_time",
			"warm_up_duration",
			"active_duration",
			"queue_duration",
			"grace_period_duration",
		    "included_in_block"
		FROM
			"governance_proposals"
		WHERE
			"proposal_id" = $1
	;
	`

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&pjd.Id, &pjd.Proposer, &pjd.Title,
		&pjd.CreateTime, &pjd.WarmUpDuration, &pjd.ActiveDuration, &pjd.QueueDuration, &pjd.GraceDuration,
		&pjd.IncludedInBlockNumber,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrapf(err, "get proposal as job data %d", id)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &pjd, nil
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

func jobDataMetadata(jd *ProposalJobData) map[string]interface{} {
	m := make(map[string]interface{})
	m["proposalId"] = jd.Id
	m["proposer"] = jd.Proposer
	return m
}

func eventJobDataMetadata(jd *ProposalEventsJobData) map[string]interface{} {
	m := make(map[string]interface{})
	m["proposalId"] = jd.Id
	return m
}
