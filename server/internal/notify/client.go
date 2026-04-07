package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const expoURL = "https://exp.host/--/api/v2/push/send"

type message struct {
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Sound string `json:"sound"`
}

func Send(pushToken, title, body string) error {
	msg := message{
		To:    pushToken,
		Title: title,
		Body:  body,
		Sound: "default",
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(expoURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("send push notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expo push returned status %d", resp.StatusCode)
	}
	return nil
}
