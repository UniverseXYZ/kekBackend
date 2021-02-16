package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/barnbridge/barnbridge-backend/types"
)

type SlackRequestBody struct {
	Text string `json:"text"`
}

func SendSlackNotification(msg string, slackNotif types.SlackNotif) error {
	if !slackNotif.Enabled {
		return nil
	}

	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, slackNotif.Webhook, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}
