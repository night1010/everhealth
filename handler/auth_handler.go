package handler

import (
	"net/http"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthHandler struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(u usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: u,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var request dto.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	user := request.ToUser()
	err := h.usecase.Register(c.Request.Context(), user)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "link sent"})
}

func (h *AuthHandler) Verify(c *gin.Context) {
	var requestUri dto.VerifyRequestUri
	var request dto.VerifyRequest
	if err := c.ShouldBindUri(&requestUri); err != nil {
		_ = c.Error(err)
		return
	}
	if err := c.ShouldBindWith(&request, binding.Form); err != nil {
		_ = c.Error(err)
		return
	}
	token := c.Query("token")
	if token == "" {
		err := apperror.NewInvalidPathQueryParamError(apperror.NewInvalidTokenError())
		_ = c.Error(err)
		return
	}
	user := request.ToUser(token)
	profile, err := request.ToProfile()
	if err != nil {
		_ = c.Error(err)
		return
	}
	user.RoleId = entity.RoleId(requestUri.Id)
	var doctorProfile *entity.DoctorProfile
	if user.RoleId == entity.RoleDoctor {
		doctorProfile = request.ToDoctorProfile()
	}

	err = h.usecase.Verify(c.Request.Context(), user, profile, doctorProfile)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "verify success"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	user := request.ToUser()
	tokenUser, err := h.usecase.Login(c.Request.Context(), user)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Data: dto.LoginResponse{Token: tokenUser.Token, RoleId: uint(tokenUser.RoleId)}})
}

func (h *AuthHandler) RequestForgotPassword(c *gin.Context) {
	var request dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	user := request.ToUser()
	token := dto.ToForgotPasswordEntity()
	err := h.usecase.ForgotPassword(c.Request.Context(), user, token)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "link sent"})
}

func (h *AuthHandler) ApplyPassword(c *gin.Context) {
	var request dto.ApplyPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	user := request.ToUser()
	token := dto.ToTokenEntity(c.Query("token"))
	err := h.usecase.ResetPassword(c.Request.Context(), user, token)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "password changed"})
}
