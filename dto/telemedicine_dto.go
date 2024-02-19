package dto

import (
	"errors"
	"time"

	"github.com/night1010/everhealth/apperror"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/logger"
	"github.com/night1010/everhealth/valueobject"
	"github.com/shopspring/decimal"
)

type TelemedicineCreateReq struct {
	DoctorId uint `json:"doctor_id" binding:"required,numeric,min=1"`
}

type SickLeaveReq struct {
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Diagnosa  string `json:"diagnosa" binding:"required,max=100"`
}

type SickLeave struct {
	StartDate time.Time
	EndDate   time.Time
	Diagnosa  string
}

func (r *SickLeaveReq) ToModel() (*SickLeave, error) {
	startDate, err := time.Parse("2006-01-02", r.StartDate)
	if err != nil {
		return nil, apperror.NewClientError(err)
	}

	endDate, err := time.Parse("2006-01-02", r.EndDate)
	if err != nil {
		return nil, apperror.NewClientError(err)
	}
	logger.Log.Info(startDate)
	if !startDate.Before(endDate) {
		return nil, apperror.NewClientError(errors.New("start Date must be less than end Date"))
	}
	return &SickLeave{StartDate: startDate, EndDate: endDate, Diagnosa: r.Diagnosa}, nil
}

type Prescription struct {
	Prescription string `json:"prescription" binding:"required"`
}

type TelemedicineParams struct {
	Name   *string `form:"name"`
	Order  *string `form:"order" binding:"omitempty,oneof=asc desc"`
	Limit  *int    `form:"limit" binding:"omitempty,numeric,min=1"`
	Page   *int    `form:"page" binding:"omitempty,numeric,min=1"`
	Status *string `form:"status" binding:"omitempty,oneof=ongoing ended"`
}

func (qp *TelemedicineParams) ToQuery() *valueobject.Query {
	query := valueobject.NewQuery()
	if qp.Name != nil {
		query.Condition("name", valueobject.ILike, *qp.Name)
	}
	if qp.Status != nil {
		query.Condition("status", valueobject.Equal, *qp.Status)
	}
	if qp.Page != nil {
		query.WithPage(*qp.Page)
	}
	if qp.Limit != nil {
		query.WithLimit(*qp.Limit)
	}
	if qp.Order != nil {
		query.WithOrder(valueobject.Order(*qp.Order))
	}

	return query
}

type TelemedicineRes struct {
	Id              uint                      `json:"id"`
	OrderedAt       time.Time                 `json:"ordered_at"`
	ExpiredAt       time.Time                 `json:"expired_at"`
	Status          entity.TelemedicineStatus `json:"status"`
	TotalPayment    decimal.Decimal           `json:"total_payment"`
	Proof           string                    `json:"proof"`
	SickLeavePdf    string                    `json:"sick_leave_pdf"`
	PrescriptionPdf string                    `json:"prescription_pdf"`
	Profile         *ProfileRes               `json:"profile"`
	Doctor          *DoctorRes                `json:"doctor"`
	Chats           []*ChatRes                `json:"chats"`
}

type ProfileRes struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	BirthDate string `json:"birth_date"`
}

type DoctorRes struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

type ChatRes struct {
	UserId      uint               `json:"user_id"`
	ChatTime    time.Time          `json:"chat_time"`
	Message     string             `json:"message"`
	MessageType entity.MessageType `json:"message_type"`
}

type TelemedicineUri struct {
	Id uint `uri:"id" binding:"required,numeric"`
}

func NewTelemedicinRes(u *entity.Telemedicine) *TelemedicineRes {
	var profile *ProfileRes
	var doctor *DoctorRes
	if u.Profile != nil {
		profile = NewProfileTeleRes(u.Profile)
	}
	if u.Doctor != nil {
		doctor = NewDoctorTeleRes(u.Doctor)
	}
	return &TelemedicineRes{Id: u.Id, OrderedAt: u.OrderedAt, ExpiredAt: u.ExpiredAt, Status: u.Status, TotalPayment: u.TotalPayment, Proof: u.Proof, PrescriptionPdf: u.PrescriptionPdf, SickLeavePdf: u.SickLeavePdf, Doctor: doctor, Profile: profile, Chats: NewChatsRes(u.Chats)}
}

func NewProfileTeleRes(u *entity.Profile) *ProfileRes {
	return &ProfileRes{Id: u.UserId, Name: u.Name, Image: u.Image, BirthDate: u.Birthdate.Format("2006-01-02")}
}

func NewDoctorTeleRes(u *entity.DoctorProfile) *DoctorRes {
	return &DoctorRes{Id: u.ProfileId, Name: u.Profile.Name, Image: u.Profile.Image, Status: entity.DoctorStatusMap[u.Status]}
}

func NewChatRes(c *entity.Chat) *ChatRes {
	return &ChatRes{
		UserId:      c.UserId,
		ChatTime:    c.ChatTime,
		Message:     c.Message,
		MessageType: c.MessageType,
	}
}

func NewChatsRes(cs []*entity.Chat) []*ChatRes {
	return newResponsesFromEntities(cs, NewChatRes)
}
