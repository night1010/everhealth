package config

import (
	"os"

	"github.com/joho/godotenv"
)

var rajaOngkirConfig *RajaOngkirConfig

type RajaOngkirConfig struct {
	Token string
}

func NewRajaOngkirConfig() *RajaOngkirConfig {
	if rajaOngkirConfig == nil {
		rajaOngkirConfig = initializeRajaOngkirConfig()
	}
	return rajaOngkirConfig
}

func initializeRajaOngkirConfig() *RajaOngkirConfig {
	_ = godotenv.Load()
	token := os.Getenv("RAJAONGKIR_TOKEN")
	return &RajaOngkirConfig{
		Token: token,
	}
}
