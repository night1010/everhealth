package config

import (
	"os"

	"github.com/joho/godotenv"
)

var restConfig *RestConfig

type RestConfig struct {
	Host string
	Port string
}

func NewRestConfig() *RestConfig {
	if restConfig == nil {
		restConfig = initializeRestConfig()
	}
	return restConfig
}

func initializeRestConfig() *RestConfig {
	_ = godotenv.Load()

	host := os.Getenv("REST_HOST")
	port := os.Getenv("REST_PORT")
	
	return &RestConfig{
		Host: host,
		Port: port,
	}
}
