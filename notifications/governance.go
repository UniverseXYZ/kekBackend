package notifications

type ProposalData struct {
	Id int64
}

func ProposalCreatedJob(id int64) (*Job, error) {
	return NewJob("proposal-created", 0, ProposalData{id})
}

// 		_, err = stmt.Exec(p.Id.Int64(), p.Proposer.String(), p.Title, p.CreateTime.Int64(), p.WarmUpDuration.Int64(), p.ActiveDuration.Int64(), p.QueueDuration.Int64(), p.GracePeriodDuration.Int64(), g.Preprocessed.BlockTimestamp)
// func FromGovernanceProposal(id int64, proposer string, title string, createTime int64, warmUpDuration int64, activeDuration int64, queueDuration int64, graceDuration int64, blockNumber int64, blockTime int64) []Notification {
// 	// TODO starts at blockTime -1 or creation time?
// 	startTime := blockTime - 1
//
// 	createNotif := NewNotification(
// 		"system",
// 		"proposal-created",
// 		blockNumber,
// 		startTime,
// 		startTime+warmUpDuration-300,
// 		fmt.Sprintf("Proposal PID-%d created by %s", id, proposer),
// 		nil,
// 	)
// 	activatingNotif := NewNotification(
// 		"system",
// 		"proposal-activating",
// 		blockNumber,
// 		startTime+warmUpDuration-300,
// 		startTime+warmUpDuration,
// 		fmt.Sprintf(fmt.Sprintf("Voting period for PID-%d starting in 5 minutes"), id),
// 		nil,
// 	)
// 	activeNotif := NewNotification(
// 		"system",
// 		"proposal-active",
// 		blockNumber,
// 		startTime+warmUpDuration,
// 		startTime+warmUpDuration+activeDuration-300,
// 		fmt.Sprintf(fmt.Sprintf("Governace proposal PID-%d is now active"), id),
// 		nil,
// 	)
//
// 	return []Notification{
// 		createNotif,
// 		activatingNotif,
// 		activeNotif,
// 	}
// }
