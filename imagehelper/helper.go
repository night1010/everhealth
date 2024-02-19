package imagehelper

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"github.com/night1010/everhealth/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"google.golang.org/api/option"
)

const (
	GoogleMethod     = "googlestorage"
	CloudinaryMethod = "cloudinary"
)

type ImageHelper interface {
	Upload(ctx context.Context, file io.Reader, folder string, key string) (string, error)
	Destroy(ctx context.Context, folder string, key string) error
}

func NewImageHelper(impl string) (ImageHelper, error) {
	if impl == CloudinaryMethod {
		cldConfig := config.NewCloudinaryConfig()
		cld, err := cloudinary.NewFromParams(cldConfig.Name, cldConfig.Key, cldConfig.Secret)
		return &cloudinaryService{cld: cld}, err
	}
	if impl == GoogleMethod {
		gglConfig := config.NewGoogleCloudStorageConfig()
		storageClient, err := storage.NewClient(context.Background(), option.WithCredentialsFile(gglConfig.JsonCred))
		return &googleStorageService{strg: storageClient, cdnUrl: gglConfig.CdnURl, bucket: gglConfig.Bucket}, err
	}
	return nil, nil
}
