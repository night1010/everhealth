package imagehelper

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

type googleStorageService struct {
	strg   *storage.Client
	cdnUrl string
	bucket string
}

func (g *googleStorageService) Upload(ctx context.Context, file io.Reader, folder string, key string) (string, error) {
	object := folder + "/" + key
	sw := g.strg.Bucket(g.bucket).Object(object).NewWriter(context.Background())
	sw.CacheControl = "no-cache"
	if _, err := io.Copy(sw, file); err != nil {
		return "", err
	}
	if err := sw.Close(); err != nil {
		return "", err
	}
	url := fmt.Sprintf(g.cdnUrl+"%s", object)
	return url, nil
}

func (g *googleStorageService) Destroy(ctx context.Context, folder string, key string) error {
	object := folder + "/" + key
	if err := g.strg.Bucket(g.bucket).Object(object).Delete(context.Background()); err != nil {
		return err
	}
	return nil
}
