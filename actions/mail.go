package actions

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/jordan-wright/email"
)

func SendEmail(params map[string]string, result string) error {
	// Parse params
	host := params["smtp_host"]
	port := params["smtp_port"]
	username := params["username"]
	password := params["password"]
	from := params["from"]
	to := strings.Split(params["to"], ",") // Разделение получателей
	subject := params["subject"]
	bodyTemplate := params["body"]

	// Insert result in template
	body := strings.ReplaceAll(bodyTemplate, "{{.Result}}", result)

	// Creating email
	e := email.NewEmail()
	e.From = from
	e.To = to
	e.Subject = subject
	e.Text = []byte(body)

	// Auth setting
	auth := smtp.PlainAuth("", username, password, host)

	// Mail sending
	err := e.Send(
		fmt.Sprintf("%s:%s", host, port),
		auth,
	)

	return err
}
