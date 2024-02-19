package config

import (
	"os"

	"github.com/joho/godotenv"
)

var cloudinaryConfig *CloudinaryConfig

type CloudinaryConfig struct {
	Secret string
	Name   string
	Key    string
}

func NewCloudinaryConfig() *CloudinaryConfig {
	if cloudinaryConfig == nil {
		cloudinaryConfig = initializeCloudinaryConfig()
	}
	return cloudinaryConfig
}

func initializeCloudinaryConfig() *CloudinaryConfig {
	_ = godotenv.Load()

	cldName := os.Getenv("CLD_NAME")
	cldKey := os.Getenv("CLD_KEY")
	cldSecret := os.Getenv("CLD_SECRET")

	return &CloudinaryConfig{
		Name:   cldName,
		Key:    cldKey,
		Secret: cldSecret,
	}

}
