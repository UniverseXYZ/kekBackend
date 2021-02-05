package api

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
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
		forVotes            string
		againstVotes        string
		bondStaked          string
		state               types.ProposalState
	)

	err := a.db.QueryRow(`
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
		       coalesce(( select sum(power) from proposal_votes(proposal_id) where support = true ), 0) as for_votes,
			   coalesce(( select sum(power) from proposal_votes(proposal_id) where support = false ), 0) as against_votes,
		       coalesce(( select bond_staked_at_ts(to_timestamp(create_time+warm_up_duration)) ), 0) as bond_staked,
			   ( select * from proposal_state(proposal_id) ) as proposal_state
		from governance_proposals
		where proposal_ID = $1
	`, proposalID).Scan(&id, &proposer, &description, &title, &createTime, &targets, &values, &signatures, &calldatas, &timestamp, &warmUpDuration, &activeDuration, &queueDuration, &gracePeriodDuration, &acceptanceThreshold, &minQuorum, &forVotes, &againstVotes, &bondStaked, &state)

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
		BondStaked:          bondStaked,
	}

	history, err := a.history(proposal)
	if err != nil {
		Error(c, err)
		return
	}

	proposal.History = history

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, proposal, map[string]interface{}{
		"block": block,
	})
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
			   coalesce(( select sum(power) from proposal_votes(proposal_id) where support = true ), 0) as for_votes,
			   coalesce(( select sum(power) from proposal_votes(proposal_id) where support = false ), 0) as against_votes
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
		if proposalState == "ACTIVE" {
			parameters = append(parameters, pq.Array([]string{"WARMUP", "ACTIVE", "ACCEPTED", "QUEUED", "GRACE"}))
		} else if proposalState == "FAILED" {
			parameters = append(parameters, pq.Array([]string{"CANCELED", "FAILED", "ABROGATED", "EXPIRED"}))
		} else {
			parameters = append(parameters, pq.Array([]string{proposalState}))
		}

		stateFilter = fmt.Sprintf("and (select proposal_state(proposal_id) ) = ANY($%d)", len(parameters))
	}

	var titleFilter string
	if title != "" {
		parameters = append(parameters, "%"+strings.ToLower(title)+"%")
		titleFilter = fmt.Sprintf("and lower(title) like $%d", len(parameters))
	}

	query = fmt.Sprintf(query, stateFilter, titleFilter)

	rows, err := a.db.Query(query, parameters...)
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

	err = a.db.QueryRow(fmt.Sprintf(countQuery, stateFilter2, titleFilter2), parameters2...).Scan(&count)
	if err != nil {
		Error(c, err)
		return
	}

	block, err := a.getHighestBlock()
	if err != nil {
		Error(c, err)
		return
	}

	OK(c, proposalList, map[string]interface{}{"count": count, "block": block})
}

func checkStateExist(state string) bool {
	proposalStates := []types.ProposalState{types.WARMUP, types.ACTIVE, types.CANCELED, types.FAILED, types.ACCEPTED, types.QUEUED, types.GRACE, types.EXPIRED, types.EXECUTED, types.ABROGATED}
	for _, s := range proposalStates {
		if s == types.ProposalState(state) {
			return true
		}
	}

	return false
}

func getTimeLeft(state types.ProposalState, createTime, warmUpDuration, activeDuration, queueDuration, gracePeriodDuration int64) *int64 {
	now := time.Now().Unix()
	var timeLeft int64

	switch state {
	case types.CANCELED, types.FAILED, types.ACCEPTED, types.EXPIRED, types.EXECUTED, types.ABROGATED:
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
