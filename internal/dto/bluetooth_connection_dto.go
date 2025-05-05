package dto

import (
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"time"
)

type CreateBluetoothConnectionRequest struct {
	UserID         int64      `json:"user_id" validate:"required"`
	MotorbikeID    int64      `json:"motorbike_id" validate:"required"`
	ConnectedAt    time.Time  `json:"connected_at" validate:"required"`
	DisconnectedAt *time.Time `json:"disconnected_at" validate:"required"`
}

func (dto CreateBluetoothConnectionRequest) ToDBModel(m model.BluetoothConnection) model.BluetoothConnection {
	m.UserID = dto.UserID
	m.MotorbikeID = dto.MotorbikeID
	m.ConnectedAt = dto.ConnectedAt
	m.DisconnectedAt = dto.DisconnectedAt

	return m
}

type UpdateBluetoothConnectionRequest struct {
	UserID         int64      `json:"user_id"`
	MotorbikeID    int64      `json:"motorbike_id"`
	ConnectedAt    time.Time  `json:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at"`
}

func (dto UpdateBluetoothConnectionRequest) ToDBModel(m model.BluetoothConnection) model.BluetoothConnection {
	m.UserID = dto.UserID
	m.MotorbikeID = dto.MotorbikeID
	m.ConnectedAt = dto.ConnectedAt
	m.DisconnectedAt = dto.DisconnectedAt

	return m
}

type BluetoothConnectionResponse struct {
	ID             int64      `json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	UserID         int64      `json:"user_id"`
	MotorbikeID    int64      `json:"motorbike_id"`
	ConnectedAt    time.Time  `json:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at"`
}

func (dto BluetoothConnectionResponse) ToResponseModel(m model.BluetoothConnection) BluetoothConnectionResponse {
	dto.ID = m.ID
	dto.CreatedAt = m.CreatedAt
	dto.UpdatedAt = m.UpdatedAt
	dto.UserID = m.UserID
	dto.MotorbikeID = m.MotorbikeID
	dto.ConnectedAt = m.ConnectedAt
	dto.DisconnectedAt = m.DisconnectedAt
	return dto
}
