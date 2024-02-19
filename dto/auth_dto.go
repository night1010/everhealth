package dto

import (
	"time"

	"github.com/night1010/everhealth/entity"
)

type RegisterRequest struct {
	Email  string        `binding:"required,email" json:"email"`
	RoleId entity.RoleId `binding:"required,oneof=1 2" json:"role_id"`
}

type VerifyRequestUri struct {
	Id uint `uri:"id" binding:"required,numeric"`
}

type VerifyRequest struct {
	Name             string `binding:"required" form:"name"`
	Birthdate        string `binding:"required" form:"dob"`
	Password         string `binding:"required" form:"password"`
	SpecializationId uint   `form:"specialization_id"`
	YearOfExperience uint   `form:"yoe"`
}

type LoginRequest struct {
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required" json:"password"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	RoleId uint   `json:"role_id"`
}

type ForgotPasswordRequest struct {
	Email string `binding:"required,email" json:"email"`
}

type ApplyPasswordRequest struct {
	Password string `binding:"required" json:"password"`
}

func (r *RegisterRequest) ToUser() *entity.User {
	return &entity.User{
		Email:      r.Email,
		IsVerified: false,
		RoleId:     r.RoleId,
	}
}

func (r *VerifyRequest) ToUser(token string) *entity.User {
	return &entity.User{
		Password: r.Password,
		Token:    token,
	}
}

func (r *VerifyRequest) ToProfile() (*entity.Profile, error) {
	parsedBod, err := time.Parse("2006-01-02", r.Birthdate)
	if err != nil {
		return nil, err
	}
	return &entity.Profile{
		Name:      r.Name,
		Birthdate: parsedBod,
	}, nil
}

func (r *VerifyRequest) ToDoctorProfile() *entity.DoctorProfile {
	return &entity.DoctorProfile{
		SpecialistId:     r.SpecializationId,
		YearOfExperience: r.YearOfExperience,
		Status:           entity.Offline,
	}
}

func (r *LoginRequest) ToUser() *entity.User {
	return &entity.User{
		Email:    r.Email,
		Password: r.Password,
	}
}

func (r *ForgotPasswordRequest) ToUser() *entity.User {
	return &entity.User{
		Email: r.Email,
	}
}

func (r *ApplyPasswordRequest) ToUser() *entity.User {
	return &entity.User{
		Password: r.Password,
	}
}

func ToForgotPasswordEntity() *entity.ForgotPasswordToken {
	return &entity.ForgotPasswordToken{
		ExpiredAt: time.Now().Add(3 * time.Minute),
		IsActive:  true,
	}
}

func ToTokenEntity(token string) *entity.ForgotPasswordToken {
	return &entity.ForgotPasswordToken{
		Token: token,
	}
}
