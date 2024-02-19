package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/appjwt"
	"github.com/night1010/everhealth/chat"
	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var WebsocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatHandler struct {
	manager     *chat.Manager
	chatUsecase usecase.ChatUsecase
	jwt         appjwt.Jwt
}

func NewChatHandler(manager *chat.Manager, chatUsecase usecase.ChatUsecase, jwt appjwt.Jwt) *ChatHandler {
	return &ChatHandler{
		manager:     manager,
		chatUsecase: chatUsecase,
		jwt:         jwt,
	}
}

func (h *ChatHandler) Handle(c *gin.Context) {
	var uri dto.RequestUri
	if err := c.ShouldBindUri(&uri); err != nil {
		_ = c.Error(err)
		return
	}

	var token dto.ChatToken
	if err := c.ShouldBindQuery(&token); err != nil {
		_ = c.Error(err)
		return
	}

	user, err := h.jwt.ValidateToken(token.Token)
	if err != nil {
		_ = c.Error(err)
		return
	}

	ctx := context.WithValue(c.Request.Context(), "user_id", user.Id)
	ctx = context.WithValue(ctx, "role_id", user.RoleId)
	c.Request = c.Request.WithContext(ctx)

	isValid, err := h.chatUsecase.ValidateTelemedicine(c.Request.Context(), uint(uri.Id))
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !isValid {
		_ = c.Error(apperror.NewForbiddenActionError("can't join invalid chat room"))
		return
	}

	conn, err := WebsocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		_ = c.Error(fmt.Errorf("failed to upgrade connection to ws: %w", err))
		return
	}

	userId := c.Request.Context().Value("user_id").(uint)

	client := chat.NewClient(conn, h.manager, uint(uri.Id), userId)
	h.manager.AddClient(client)

	go client.ReadMessages()
	go client.WriteMessages()
}
