package misc

import (
	"embed"
	"time"

	"github.com/go-mail/mail/v2"
)

var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {

	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m Mailer) Send(recipient string, url string) error {

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", "test")
	msg.SetBody("text/plain", "use this link to reset your password"+url)
	//msg.AddAlternative("text/html", htmlBody.String())

	err := m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}
