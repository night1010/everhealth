package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const defaultJwtExpiryDuration = 1

var jwtConfig *JwtConfig

type JwtConfig struct {
	ExpiryDuration time.Duration
	Secret         string
}

func NewJwtConfig() *JwtConfig {
	if jwtConfig == nil {
		jwtConfig = initializeJwtConfig()
	}
	return jwtConfig
}

func initializeJwtConfig() *JwtConfig {
	_ = godotenv.Load()

	jwtExpiryDurationString := os.Getenv("JWT_EXPIRY_DURATION")
	jwtExpiryDuration, err := strconv.Atoi(jwtExpiryDurationString)
	if err != nil {
		jwtExpiryDuration = defaultJwtExpiryDuration
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	return &JwtConfig{
		ExpiryDuration: time.Duration(jwtExpiryDuration),
		Secret:         jwtSecret,
	}

}
