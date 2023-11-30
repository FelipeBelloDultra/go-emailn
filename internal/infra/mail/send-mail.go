package mail

import (
	"os"

	"github.com/FelipeBelloDultra/emailn/internal/domain/campaign"
	"gopkg.in/gomail.v2"
)

func SendMail(campaign *campaign.Campaign) error {
	dialer := gomail.NewDialer(os.Getenv("EMAIL_SMTP"), 1025, os.Getenv("EMAIL_USER"), os.Getenv("EMAIL_PASSWORD"))

	var emails []string
	for _, contact := range campaign.Contacts {
		emails = append(emails, contact.Email)
	}

	message := gomail.NewMessage()

	message.SetHeader("From", "Welcome <felipe_bello_dultra@hotmail.com>")
	message.SetHeader("To", emails...)
	message.SetHeader("Subject", campaign.Name)
	message.SetHeader("text/html", campaign.Content)

	return dialer.DialAndSend(message)
}
