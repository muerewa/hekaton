package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

func TemplateHandler(message, result string) (string, error) {
	// Обработка шаблона сообщения
	tmpl, err := template.New("message").Parse(message)
	if err != nil {
		return "", fmt.Errorf("parsing template error: %v", err)
	}

	var res bytes.Buffer
	data := struct {
		Result string
	}{
		Result: result,
	}

	err = tmpl.Execute(&res, data)
	if err != nil {
		return "", fmt.Errorf("parsing template error: %v", err)
	}
	return res.String(), nil
}

func SendTelegramMessage(name string, params map[string]string, result string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", params["token"])
	template, err := TemplateHandler(params["message"], result)
	if err != nil {
		return err
	}
	text := fmt.Sprintf("%s: %s", name, template)

	payload := map[string]string{
		"chat_id": params["chat_id"],
		"text":    text,
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
