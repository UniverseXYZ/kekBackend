package api

import (
	"database/sql"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/api/types"
	types2 "github.com/barnbridge/barnbridge-backend/types"
)

func (a *API) ProposalDetailsHandler(c *gin.Context) {
	proposalID := c.Param("proposalID")

	var (
		id                  uint64
		proposer            string
		description         string
		title               string
		createTime          uint64
		targets             types2.JSONStringArray
		values              types2.JSONStringArray
		signatures          types2.JSONStringArray
		calldatas           types2.JSONStringArray
		timestamp           int64
		warmUpDuration      uint64
		activeDuration      uint64
		queueDuration       uint64
		gracePeriodDuration uint64
		acceptanceThreshold uint64
		minQuorum           uint64
	)
	err := a.core.DB().QueryRow(`select proposal_id,proposer,description,title,create_time,targets,"values",signatures,calldatas,block_timestamp,warm_up_duration, active_duration ,
       queue_duration ,  grace_period_duration, acceptance_threshold, min_quorum 
       from governance_proposals where proposal_ID = $1 `, proposalID).Scan(&id, &proposer, &description, &title, &createTime, &targets, &values, &signatures,
		&calldatas, &timestamp, &warmUpDuration, &activeDuration, &queueDuration, &gracePeriodDuration, &acceptanceThreshold, &minQuorum)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	proposal := types.Proposal{
		Id:                  id,
		Proposer:            proposer,
		Description:         description,
		Title:               title,
		CreateTime:          createTime,
		Targets:             targets,
		Values:              values,
		Signatures:          signatures,
		Calldatas:           calldatas,
		BlockTimestamp:      timestamp,
		WarmUpDuration:      warmUpDuration,
		ActiveDuration:      activeDuration,
		QueueDuration:       queueDuration,
		GracePeriodDuration: gracePeriodDuration,
		AcceptanceThreshold: acceptanceThreshold,
		MinQuorum:           minQuorum,
	}
	state, err := a.calculateState(proposal)
	if err != nil {
		Error(c, err)
		return
	}
	proposal.State = state

	OK(c, proposal)
}

func (a *API) AllProposalHandler(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	offset := c.DefaultQuery("offset", strconv.FormatInt(math.MaxInt32, 10))
	title := c.DefaultQuery("title", "")
	proposalState := strings.ToLower(c.DefaultQuery("state", "all"))

	var rows *sql.Rows
	var err error

	if proposalState != "all" && !checkStateExist(proposalState) {
		BadRequest(c, errors.New("unknown state"))
		return
	}

	if title == "" {
		rows, err = a.core.DB().Query(`
			select proposal_ID, proposer, description, title, create_time, targets, "values", signatures, calldatas, block_timestamp,warm_up_duration, active_duration ,
       			queue_duration ,  grace_period_duration, acceptance_threshold, min_quorum
			from governance_proposals 
			where proposal_id <= $1 
			order by proposal_id desc 
			limit $2
		`, offset, limit)
	} else {
		title = "%" + strings.ToLower(title) + "%"
		rows, err = a.core.DB().Query(`
			select proposal_ID, proposer, description, title, create_time, targets, "values", signatures, calldatas,block_timestamp,warm_up_duration, active_duration ,
       queue_duration ,  grace_period_duration, acceptance_threshold, min_quorum 
			from governance_proposals 
			where proposal_id <= $1 
			  and lower(title) like $2 
			order by proposal_id desc 
			limit $3
		`, offset, title, limit)
	}

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	defer rows.Close()

	var proposalList []types.Proposal

	for rows.Next() {
		var (
			id                  uint64
			proposer            string
			description         string
			title               string
			createTime          uint64
			targets             types2.JSONStringArray
			values              types2.JSONStringArray
			signatures          types2.JSONStringArray
			calldatas           types2.JSONStringArray
			timestamp           int64
			warmUpDuration      uint64
			activeDuration      uint64
			queueDuration       uint64
			gracePeriodDuration uint64
			acceptanceThreshold uint64
			minQuorum           uint64
		)
		err := rows.Scan(&id, &proposer, &description, &title, &createTime, &targets, &values, &signatures, &calldatas, &timestamp, &warmUpDuration, &activeDuration, &queueDuration, &gracePeriodDuration, &acceptanceThreshold, &minQuorum)
		if err != nil {
			Error(c, err)
			return
		}

		proposal := types.Proposal{
			Id:                  id,
			Proposer:            proposer,
			Description:         description,
			Title:               title,
			CreateTime:          createTime,
			Targets:             targets,
			Values:              values,
			Signatures:          signatures,
			Calldatas:           calldatas,
			BlockTimestamp:      timestamp,
			WarmUpDuration:      warmUpDuration,
			ActiveDuration:      activeDuration,
			QueueDuration:       queueDuration,
			GracePeriodDuration: gracePeriodDuration,
			AcceptanceThreshold: acceptanceThreshold,
			MinQuorum:           minQuorum,
		}
		proposal.State, err = a.calculateState(proposal)
		if err != nil {
			Error(c, err)
			return
		}

		if proposalState == "all" || proposal.State == proposalState {
			proposalList = append(proposalList, proposal)
		}
	}

	if len(proposalList) == 0 {
		NotFound(c)
		return
	}

	var count int
	err = a.core.DB().QueryRow(`select count(*) from governance_proposals`).Scan(&count)
	if err != nil {
		Error(c, err)
	}

	OK(c, proposalList, map[string]interface{}{"count": count})
}

func (a *API) calculateState(p types.Proposal) (string, error) {
	now := time.Now()
	timestampNow := uint64(now.Unix())
	warmUpDuration := p.CreateTime + p.WarmUpDuration
	activeDuration := p.CreateTime + p.WarmUpDuration + p.ActiveDuration
	queuedDuration := p.CreateTime + p.WarmUpDuration + p.ActiveDuration + p.QueueDuration
	graceDuration := queuedDuration + p.GracePeriodDuration

	events, err := a.getAllEvents(p.Id)
	if err != nil {
		return "", err
	}

	existCanceled := a.checkForCanceledEvent(events)
	if existCanceled {
		return "Canceled", nil
	}

	existExecuted := a.checkForExecutedEvent(events)
	if existExecuted {
		return "Executed", nil
	}

	if timestampNow <= warmUpDuration {
		return "WarmUp", nil
	}

	if timestampNow <= activeDuration {
		return "Active", nil
	}

	if timestampNow < queuedDuration {
		existQueued := a.checkForQueuedEvent(events)
		if existQueued {
			return "Queued", nil
		} else {
			state, err := a.checkVotesForProposal(p.Id)
			if err != nil {
				return "", err
			}
			return state, err
		}
	}

	if timestampNow < graceDuration {
		return "Grace", nil
	}

	return "Expired", nil

}
func (a *API) getAllEvents(proposalID uint64) ([]types.Event, error) {

	rows, err := a.core.DB().Query(`select proposal_id ,caller,event_type,event_data,block_timestamp from governance_events where proposal_id = $1`, proposalID)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	var eventsList []types.Event

	for rows.Next() {
		var event types.Event
		err := rows.Scan(&event.ProposalID, &event.Caller, &event.EventType, &event.Eta, &event.CreateTime)
		if err != nil {

			return nil, err
		}

		eventsList = append(eventsList, event)
	}
	return eventsList, nil
}

func (a *API) checkForQueuedEvent(events []types.Event) bool {
	for _, event := range events {
		if strings.ToLower(event.EventType) == strings.ToLower("QUEUED") {
			return true
		}
	}
	return false
}

func (a *API) checkForCanceledEvent(events []types.Event) bool {
	for _, event := range events {
		if strings.ToLower(event.EventType) == strings.ToLower("CANCELED") {
			return true
		}
	}
	return false
}

func (a *API) checkForExecutedEvent(events []types.Event) bool {
	for _, event := range events {
		if strings.ToLower(event.EventType) == strings.ToLower("EXECUTED") {
			return true
		}
	}
	return false
}

func (a *API) checkVotesForProposal(proposalID uint64) (string, error) {
	var forVotes, againstVotes int
	var state string
	err := a.core.DB().QueryRow(`select count(*) from proposal_votes($1) where support = 'true'`, proposalID).Scan(&forVotes)
	if err != nil && err != sql.ErrNoRows {
		return state, err
	}
	err = a.core.DB().QueryRow(`select count(*) from proposal_votes($1) where support = 'false'`, proposalID).Scan(&againstVotes)
	if err != nil && err != sql.ErrNoRows {
		return state, err
	}
	if forVotes < againstVotes {
		state = "Failed"
	} else {
		state = "Accepted"
	}
	return state, nil
}

func checkStateExist(state string) bool {
	proposalStates := [9]string{"WarmUp", "Active", "Canceled", "Failed", "Accepted", "Queued", "Grace", "Expired", "Executed"}
	for _, s := range proposalStates {
		if strings.ToLower(s) == strings.ToLower(state) {
			return true
		}

	}
	return false
}
