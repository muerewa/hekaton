package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func SendTelegramMessage(params map[string]string, result string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", params["token"])
	payload := map[string]string{
		"chat_id": params["chat_id"],
		"text":    params["message"],
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status)
	}
	return nil
}
