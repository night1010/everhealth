package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHAndler(u usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: u,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	user, profile, doctorProfile, err := h.usecase.UserProfile(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToUserProfileDTO(user, profile, doctorProfile))
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var request dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	err := h.usecase.ResetPassword(c.Request.Context(), request.OldPassword, request.NewPassword)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ResetPasswordResponse{Message: "password changed"})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var request dto.UpdateProfileRequest
	var doctorProfile *entity.DoctorProfile
	err := c.ShouldBindWith(&request, binding.Form)
	if err != nil {
		_ = c.Error(err)
		return
	}
	roleId := c.Request.Context().Value("role_id").(entity.RoleId)
	if roleId == entity.RoleDoctor {
		if request.Fee == "" || request.YearOfExperience == 0 {
			_ = c.Error(apperror.NewClientError(fmt.Errorf("doctor should include doctor data")))
			return
		}
		doctorProfile, err = request.ToDoctorProfile()
		if err != nil {
			_ = c.Error(err)
			return
		}
	}
	profile := request.ToProfile()
	if err != nil {
		_ = c.Error(err)
		return
	}
	err = h.usecase.UpdateProfile(c.Request.Context(), profile, doctorProfile)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "update success"})
}

func (h *UserHandler) ListDoctor(c *gin.Context) {
	var request dto.GetDoctorsParam
	if err := c.ShouldBindQuery(&request); err != nil {
		_ = c.Error(err)
		return
	}
	query, err := request.ToQuery()
	if err != nil {
		_ = c.Error(err)
		return
	}
	pagedResult, err := h.usecase.ListAllDoctors(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	var listDoctors []*dto.GetDoctorsResponse
	data := pagedResult.Data.([]*entity.DoctorProfile)
	for _, doctor := range data {
		listDoctors = append(listDoctors, &dto.GetDoctorsResponse{
			Id:               doctor.ProfileId,
			Name:             doctor.Profile.Name,
			Image:            doctor.Profile.Image,
			Specialization:   doctor.Specialist.Name,
			SpecializationId: strconv.Itoa(int(doctor.SpecialistId)),
			Fee:              doctor.Fee.String(),
			YearOfExperience: doctor.YearOfExperience,
			Certificate:      doctor.Certificate,
			Status:           entity.DoctorStatusMap[doctor.Status],
		})
	}
	c.JSON(200, dto.Response{
		Data:        listDoctors,
		CurrentPage: &pagedResult.CurrentPage,
		CurrentItem: &pagedResult.CurrentItems,
		TotalPage:   &pagedResult.TotalPage,
		TotalItem:   &pagedResult.TotalItem,
	})
}

func (h *UserHandler) DoctorDetail(c *gin.Context) {
	var doctorUri dto.DoctorUri
	if err := c.ShouldBindUri(&doctorUri); err != nil {
		_ = c.Error(err)
		return
	}
	profile, doctorProfile, err := h.usecase.DoctorDetail(c.Request.Context(), doctorUri.Id)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.ToDoctorDetailDTO(profile, doctorProfile))
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	var request dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		_ = c.Error(err)
		return
	}
	output := request.ToStatus()
	err := h.usecase.UpdateStatus(c.Request.Context(), output)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dto.Response{Message: "status updated"})
}

func (h *UserHandler) GetAllUser(c *gin.Context) {
	var requestParam dto.UserQueryParamReq
	if err := c.ShouldBindQuery(&requestParam); err != nil {
		_ = c.Error(err)
		return
	}
	query := requestParam.ToQuery()
	pageResult, err := h.usecase.GetAllUser(c.Request.Context(), query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	users := pageResult.Data.([]*entity.User)
	usersRes := []*dto.UserRes{}
	for _, user := range users {
		userRes := dto.NewUserRes(user)
		usersRes = append(usersRes, userRes)
	}
	c.JSON(http.StatusOK, dto.Response{Data: usersRes,
		TotalPage: &pageResult.TotalPage, TotalItem: &pageResult.TotalItem, CurrentPage: &pageResult.CurrentPage, CurrentItem: &pageResult.CurrentItems})
}
