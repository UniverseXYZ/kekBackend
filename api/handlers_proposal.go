package api

import (
	"database/sql"
	"fmt"
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
		createTime          int64
		targets             types2.JSONStringArray
		values              types2.JSONStringArray
		signatures          types2.JSONStringArray
		calldatas           types2.JSONStringArray
		timestamp           int64
		warmUpDuration      int64
		activeDuration      int64
		queueDuration       int64
		gracePeriodDuration int64
		acceptanceThreshold int64
		minQuorum           int64
		state               types.ProposalState
		forVotes            string
		againstVotes        string
	)

	err := a.core.DB().QueryRow(`
		select proposal_id,
			   proposer,
			   description,
			   title,
			   create_time,
			   targets,
			   "values",
			   signatures,
			   calldatas,
			   block_timestamp,
			   warm_up_duration,
			   active_duration,
			   queue_duration,
			   grace_period_duration,
			   acceptance_threshold,
			   min_quorum,
			   ( select * from proposal_state(proposal_id) ) as proposal_state,
		       coalesce(( select power from proposal_votes(proposal_id) where support = true ), 0) as for_votes,
			   coalesce(( select power from proposal_votes(proposal_id) where support = false ), 0) as against_votes
		from governance_proposals
		where proposal_ID = $1
	`, proposalID).Scan(&id, &proposer, &description, &title, &createTime, &targets, &values, &signatures, &calldatas, &timestamp, &warmUpDuration, &activeDuration, &queueDuration, &gracePeriodDuration, &acceptanceThreshold, &minQuorum, &state, &forVotes, &againstVotes)

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
		State:               state,
		StateTimeLeft:       getTimeLeft(state, createTime, warmUpDuration, activeDuration, queueDuration, gracePeriodDuration),
		ForVotes:            forVotes,
		AgainstVotes:        againstVotes,
	}

	OK(c, proposal)
}

func (a *API) AllProposalHandler(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	page := c.DefaultQuery("page", "1")
	title := c.DefaultQuery("title", "")
	proposalState := strings.ToUpper(c.DefaultQuery("state", "all"))

	if proposalState != "ALL" && !checkStateExist(proposalState) {
		BadRequest(c, errors.New("unknown state"))
		return
	}

	offset, err := calculateOffset(limit, page)
	if err != nil {
		Error(c, err)
		return
	}

	query := `
		select proposal_ID,
			   proposer,
			   description,
			   title,
			   create_time,
			   warm_up_duration,
			   active_duration,
			   queue_duration,
			   grace_period_duration,
			   ( select proposal_state(proposal_id) ) as proposal_state,
			   coalesce(( select power from proposal_votes(proposal_id) where support = true ), 0) as for_votes,
			   coalesce(( select power from proposal_votes(proposal_id) where support = false ), 0) as against_votes
		from governance_proposals
		where 1=1 
		%s %s
		order by proposal_id desc
		offset $1
		limit $2
	`

	var parameters = []interface{}{offset, limit}

	var stateFilter string
	if proposalState != "ALL" {
		parameters = append(parameters, proposalState)
		stateFilter = fmt.Sprintf("and ( select proposal_state(proposal_id) ) = $%d", len(parameters))
	}

	var titleFilter string
	if title != "" {
		parameters = append(parameters, "%"+strings.ToLower(title)+"%")
		titleFilter = fmt.Sprintf("and lower(title) like $%d", len(parameters))
	}

	query = fmt.Sprintf(query, stateFilter, titleFilter)

	rows, err := a.core.DB().Query(query, parameters...)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	defer rows.Close()

	var proposalList []types.ProposalLite

	for rows.Next() {
		var (
			id                  uint64
			proposer            string
			description         string
			title               string
			createTime          int64
			warmUpDuration      int64
			activeDuration      int64
			queueDuration       int64
			gracePeriodDuration int64
			state               string
			forVotes            string
			againstVotes        string
		)
		err := rows.Scan(&id, &proposer, &description, &title, &createTime, &warmUpDuration, &activeDuration, &queueDuration, &gracePeriodDuration, &state, &forVotes, &againstVotes)
		if err != nil {
			Error(c, err)
			return
		}

		proposal := types.ProposalLite{
			Id:            id,
			Proposer:      proposer,
			Description:   description,
			Title:         title,
			CreateTime:    createTime,
			State:         state,
			StateTimeLeft: getTimeLeft(types.ProposalState(state), createTime, warmUpDuration, activeDuration, queueDuration, gracePeriodDuration),
			ForVotes:      forVotes,
			AgainstVotes:  againstVotes,
		}

		proposalList = append(proposalList, proposal)
	}

	if len(proposalList) == 0 {
		NotFound(c)
		return
	}

	var count int
	var parameters2 []interface{}

	var stateFilter2 string
	if proposalState != "ALL" {
		parameters2 = append(parameters2, proposalState)
		stateFilter2 = fmt.Sprintf("and ( select proposal_state(proposal_id) ) = $%d", len(parameters2))
	}

	var titleFilter2 string
	if title != "" {
		parameters2 = append(parameters2, "%"+strings.ToLower(title)+"%")
		titleFilter2 = fmt.Sprintf("and lower(title) like $%d", len(parameters2))
	}

	var countQuery = `select count(*) from governance_proposals where 1=1 %s %s`

	err = a.core.DB().QueryRow(fmt.Sprintf(countQuery, stateFilter2, titleFilter2), parameters2...).Scan(&count)
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, proposalList, map[string]interface{}{"count": count})
}

func checkStateExist(state string) bool {
	proposalStates := [9]string{"WARMUP", "ACTIVE", "CANCELED", "FAILED", "ACCEPTED", "QUEUED", "GRACE", "EXPIRED", "EXECUTED"}
	for _, s := range proposalStates {
		if strings.ToLower(s) == strings.ToLower(state) {
			return true
		}

	}
	return false
}

func getTimeLeft(state types.ProposalState, createTime, warmUpDuration, activeDuration, queueDuration, gracePeriodDuration int64) *int64 {
	now := time.Now().Unix()
	var timeLeft int64

	switch state {
	case types.CANCELED, types.FAILED, types.ACCEPTED, types.EXPIRED, types.EXECUTED:
		return nil
	case types.WARMUP:
		timeLeft = createTime + warmUpDuration - now
	case types.ACTIVE:
		timeLeft = createTime + warmUpDuration + activeDuration - now
	case types.QUEUED:
		timeLeft = createTime + warmUpDuration + activeDuration + queueDuration - now
	case types.GRACE:
		timeLeft = createTime + warmUpDuration + activeDuration + queueDuration + gracePeriodDuration - now
	}

	return &timeLeft
}

func (a *API) handleProposalHistory(c *gin.Context) {
	proposalID := c.Param("proposalID")
	var proposal types.Proposal
	err := a.core.DB().QueryRow(`
		select
			   block_timestamp,
			   warm_up_duration,
			   active_duration,
			   queue_duration,
			   grace_period_duration,
			   ( select * from proposal_state(proposal_id) ) as proposal_state
		from governance_proposals
		where proposal_ID = $1
	`, proposalID).Scan(&proposal.BlockTimestamp, &proposal.WarmUpDuration, &proposal.ActiveDuration, &proposal.QueueDuration, &proposal.GracePeriodDuration, &proposal.State)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	eventsList, err := a.getALlEvents(proposalID)

	if err != nil {
		Error(c, err)
	}

	for _, event := range eventsList {
		if event.EventType == "CANCELED" {
			history, err := getHistoryBeforeCanceled(int64(event.CreateTime), proposal, eventsList)
			if err != nil {
				Error(c, err)
				return
			}
			OK(c, history)
			return
		}
	}

	var lastStateTs int64
	timeLeft := getTimeLeft(proposal.State, proposal.CreateTime, proposal.WarmUpDuration, proposal.ActiveDuration, proposal.QueueDuration, proposal.GracePeriodDuration)
	if timeLeft == nil {
		lastStateTs = time.Now().Unix()
	} else {
		lastStateTs = *timeLeft + time.Now().Unix()
	}

	history, err := getHistoryOfProposal(lastStateTs, proposal, eventsList)
	if err != nil {
		Error(c, err)
		return
	}
	OK(c, history)
	return
}

func (a *API) getALlEvents(proposalID string) ([]types.Event, error) {
	rows, err := a.core.DB().Query(`select proposal_id ,caller,event_type,event_data,block_timestamp from governance_events  where proposal_id = $1 order by block_timestamp`, proposalID)

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

func getHistoryOfProposal(lastStateTs int64, proposal types.Proposal, events []types.Event) ([]types.ProposalHistory, error) {
	var proposalHistory []types.ProposalHistory
	createdTs := proposal.CreateTime
	warmUpTs := createdTs + proposal.WarmUpDuration
	votingTs := warmUpTs + proposal.ActiveDuration
	queuedTs := votingTs + proposal.QueueDuration
	graceTs := queuedTs + proposal.GracePeriodDuration

	if lastStateTs > graceTs {
		execEvent := getExecutedEvent(events)
		if execEvent.CreateTime != 0 {
			proposalHistory = append(proposalHistory, types.ProposalHistory{
				ProposalState: types.EXECUTED,
				EndTime:       time.Unix(int64(lastStateTs), 0),
			})
		} else {
			proposalHistory = append(proposalHistory, types.ProposalHistory{
				ProposalState: types.EXPIRED,
				EndTime:       time.Unix(int64(lastStateTs), 0),
			})
		}
	}

	if lastStateTs > queuedTs {
		queuedEvent := getQueuedEvent(events)

		if time.Now().Unix() > graceTs && queuedEvent.Eta != nil {
			proposalHistory = append(proposalHistory, types.ProposalHistory{
				ProposalState: types.GRACE,
				EndTime:       time.Unix(int64(graceTs), 0),
			})
		}

		if queuedEvent.CreateTime == 0 {
			proposalHistory = append(proposalHistory, types.ProposalHistory{
				ProposalState: types.FAILED,
				EndTime:       time.Unix(int64(queuedEvent.CreateTime), 0),
			})

		} else {
			proposalHistory = append(proposalHistory, types.ProposalHistory{
				ProposalState: types.ACCEPTED,
				EndTime:       time.Unix(int64(queuedEvent.CreateTime), 0),
			}, types.ProposalHistory{
				ProposalState: types.QUEUED,
				EndTime:       time.Unix(queuedTs, 0),
			})
		}

	}

	if lastStateTs < queuedTs {
		proposalHistory = append(proposalHistory, types.ProposalHistory{
			ProposalState: "VOTING",
			EndTime:       time.Unix(votingTs, 0),
		})
	}

	if lastStateTs > votingTs {
		proposalHistory = append(proposalHistory, types.ProposalHistory{
			ProposalState: types.WARMUP,
			EndTime:       time.Unix(warmUpTs, 0),
		})
	}

	if lastStateTs > warmUpTs {
		proposalHistory = append(proposalHistory, types.ProposalHistory{
			ProposalState: types.CREATED,
			EndTime:       time.Unix(createdTs, 0),
		})
	}

	return proposalHistory, nil
}

func getQueuedEvent(events []types.Event) types.Event {
	var queuedEvent types.Event
	for _, event := range events {
		if event.EventType == "QUEUED" {
			return event
		}
	}
	return queuedEvent
}

func getExecutedEvent(events []types.Event) types.Event {
	var execEvent types.Event
	for _, event := range events {
		if event.EventType == "EXECUTED" {
			return event
		}
	}
	return execEvent
}
func getHistoryBeforeCanceled(canceledTS int64, proposal types.Proposal, events []types.Event) ([]types.ProposalHistory, error) {
	var proposalHistory []types.ProposalHistory
	createdTs := proposal.CreateTime
	warmUpTs := createdTs + proposal.WarmUpDuration
	votingTs := warmUpTs + proposal.ActiveDuration
	queuedTs := votingTs + proposal.QueueDuration
	graceTs := queuedTs + proposal.GracePeriodDuration

	proposalHistory = append(proposalHistory, types.ProposalHistory{
		ProposalState: types.CANCELED,
		EndTime:       time.Unix(int64(canceledTS), 0),
	})

	if canceledTS > queuedTs {
		queuedEvent := getQueuedEvent(events)

		if time.Now().Unix() > graceTs && queuedEvent.Eta != nil {
			proposalHistory = append(proposalHistory, types.ProposalHistory{
				ProposalState: types.GRACE,
				EndTime:       time.Unix(int64(graceTs), 0),
			})
		}

		if queuedEvent.CreateTime != 0 {
			proposalHistory = append(proposalHistory, types.ProposalHistory{
				ProposalState: types.QUEUED,
				EndTime:       time.Unix(queuedTs, 0),
			}, types.ProposalHistory{
				ProposalState: types.ACCEPTED,
				EndTime:       time.Unix(int64(queuedEvent.CreateTime), 0),
			})

		}
	}

	if canceledTS < queuedTs {
		proposalHistory = append(proposalHistory, types.ProposalHistory{
			ProposalState: "VOTING",
			EndTime:       time.Unix(votingTs, 0),
		})
	}

	if canceledTS > votingTs {
		proposalHistory = append(proposalHistory, types.ProposalHistory{
			ProposalState: types.WARMUP,
			EndTime:       time.Unix(warmUpTs, 0),
		})
	}

	if canceledTS > warmUpTs {
		proposalHistory = append(proposalHistory, types.ProposalHistory{
			ProposalState: types.CREATED,
			EndTime:       time.Unix(createdTs, 0),
		})
	}

	return proposalHistory, nil
}
