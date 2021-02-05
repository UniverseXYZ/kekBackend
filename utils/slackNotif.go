package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const webHookUrl = "https://hooks.slack.com/services/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"

type SlackRequestBody struct {
	Text string `json:"text"`
}

func SendSlackNotification(msg string) error {
	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, webHookUrl, bytes.NewBuffer(slackBody))
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
