package models

import "gorm.io/gorm"

type MailStatus uint

const (
	Failed MailStatus = iota
	Success
)

type Mail struct {
	gorm.Model
	UserId   uint       `json:"user_id" gorm:"unique" validate:"required"`
	TeamCode uint       `json:"team_code" validate:"required"`
	Email    string     `json:"email" gorm:"unique" validate:"required"`
	Status   MailStatus `json:"status" validate:"required"`
}
