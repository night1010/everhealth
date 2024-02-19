package entity

import (
	"time"

	"gorm.io/gorm"
)

type RoleId uint

const (
	RoleUser RoleId = iota + 1
	RoleDoctor
	RoleAdmin
	RoleSuperAdmin
)

type Role struct {
	Id        RoleId `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
