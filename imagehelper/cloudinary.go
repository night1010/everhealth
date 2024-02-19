package imagehelper

import (
	"context"
	"io"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type cloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func (c *cloudinaryService) Upload(ctx context.Context, file io.Reader, folder string, key string) (string, error) {
	uploadParams := uploader.UploadParams{
		PublicID: key,
		Folder:   folder,
	}
	resp, err := c.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", err
	}
	return resp.SecureURL, nil
}

func (c *cloudinaryService) Destroy(ctx context.Context, folder string, key string) error {
	publicId := folder + "/" + key
	destroyParams := uploader.DestroyParams{
		PublicID: publicId,
	}
	_, err := c.cld.Upload.Destroy(context.Background(), destroyParams)
	return err
}
