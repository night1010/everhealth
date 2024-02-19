package config

import (
	"os"

	"github.com/joho/godotenv"
)

var emailConfig *EmailConfig

type EmailConfig struct {
	Name       string
	Address    string
	Password   string
	PrefixLink string
}

func NewEmailConfig() *EmailConfig {
	if emailConfig == nil {
		emailConfig = initializeEmailConfig()
	}
	return emailConfig
}

func initializeEmailConfig() *EmailConfig {
	_ = godotenv.Load()

	name := os.Getenv("EMAIL_SENDER_NAME")
	address := os.Getenv("EMAIL_SENDER_ADDRESS")
	password := os.Getenv("EMAIL_SENDER_PASSWORD")
	prefixLink := os.Getenv("PREFIX_LINK")

	return &EmailConfig{
		Name:       name,
		Address:    address,
		Password:   password,
		PrefixLink: prefixLink,
	}

}
