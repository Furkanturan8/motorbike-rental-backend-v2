package dto

import (
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"time"
)

type CreateRideRequest struct {
	UserID      uint       `json:"user_id" validate:"required"`
	MotorbikeID uint       `json:"motorbike_id" validate:"required"`
	StartTime   time.Time  `json:"start_time" validate:"required"`
	EndTime     *time.Time `json:"end_time" validate:"required"`
	Duration    string     `json:"duration" validate:"required"`
	Cost        float64    `json:"cost" validate:"required"`
}

func (dto CreateRideRequest) ToDBModel(m model.Ride) model.Ride {
	m.UserID = dto.UserID
	m.MotorbikeID = dto.MotorbikeID
	m.StartTime = dto.StartTime
	m.EndTime = dto.EndTime
	m.Duration = dto.Duration
	m.Cost = dto.Cost
	return m
}

type UpdateRideRequest struct {
	UserID      uint       `json:"user_id"`
	MotorbikeID uint       `json:"motorbike_id"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Duration    string     `json:"duration"`
	Cost        float64    `json:"cost"`
}

func (dto UpdateRideRequest) ToDBModel(m model.Ride) model.Ride {
	m.UserID = dto.UserID
	m.MotorbikeID = dto.MotorbikeID
	m.StartTime = dto.StartTime
	m.EndTime = dto.EndTime
	m.Duration = dto.Duration
	m.Cost = dto.Cost
	return m
}

type RideResponse struct {
	ID          int64      `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UserID      uint       `json:"user_id"`
	MotorbikeID uint       `json:"motorbike_id"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Duration    string     `json:"duration"`
	Cost        float64    `json:"cost"`
}

func (dto RideResponse) ToResponseModel(m model.Ride) RideResponse {
	dto.ID = m.ID
	dto.CreatedAt = m.CreatedAt
	dto.UpdatedAt = m.UpdatedAt
	dto.UserID = m.UserID
	dto.MotorbikeID = m.MotorbikeID
	dto.StartTime = m.StartTime
	dto.EndTime = m.EndTime
	dto.Duration = m.Duration
	dto.Cost = m.Cost
	return dto
}
