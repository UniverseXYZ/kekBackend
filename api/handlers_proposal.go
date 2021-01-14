package api

import (
	"database/sql"
	"math"
	"strconv"
	"strings"

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
		state               string
	)
	err := a.core.DB().QueryRow(`select proposal_id,proposer,description,title,create_time,targets,"values",signatures,calldatas,block_timestamp,warm_up_duration, active_duration ,
       queue_duration ,  grace_period_duration, acceptance_threshold, min_quorum 
       from governance_proposals where proposal_ID = $1 
       and (select * from proposal_state($1) as proposal_state)`, proposalID).Scan(&id, &proposer, &description, &title, &createTime, &targets, &values, &signatures,
		&calldatas, &timestamp, &warmUpDuration, &activeDuration, &queueDuration, &gracePeriodDuration, &acceptanceThreshold, &minQuorum, &state)

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
	}

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

	if proposalState == "all" {
		if title == "" {
			rows, err = a.core.DB().Query(`
			select proposal_ID, proposer, description, title, create_time, targets, "values", signatures, calldatas, block_timestamp,warm_up_duration, active_duration ,
       			queue_duration ,  grace_period_duration, acceptance_threshold, min_quorum
			from governance_proposals 
			where proposal_id <= $1 
			order by proposal_id desc 
			limit $2
			and (select * from proposal_state(proposal_id) as proposal_state)
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
		and (select * from proposal_state(proposal_id) as proposal_state)
		`, offset, title, limit)
		}
	} else {
		if title == "" {
			rows, err = a.core.DB().Query(`
			select proposal_ID, proposer, description, title, create_time, targets, "values", signatures, calldatas, block_timestamp,warm_up_duration, active_duration ,
       			queue_duration ,  grace_period_duration, acceptance_threshold, min_quorum
			from governance_proposals 
			where proposal_id <= $1 
			and proposal_state(proposal_id) = $3
			order by proposal_id desc 
			limit $2
			and (select * from proposal_state(proposal_id) as proposal_state)
		`, offset, limit, proposalState)
		} else {
			title = "%" + strings.ToLower(title) + "%"
			rows, err = a.core.DB().Query(`
			select proposal_ID, proposer, description, title, create_time, targets, "values", signatures, calldatas,block_timestamp,warm_up_duration, active_duration ,
       queue_duration ,  grace_period_duration, acceptance_threshold, min_quorum 
			from governance_proposals 
			where proposal_id <= $1 
			  and lower(title) like $2 
					and proposal_state(proposal_id) = $4
			order by proposal_id desc 
			limit $3
		and (select * from proposal_state(proposal_id) as proposal_state)
		`, offset, title, limit, proposalState)
		}
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
			state               string
		)
		err := rows.Scan(&id, &proposer, &description, &title, &createTime, &targets, &values, &signatures, &calldatas, &timestamp, &warmUpDuration, &activeDuration, &queueDuration, &gracePeriodDuration, &acceptanceThreshold, &minQuorum, &state)
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
			State:               state,
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

func checkStateExist(state string) bool {
	proposalStates := [9]string{"WarmUp", "Active", "Canceled", "Failed", "Accepted", "Queued", "Grace", "Expired", "Executed"}
	for _, s := range proposalStates {
		if strings.ToLower(s) == strings.ToLower(state) {
			return true
		}

	}
	return false
}
