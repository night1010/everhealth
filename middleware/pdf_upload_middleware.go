package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/night1010/everhealth/apperror"
	"github.com/gin-gonic/gin"
)

func PDFUpload() gin.HandlerFunc {
	return func(c *gin.Context) {
		const (
			MB = 1 << 20
		)

		type Sizer interface {
			Size() int64
		}

		file, _, err := c.Request.FormFile("pdf")
		if err != nil {
			if err == http.ErrMissingFile || c.Request.Method == http.MethodPut {
				c.Next()
				return
			}
			c.Error(apperror.NewClientError(err))
			c.Abort()
			return
		}
		defer file.Close()
		if err := c.Request.ParseMultipartForm(1 * MB); err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 1*MB) // 1 Mb for PDF

		fileHeader := make([]byte, 512)

		if _, err := file.Read(fileHeader); err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if _, err := file.Seek(0, 0); err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		if file.(Sizer).Size() > 1000000 {
			c.Error(apperror.NewClientError(errors.New("PDF must be below 1 MB")))
			c.Abort()
			return
		}
		if http.DetectContentType(fileHeader) != "application/pdf" {
			c.Error(apperror.NewClientError(errors.New("File type must be PDF")))
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), "pdf", file)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
