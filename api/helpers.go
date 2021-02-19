package api

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/barnbridge/barnbridge-backend/api/types"
	"github.com/pkg/errors"
)

func calculateOffset(limit string, page string) (string, error) {
	l, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return "", err
	}

	p, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		return "", err
	}

	offset := (p - 1) * l

	return strconv.FormatInt(offset, 10), nil
}

func (a *API) getProposalEvents(id uint64) ([]types.Event, error) {
	rows, err := a.db.Query(`
		select proposal_id,
		       caller,
		       event_type,
		       event_data,
		       block_timestamp,
		       tx_hash
		from governance_events 
		where proposal_id = $1`, id)
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "could not query proposal events")
	}

	var eventsList []types.Event

	for rows.Next() {
		var event types.Event
		err := rows.Scan(&event.ProposalID, &event.Caller, &event.EventType, &event.Eta, &event.CreateTime, &event.TxHash)
		if err != nil {
			return nil, errors.Wrap(err, "could not scan proposal event")
		}

		eventsList = append(eventsList, event)
	}

	return eventsList, nil
}

func (a *API) getHighestBlock() (*int64, error) {
	var number int64

	err := a.db.QueryRow(`select number from blocks order by number desc limit 1;`).Scan(&number)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "could not get highest block")
	}

	return &number, nil
}

func (a *API) getERC721Transfer(tokenType string, userAddress string, offset string, limit string) ([]types.ERC721Transfer, error) {
	query := `
			select  token_address,
					token_type,
					sender,
					receiver,
			       token_id,
					tx_hash,
					tx_index,
					log_index,
					block_timestamp,
					included_in_block
		from erc721_transfers
		where 1=1 %s %s order by block_timestamp desc
		offset $1 
		limit $2 `

	var parameters = []interface{}{offset, limit}

	if userAddress != "" {
		parameters = append(parameters, userAddress)
		userFilter := fmt.Sprintf("and owner_address = $%d", len(parameters))
		query = fmt.Sprintf(query, userFilter)
	} else {
		query = fmt.Sprintf(query, "")
	}

	if tokenType != "all" {
		parameters = append(parameters, tokenType)
		tokenTypeFilter := fmt.Sprintf("and token_type = $%d", len(parameters))
		query = fmt.Sprintf(query, tokenTypeFilter)
	} else {
		query = fmt.Sprintf(query, "")
	}

	rows, err := a.db.Query(query, parameters...)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}

	var transfers []types.ERC721Transfer

	for rows.Next() {
		var t types.ERC721Transfer
		err := rows.Scan(&t.TokenAddress, &t.TokenType, &t.Sender, &t.Receiver, &t.TokenID, &t.TransactionHash, &t.TransactionIndex, &t.LogIndex, &t.BlockTimestamp, &t.BlockNumber)
		if err != nil {
			return nil, err
		}

		transfers = append(transfers, t)
	}

	return transfers, nil
}
