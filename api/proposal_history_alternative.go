package api

// func (a *API) generateFullHistory(p types.Proposal) ([]types.HistoryEvent, error) {
// 	events, err := a.getProposalEvents(p.Id)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "could not get proposal events")
// 	}
//
// 	var history []types.HistoryEvent
//
// 	for _, e := range events {
// 		h := types.HistoryEvent{
// 			Name:    e.EventType,
// 			StartTs: e.CreateTime,
// 			EndTs:   0,
// 		}
//
// 		if h.Name == string(types.CREATED) {
// 			h.StartTs = h.StartTs - 1
// 		}
//
// 		history = append(history, h)
// 	}
//
// 	history = append(history, types.HistoryEvent{
// 		Name:    string(types.WARMUP),
// 		StartTs: p.CreateTime,
// 	})
//
// 	history = append(history, types.HistoryEvent{
// 		Name:    string(types.ACTIVE),
// 		StartTs: p.CreateTime + p.WarmUpDuration + 1,
// 	})
//
// 	history = append(history, types.HistoryEvent{
// 		Name:    string(types.ACCEPTED),
// 		StartTs: p.CreateTime + p.WarmUpDuration + p.ActiveDuration + 1,
// 	})
//
// 	history = append(history, types.HistoryEvent{
// 		Name:    string(types.FAILED),
// 		StartTs: p.CreateTime + p.WarmUpDuration + p.ActiveDuration + 1,
// 	})
//
// 	history = append(history, types.HistoryEvent{
// 		Name:    string(types.GRACE),
// 		StartTs: p.CreateTime + p.WarmUpDuration + p.ActiveDuration + p.QueueDuration,
// 	})
//
// 	history = append(history, types.HistoryEvent{
// 		Name:    string(types.EXPIRED),
// 		StartTs: p.CreateTime + p.WarmUpDuration + p.ActiveDuration + p.QueueDuration + p.GracePeriodDuration + 1,
// 	})
//
// 	return history, nil
// }
//
// func (a *API) filterHistory(p types.Proposal, history []types.HistoryEvent) []types.HistoryEvent {
// 	// now ... we have to eliminate all events where starTs > time.Now()
// 	var history2 []types.HistoryEvent
// 	for _, e := range history {
// 		if e.StartTs <= time.Now().Unix() {
// 			history2 = append(history2, e)
// 		}
// 	}
//
// 	var history3 []types.HistoryEvent
//
// 	// second step: remove all events that follow a final event (Canceled, Executed, Expired)
// 	canceled := findEvent(history2, types.CANCELED)
// 	if canceled != nil {
// 		for _, e := range history2 {
// 			if e.StartTs <= canceled.StartTs {
// 				history3 = append(history3, e)
// 			}
// 		}
// 	} else {
// 		history3 = history2
// 	}
//
// 	var history4 []types.HistoryEvent
// 	if findEvent(history3, types.QUEUED) != nil {
// 		for _, e := range history3 {
// 			if e.Name != string(types.FAILED) {
// 				history4 = append(history4, e)
// 			}
// 		}
// 	} else {
// 		if p.State == types.FAILED {
// 			failed := findEvent(history3, types.FAILED)
// 			for _, e := range history3 {
// 				if e.StartTs <= failed.StartTs && e.Name != string(types.ACCEPTED) {
// 					history4 = append(history4, e)
// 				}
// 			}
// 		} else {
// 			history4 = history3
// 		}
// 	}
//
// 	var history5 []types.HistoryEvent
// 	executed := findEvent(history4, types.EXECUTED)
// 	if executed != nil {
// 		for _, e := range history4 {
// 			history5 = append(history5, e)
// 		}
// 	} else {
// 		history5 = history4
// 	}
//
// 	return history5
// }
//
// func (a *API) buildHistory(p types.Proposal) ([]types.HistoryEvent, error) {
// 	fullHistory, err := a.generateFullHistory(p)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	history := a.filterHistory(p, fullHistory)
//
// 	sort.Slice(history, func(i, j int) bool {
// 		if history[i].Name == string(types.CREATED) && history[j].Name == string(types.WARMUP) {
// 			return false
// 		} else if history[j].Name == string(types.CREATED) && history[i].Name == string(types.WARMUP) {
// 			return true
// 		}
//
// 		return history[i].StartTs > history[j].StartTs
// 	})
//
// 	for i := 1; i <= len(history)-1; i++ {
// 		history[i].EndTs = history[i-1].StartTs - 1
// 	}
//
// 	return history, nil
// }
//
// func findEvent(history []types.HistoryEvent, name types.ProposalState) *types.HistoryEvent {
// 	for _, e := range history {
// 		if e.Name == string(name) {
// 			return &e
// 		}
// 	}
//
// 	return nil
// }
