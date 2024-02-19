package dto

import (
	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"github.com/shopspring/decimal"
)

type Specialization struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type UserProfileResponse struct {
	Name             string         `json:"name"`
	Image            string         `json:"image"`
	Email            string         `json:"email"`
	Birthdate        string         `json:"dob"`
	YearOfExperience uint           `json:"yoe,omitempty"`
	Fee              string         `json:"fee,omitempty"`
	Specialization   Specialization `json:"specialization,omitempty"`
	Certificate      string         `json:"certificate,omitempty"`
	Status           string         `json:"status"`
}

type ResetPasswordRequest struct {
	OldPassword string `binding:"required" json:"old_password"`
	NewPassword string `binding:"required" json:"new_password"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type UpdateProfileRequest struct {
	Name             string `binding:"required" form:"name"`
	YearOfExperience uint   `form:"yoe"`
	Fee              string `form:"fee"`
}

type GetDoctorsParam struct {
	Name           *string `form:"name"`
	Specialization *int    `form:"specialization" binding:"omitempty,numeric,min=1"`
	SortBy         *string `form:"sort_by" binding:"omitempty,oneof=name fee"`
	Order          *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit          *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page           *int    `form:"page" binding:"omitempty,numeric,min=1"`
}

type GetDoctorsResponse struct {
	Id               uint   `json:"id"`
	Name             string `json:"name"`
	Image            string `json:"image"`
	Specialization   string `json:"specialization"`
	SpecializationId string `json:"specialization_id"`
	Fee              string `json:"fee"`
	YearOfExperience uint   `json:"yoe"`
	Certificate      string `json:"certificate"`
	Status           string `json:"status"`
}

type DoctorUri struct {
	Id uint `uri:"id" binding:"required,numeric"`
}

type DoctorDetailResponse struct {
	Name             string `json:"name"`
	Image            string `json:"image"`
	Specialization   string `json:"specialization"`
	Fee              string `json:"fee"`
	YearOfExperience uint   `json:"yoe"`
	Certificate      string `json:"certificate"`
	Status           string `json:"status"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=online offline"`
}

func (r *UpdateStatusRequest) ToStatus() entity.StatusDoctor {
	if r.Status == "offline" {
		return entity.Offline
	} else if r.Status == "online" {
		return entity.Online
	}
	return 0
}

func (dp *GetDoctorsParam) ToQuery() (*valueobject.Query, error) {
	query := valueobject.NewQuery()

	if dp.Page != nil {
		query.WithPage(*dp.Page)
	}
	if dp.Limit != nil {
		query.WithLimit(*dp.Limit)
	}

	if dp.Order != nil {
		query.WithOrder(valueobject.Order(*dp.Order))
	}

	if dp.SortBy != nil {
		query.WithSortBy(*dp.SortBy)
	}

	if dp.Name != nil {
		query.Condition("name", valueobject.ILike, *dp.Name)
	}

	if dp.Specialization != nil {
		query.Condition("specialization", valueobject.Equal, *dp.Specialization)
	}

	return query, nil
}

func ToUserProfileDTO(user *entity.User, profile *entity.Profile, doctorProfile *entity.DoctorProfile) UserProfileResponse {
	var userProfileResponse UserProfileResponse
	userProfileResponse.Name = profile.Name
	userProfileResponse.Image = profile.Image
	userProfileResponse.Email = user.Email
	userProfileResponse.Birthdate = profile.Birthdate.Format("2006-01-02")
	if doctorProfile != nil {
		userProfileResponse.YearOfExperience = doctorProfile.YearOfExperience
		userProfileResponse.Fee = doctorProfile.Fee.String()
		userProfileResponse.Specialization.Id = doctorProfile.SpecialistId
		userProfileResponse.Specialization.Name = doctorProfile.Specialist.Name
		userProfileResponse.Status = entity.DoctorStatusMap[doctorProfile.Status]
		userProfileResponse.Certificate = string(doctorProfile.Certificate)
	}
	return userProfileResponse
}

func (r *UpdateProfileRequest) ToProfile() *entity.Profile {
	return &entity.Profile{
		Name: r.Name,
	}
}

func (r *UpdateProfileRequest) ToDoctorProfile() (*entity.DoctorProfile, error) {
	fee, err := decimal.NewFromString(r.Fee)
	if err != nil {
		return nil, apperror.NewClientError(apperror.NewResourceStateError("invalid fee"))
	}
	return &entity.DoctorProfile{
		YearOfExperience: r.YearOfExperience,
		Fee:              fee,
	}, nil
}

func ToDoctorDetailDTO(profile *entity.Profile, doctorProfile *entity.DoctorProfile) DoctorDetailResponse {
	var doctorDetailResponse DoctorDetailResponse
	doctorDetailResponse.Name = profile.Name
	doctorDetailResponse.Image = profile.Image
	doctorDetailResponse.Specialization = doctorProfile.Specialist.Name
	doctorDetailResponse.Fee = doctorProfile.Fee.String()
	doctorDetailResponse.YearOfExperience = doctorProfile.YearOfExperience
	doctorDetailResponse.Certificate = doctorProfile.Certificate
	doctorDetailResponse.Status = entity.DoctorStatusMap[doctorProfile.Status]
	return doctorDetailResponse
}

type UserQueryParamReq struct {
	Email      *string        `form:"email"`
	IsVerified *bool          `form:"is_verified"`
	RoleId     *entity.RoleId `form:"role_id" binding:"omitempty,oneof=1 2 3"`
	SortBy     *string        `form:"sort_by" binding:"omitempty,oneof=email"`
	Order      *string        `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit      *int           `form:"limit" binding:"omitempty,numeric,min=1"`
	Page       *int           `form:"page" binding:"omitempty,numeric,min=1"`
}

func (p *UserQueryParamReq) ToQuery() *valueobject.Query {
	query := valueobject.NewQuery()

	if p.Page != nil {
		query.WithPage(*p.Page)
	}
	if p.Limit != nil {
		query.WithLimit(*p.Limit)
	}

	if p.Order != nil {
		query.WithOrder(valueobject.Order(*p.Order))
	}

	if p.SortBy != nil {
		query.WithSortBy(*p.SortBy)
	} else {
		query.WithSortBy("id")
	}

	if p.Email != nil {
		query.Condition("email", valueobject.ILike, *p.Email)
	}
	if p.IsVerified != nil {
		query.Condition("is_verified", valueobject.Equal, *p.IsVerified)
	}
	if p.RoleId != nil {
		query.Condition("role_id", valueobject.Equal, *p.RoleId)
	}

	return query
}

type UserRes struct {
	Id         uint     `json:"id"`
	Email      string   `json:"email"`
	Role       *RoleRes `json:"role"`
	Name       string   `json:"name"`
	IsVerified bool     `json:"is_verified"`
}

func NewUserRes(u *entity.User) *UserRes {
	var name string = "Unverified User"
	role := RoleRes{Id: uint(u.RoleId), Name: u.Role.Name}
	if u.Profile != nil {
		name = u.Profile.Name
	}
	if u.AdminContact != nil {
		name = u.AdminContact.Name
	}
	return &UserRes{Id: u.Id, Email: u.Email, Role: &role, Name: name, IsVerified: u.IsVerified}
}

type RoleRes struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
