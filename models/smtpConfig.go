package models

import "os"

type SmtpConfig struct {

	SmtpAddress		string
	SmtpPort		string
	SmtpEmail		string
	SmtpPass		string
}

func NewSmtpConfig() *SmtpConfig {
	return &SmtpConfig{
		SmtpAddress: os.Getenv("SMTP_SERV"),
		SmtpPort:    os.Getenv("SMTP_PORT"),
		SmtpEmail:   os.Getenv("FROM"),
		SmtpPass:    os.Getenv("PASS"),
	}
}

