package models

import "gopkg.in/go-playground/validator.v9"

var validate *validator.Validate

type Message struct {
	UserId   uint   `json:"user_id" validate:"required"`
	TeamCode uint   `json:"team_code" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Counter  uint   `json:"counter" validate:"required"`
}
