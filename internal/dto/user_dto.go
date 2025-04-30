package dto

import (
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/utils"
)

type CreateUserRequest struct {
	Email     string       `json:"email" validate:"required_without=Phone,omitempty,max=64,email"`
	Phone     string       `json:"phone" validate:"required_without=Email,omitempty,max=64"`
	FirstName string       `json:"first_name" validate:"required,max=100"`
	LastName  string       `json:"last_name" validate:"required,max=100"`
	Password  string       `json:"password" validate:"required,min=3,max=100"`
	Status    model.Status `json:"status" validate:"omitempty,oneof=active inactive"`
	Role      model.Role   `json:"role"`
}

func (dto CreateUserRequest) ToDBModel(m model.User) model.User {
	m.Email = dto.Email
	if utils.ValidatePhone(dto.Phone) {
		m.Phone = dto.Phone
	} else {
		m.Phone = ""
	}
	m.FirstName = dto.FirstName
	m.LastName = dto.LastName
	if dto.Password != "" {
		_ = m.SetPassword(dto.Password)
	}
	if dto.Role == "" {
		m.Role = model.UserRole
	} else {
		m.Role = dto.Role
	}
	if dto.Status == "" {
		m.Status = model.StatusActive
	} else {
		m.Status = dto.Status
	}

	return m
}

type UpdateUserRequest struct {
	Email           string       `json:"email" validate:"omitempty,max=64,email"`
	Phone           string       `json:"phone" validate:"omitempty,max=64,phone"`
	FirstName       string       `json:"first_name" validate:"omitempty,max=100"`
	LastName        string       `json:"last_name" validate:"omitempty,max=100"`
	CurrentPassword string       `json:"current_password" validate:"omitempty,min=3,max=100"`
	NewPassword     string       `json:"new_password" validate:"omitempty,min=3,max=100"`
	Status          model.Status `json:"status" validate:"omitempty,oneof=active inactive"`
	Role            model.Role   `json:"role"`
}

func (dto UpdateUserRequest) ToDBModel(m model.User) model.User {
	m.Email = dto.Email
	if utils.ValidatePhone(dto.Phone) {
		m.Phone = dto.Phone
	} else {
		m.Phone = ""
	}
	m.FirstName = dto.FirstName
	m.LastName = dto.LastName
	if dto.Role == "" {
		m.Role = model.UserRole
	} else {
		m.Role = dto.Role
	}
	if dto.Status == "" {
		m.Status = model.StatusActive
	} else {
		m.Status = dto.Status
	}

	return m
}

type UserResponse struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Status    string `json:"status"`
}

func (dto UserResponse) ToResponseModel(m model.User) UserResponse {
	dto.ID = m.ID
	dto.Email = m.Email
	dto.Phone = m.Phone
	dto.FirstName = m.FirstName
	dto.LastName = m.LastName
	dto.Role = string(m.Role)
	dto.Status = string(m.Status)

	return dto
}
