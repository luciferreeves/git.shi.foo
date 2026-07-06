package mail

import (
	"fmt"
	"net/smtp"

	"git.shi.foo/config"
	"git.shi.foo/utils/logger"
)

func Send(to string, subject string, body string) error {
	if config.Mail.Host == "" {
		logger.Warnf(LogPrefix, NotConfigured, to)
		return nil
	}

	address := fmt.Sprintf("%s:%d", config.Mail.Host, config.Mail.Port)
	message := buildMessage(to, subject, body)

	var authentication smtp.Auth
	if config.Mail.Username != "" {
		authentication = smtp.PlainAuth("", config.Mail.Username, config.Mail.Password, config.Mail.Host)
	}

	if deliveryError := smtp.SendMail(address, authentication, config.Mail.From, []string{to}, message); deliveryError != nil {
		logger.Errorf(LogPrefix, DeliveryFailed, to, deliveryError)
		return deliveryError
	}

	logger.Infof(LogPrefix, MailSent, to)
	return nil
}

func buildMessage(to string, subject string, body string) []byte {
	headers := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n",
		config.Mail.From, to, subject,
	)
	return []byte(headers + body)
}
