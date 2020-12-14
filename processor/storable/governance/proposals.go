package governance

import (
	"database/sql"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/contracts"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (g GovStorable) handleProposals(logs []web3types.Log, tx *sql.Tx) error {
	var proposals []Proposal
	var actions []Action

	ctr, err := contracts.NewGovernance(common.HexToAddress(g.config.GovernanceAddress), &g.GovernanceClient)
	if err != nil {
		return err
	}

	for _, log := range logs {
		if utils.CleanUpHex(log.Topics[0]) == utils.CleanUpHex(g.govAbi.Events["ProposalCreated"].ID.String()) {
			proposalID, err := utils.HexStrToBigInt(log.Topics[1])
			if err != nil {
				return err
			}

			p, err := ctr.Proposals(nil, proposalID)
			if err != nil {
				return errors.Wrap(err, "could not get the proposals from contract")
			}

			proposals = append(proposals, p)

			a, err := ctr.GetActions(nil, proposalID)
			if err != nil {
				return errors.Wrap(err, "could not get the actions from contract")
			}

			actions = append(actions, a)
		}
	}

	//db

	return nil
}
