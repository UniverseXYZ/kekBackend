package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type Notifier struct {
	config Config
}

var instance *Notifier

func Init(config Config) {
	if instance != nil {
		return
	}

	instance = &Notifier{config: config}
}

type SlackRequestBody struct {
	Text string `json:"text"`
}

func SendNotification(msg string) error {
	if instance == nil || instance.config.Webhook == "" {
		return nil
	}

	slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, instance.config.Webhook, bytes.NewBuffer(slackBody))
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
		return errors.New("non-ok response returned from Slack")
	}

	return nil
}
