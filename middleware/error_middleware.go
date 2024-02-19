package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/dto"
	"github.com/gin-gonic/gin"
)

func Error() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) < 1 {
			return
		}

		err := c.Errors.Last().Err

		var sErr *json.SyntaxError
		var uErr *json.UnmarshalTypeError
		var cErr *apperror.ClientError

		isClientError := false
		if errors.As(err, &cErr) {
			isClientError = true
			err = cErr.UnWrap()
		}

		message := strings.Split(err.Error(), "\n")
		switch {
		case err.Error() == "invalid request":
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Error: []string{"invalid request"},
			})
		case errors.Is(err, io.EOF):
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Error: message,
			})
		case errors.As(err, &sErr):
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Error: message,
			})
		case errors.As(err, &uErr):
			c.AbortWithStatusJSON(http.StatusBadRequest, dto.Response{
				Error: message,
			})
		case isClientError:
			c.AbortWithStatusJSON(cErr.HttpStatusCode(), dto.Response{
				Error: message,
			})
		case errors.Is(err, context.DeadlineExceeded):
			c.AbortWithStatusJSON(http.StatusGatewayTimeout, dto.Response{
				Error: message,
			})
		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
				Error: message,
			})
		}
	}
}
