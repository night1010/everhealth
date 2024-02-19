package config

import (
	"os"

	"github.com/joho/godotenv"
)

var frontEnd *FrontendConfig

type FrontendConfig struct {
	Host string
}

func NewFrontendConfig() *FrontendConfig {
	if frontEnd == nil {
		frontEnd = initializeFrontendConfig()
	}
	return frontEnd
}

func initializeFrontendConfig() *FrontendConfig {
	_ = godotenv.Load()

	host := os.Getenv("FRONTEND_HOST")

	return &FrontendConfig{
		Host: host,
	}
}
