package config

import (
	"os"

	"github.com/joho/godotenv"
)

var googleCloudStorageConfig *GoogleCloudStorageConfig

type GoogleCloudStorageConfig struct {
	CdnURl   string
	JsonCred string
	Bucket   string
}

func NewGoogleCloudStorageConfig() *GoogleCloudStorageConfig {
	if googleCloudStorageConfig == nil {
		googleCloudStorageConfig = initializeGoogleCloudStorageConfig()
	}
	return googleCloudStorageConfig
}

func initializeGoogleCloudStorageConfig() *GoogleCloudStorageConfig {
	_ = godotenv.Load()

	cdnUrl := os.Getenv("GOOGLE_CLOUD_URL_CDN")
	bucket := os.Getenv("GOOGLE_CLOUD_BUCKET")
	return &GoogleCloudStorageConfig{
		CdnURl:   cdnUrl,
		JsonCred: "keys.json",
		Bucket:   bucket,
	}

}
