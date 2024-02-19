package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const defaultGracefulShutdownTimeout = 5
const defaultRequestTimeout = 5

var appConfig *AppConfig

type AppConfig struct {
	Env                     string
	Name                    string
	GracefulShutdownTimeout time.Duration
	RequestTimeout          time.Duration
}

func (c *AppConfig) IsInDevMode() bool {
	return c.Env == "dev"
}

func (c *AppConfig) IsInDebugMode() bool {
	return c.Env == "debug"
}

func NewAppConfig() *AppConfig {
	if appConfig == nil {
		appConfig = initializeAppConfig()
	}
	return appConfig
}

func initializeAppConfig() *AppConfig {
	_ = godotenv.Load()

	env := os.Getenv("APP_ENV")
	name := os.Getenv("APP_NAME")

	gracefulShutdownTimeoutString := os.Getenv("GRACEFUL_SHUTDOWN_TIMEOUT")
	gracefulShutdownTimeout, err := strconv.Atoi(gracefulShutdownTimeoutString)
	if err != nil {
		gracefulShutdownTimeout = defaultGracefulShutdownTimeout
	}

	requestTimeoutString := os.Getenv("REQUEST_TIMEOUT")
	requestTimeout, err := strconv.Atoi(requestTimeoutString)
	if err != nil {
		requestTimeout = defaultRequestTimeout
	}

	return &AppConfig{
		Env:                     env,
		Name:                    name,
		GracefulShutdownTimeout: time.Duration(gracefulShutdownTimeout),
		RequestTimeout:          time.Duration(requestTimeout),
	}
}
