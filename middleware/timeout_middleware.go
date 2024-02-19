package middleware

import (
	"context"
	"time"

	"github.com/night1010/everhealth/config"
	"github.com/gin-gonic/gin"
)

func Timeout() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		appConfig := config.NewAppConfig()

		ctx, cancel := context.WithTimeout(ctx, appConfig.RequestTimeout*time.Second)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
