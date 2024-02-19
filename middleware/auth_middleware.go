package middleware

import (
	"context"
	"strings"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/appjwt"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/util"
	"github.com/gin-gonic/gin"
)

func Auth(roles ...entity.RoleId) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		token, err := extractBearerToken(bearerToken)
		if err != nil {
			c.Abort()
			_ = c.Error(err)
			return
		}

		newJwt := appjwt.NewJwt()
		claims, err := newJwt.ValidateToken(token)
		if err != nil {
			c.Abort()
			_ = c.Error(apperror.NewInvalidTokenError())
			return
		}

		ctx := context.WithValue(c.Request.Context(), "user_id", claims.Id)
		ctx = context.WithValue(ctx, "role_id", claims.RoleId)
		c.Request = c.Request.WithContext(ctx)

		if !util.IsMemberOf(roles, claims.RoleId) {
			c.Abort()
			_ = c.Error(apperror.NewForbiddenActionError("permission denied"))
			return
		}

		c.Next()
	}
}

func extractBearerToken(bearerToken string) (string, error) {
	if bearerToken == "" {
		return "", apperror.NewMissingTokenError()
	}
	token := strings.Split(bearerToken, " ")
	if len(token) != 2 {
		return "", apperror.NewInvalidTokenError()
	}
	return token[1], nil
}
