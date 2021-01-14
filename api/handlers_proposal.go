package api

import (
	"database/sql"
	"fmt"
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
			   ( select * from proposal_state(proposal_id) ) as proposal_state
		from governance_proposals
		where proposal_ID = $1
	`, proposalID).Scan(&id, &proposer, &description, &title, &createTime, &targets, &values, &signatures, &calldatas, &timestamp, &warmUpDuration, &activeDuration, &queueDuration, &gracePeriodDuration, &acceptanceThreshold, &minQuorum, &state)

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
			   ( select proposal_state(proposal_id) ) as proposal_state
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
