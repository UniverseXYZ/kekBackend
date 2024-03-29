package api

import (
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/kekDAO/kekBackend/api/types"
)

func (a *API) history(p types.Proposal) ([]types.HistoryEvent, error) {
	history, err := a.buildHistory(p)
	if err != nil {
		return nil, err
	}

	sort.Slice(history, func(i, j int) bool {
		if history[i].Name == string(types.CREATED) && history[j].Name == string(types.WARMUP) {
			return false
		} else if history[j].Name == string(types.CREATED) && history[i].Name == string(types.WARMUP) {
			return true
		}

		if history[i].Name == string(types.ACCEPTED) && history[j].Name == string(types.QUEUED) {
			return false
		} else if history[j].Name == string(types.ACCEPTED) && history[i].Name == string(types.QUEUED) {
			return true
		}

		return history[i].StartTs > history[j].StartTs
	})

	for i := 1; i <= len(history)-1; i++ {
		history[i].EndTs = history[i-1].StartTs - 1
	}

	history[0].EndTs = latestEventEndAt(p, history[0])

	return history, nil
}

func latestEventEndAt(p types.Proposal, event types.HistoryEvent) int64 {
	switch event.Name {
	case string(types.WARMUP):
		return event.StartTs + p.WarmUpDuration
	case string(types.ACTIVE):
		return event.StartTs + p.ActiveDuration
	case string(types.QUEUED):
		return event.StartTs + p.QueueDuration
	case string(types.GRACE):
		return event.StartTs + p.GracePeriodDuration
	default:
		return 0
	}
}

// we have events for the following actions: CREATED, QUEUED, CANCELED, EXECUTED
func (a *API) buildHistory(p types.Proposal) ([]types.HistoryEvent, error) {
	events, err := a.getProposalEvents(p.Id)
	if err != nil {
		return nil, errors.Wrap(err, "could not get proposal events")
	}

	var history []types.HistoryEvent

	history = append(history, types.HistoryEvent{
		Name:    string(types.CREATED),
		StartTs: p.CreateTime,
		TxHash:  events[0].TxHash,
	})

	sort.Slice(events, func(i, j int) bool {
		return events[i].CreateTime < events[j].CreateTime
	})
	events = events[1:]

	history = append(history, types.HistoryEvent{
		Name:    string(types.WARMUP),
		StartTs: p.CreateTime,
	})

	var nextDeadline int64

	// at this point, only a CANCELED event can occur that would be final for this history
	// so we check the list of events to see if there's any CANCELED event that occurred before the WARMUP period ended
	nextDeadline = p.CreateTime + p.WarmUpDuration
	if len(events) > 0 && events[0].CreateTime < nextDeadline && events[0].EventType == string(types.CANCELED) {
		history = append(history, types.HistoryEvent{
			Name:    string(types.CANCELED),
			StartTs: events[0].CreateTime,
			TxHash:  events[0].TxHash,
		})

		return history, nil
	}

	if nextDeadline >= time.Now().Unix() {
		return history, nil
	}

	// if the proposal was not canceled in the WARMUP period, then it means we reached ACTIVE state so we add that to
	// the history
	history = append(history, types.HistoryEvent{
		Name:    string(types.ACTIVE),
		StartTs: nextDeadline + 1,
	})

	// just like in WARMUP period, the only final event that could occur in this case is CANCELED
	nextDeadline = p.CreateTime + p.WarmUpDuration + p.ActiveDuration
	if len(events) > 0 && events[0].CreateTime < nextDeadline && events[0].EventType == string(types.CANCELED) {
		history = append(history, types.HistoryEvent{
			Name:    string(types.CANCELED),
			StartTs: events[0].CreateTime,
			TxHash:  events[0].TxHash,
		})

		return history, nil
	}

	if nextDeadline >= time.Now().Unix() {
		return history, nil
	}

	// if the proposal was not canceled during the ACTIVE period, that means it reached one of ACCEPTED/FAILED
	// for this we need to check the votes
	failed, err := isFailedProposal(p)
	if err != nil {
		return nil, errors.Wrap(err, "could not check if proposal failed")
	}

	if failed {
		history = append(history, types.HistoryEvent{
			Name:    string(types.FAILED),
			StartTs: nextDeadline + 1,
		})

		return history, nil
	} else {
		history = append(history, types.HistoryEvent{
			Name:    string(types.ACCEPTED),
			StartTs: nextDeadline + 1,
		})
	}

	// after the proposal reached accepted state, nothing else can happen unless somebody calls the queue function
	// which emits a QUEUED event
	if len(events) == 0 {
		return history, nil
	}

	if events[0].EventType == string(types.QUEUED) {
		history = append(history, types.HistoryEvent{
			Name:    string(types.QUEUED),
			StartTs: p.CreateTime + p.WarmUpDuration + p.ActiveDuration + 1,
			TxHash:  events[0].TxHash,
		})
	} else {
		// the next event must be QUEUED, but just in case ...
		return history, nil
	}

	events = events[1:]

	nextDeadline = p.CreateTime + p.WarmUpDuration + p.ActiveDuration + p.QueueDuration
	if nextDeadline >= time.Now().Unix() {
		return history, nil
	}

	// at this point the queue period passed and we must check if there was a abrogation proposal that succeeded
	hasCP, err := a.abrogationProposalExists(p)
	if err != nil {
		return nil, err
	}

	if hasCP {
		forVotes, kekStaked, err := a.abrogationProposalData(p)
		if err != nil {
			return nil, err
		}

		passed, err := abrogationProposalPassed(forVotes, kekStaked)
		if err != nil {
			return nil, errors.Wrap(err, "could not check if abrogation proposal passed")
		}

		if passed {
			history = append(history, types.HistoryEvent{
				Name:    string(types.ABROGATED),
				StartTs: nextDeadline,
			})

			return history, nil
		}
	}

	// no abrogation proposal or did not pass
	history = append(history, types.HistoryEvent{
		Name:    string(types.GRACE),
		StartTs: nextDeadline,
	})

	nextDeadline = nextDeadline + p.GracePeriodDuration
	if len(events) > 0 && events[0].CreateTime <= nextDeadline && events[0].EventType == string(types.EXECUTED) {
		history = append(history, types.HistoryEvent{
			Name:    string(types.EXECUTED),
			StartTs: events[0].CreateTime,
			TxHash:  events[0].TxHash,
		})

		return history, nil
	}

	if nextDeadline >= time.Now().Unix() {
		return history, nil
	}

	history = append(history, types.HistoryEvent{
		Name:    string(types.EXPIRED),
		StartTs: nextDeadline,
	})

	return history, nil
}

func (a *API) abrogationProposalData(p types.Proposal) (string, string, error) {
	var forVotes, kekStaked string
	err := a.db.QueryRow(`
		select 
		       coalesce(( select sum(power) from abrogation_proposal_votes($1) where support = true ), 0) as for_votes, 
		       kek_staked_at_ts(to_timestamp((select create_time from governance_abrogation_proposals where proposal_id = $1))) as kek_staked
	`, p.Id).Scan(&forVotes, &kekStaked)

	if err != nil {
		return "", "", errors.Wrap(err, "could not scan number of votes on abrogation proposal")
	}

	return forVotes, kekStaked, nil
}

func (a *API) abrogationProposalExists(p types.Proposal) (bool, error) {
	var count int64
	err := a.db.QueryRow(`select count(*) from governance_abrogation_proposals where proposal_id = $1`, p.Id).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "could not check for abrogation proposal")
	}

	return count > 0, nil
}

func abrogationProposalPassed(forVotes string, kekStaked string) (bool, error) {
	pro, err := decimal.NewFromString(forVotes)
	if err != nil {
		return false, errors.Wrap(err, "could not convert forVotes to decimal")
	}

	b, err := decimal.NewFromString(kekStaked)
	if err != nil {
		return false, errors.Wrap(err, "could not convert kekStaked to decimal")
	}

	return pro.GreaterThan(b.DivRound(decimal.NewFromInt(2), 18)), nil
}

func isFailedProposal(p types.Proposal) (bool, error) {
	pro, err := decimal.NewFromString(p.ForVotes)
	if err != nil {
		return false, errors.Wrap(err, "could not convert forVotes to decimal")
	}

	against, err := decimal.NewFromString(p.AgainstVotes)
	if err != nil {
		return false, errors.Wrap(err, "could not convert againstVotes to decimal")
	}

	kekStaked, err := decimal.NewFromString(p.KekStaked)
	if err != nil {
		return false, errors.Wrap(err, "could not convert kekStaked to decimal")
	}

	minQuorum := decimal.NewFromInt(p.MinQuorum)
	acceptance := decimal.NewFromInt(p.AcceptanceThreshold)

	if pro.Add(against).LessThan(kekStaked.Mul(minQuorum).DivRound(decimal.NewFromInt(100), 18)) {
		return true, nil
	}

	total := pro.Add(against)
	minForVotes := total.Mul(acceptance).DivRound(decimal.NewFromInt(100), 18)

	return pro.LessThanOrEqual(minForVotes), nil
}
