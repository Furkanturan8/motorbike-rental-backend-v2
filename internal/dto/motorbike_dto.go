package dto

import (
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"time"
)

// Fotoğraflar için ayrı bir dto
type PhotoCreateDto struct {
	PhotoURL string `json:"photo_url" validate:"required,url"`
}

type CreateMotorbikeRequest struct {
	Model             string           `json:"model" validate:"required"`
	LocationLatitude  float64          `json:"location_latitude" validate:"required"`
	LocationLongitude float64          `json:"location_longitude" validate:"required"`
	Status            string           `json:"status" validate:"required,oneof=available maintenance rented"`
	Photos            []PhotoCreateDto `json:"photos"`
	LockStatus        string           `json:"lock_status" validate:"required,oneof=locked unlocked"`
}

func (dto CreateMotorbikeRequest) ToDBModel(m model.Motorbike) model.Motorbike {
	m.Model = dto.Model
	m.LocationLatitude = dto.LocationLatitude
	m.LocationLongitude = dto.LocationLongitude
	m.Status = model.MotorBikeStatus(dto.Status)
	m.LockStatus = model.LockStatus(dto.LockStatus)

	return m
}

// Fotoğraf modellerine dönüştürme
func (dto CreateMotorbikeRequest) ToPhotoModels(motorbikeID int) []model.MotorbikePhoto {
	var photos []model.MotorbikePhoto
	for _, photoVM := range dto.Photos {
		photos = append(photos, model.MotorbikePhoto{
			MotorbikeID: motorbikeID,
			PhotoURL:    photoVM.PhotoURL,
		})
	}
	return photos
}

type UpdateMotorbikeRequest struct {
	Model             string           `json:"model"`
	LocationLatitude  float64          `json:"location_latitude"`
	LocationLongitude float64          `json:"location_longitude"`
	Status            string           `json:"status" validate:"required,oneof=available maintenance rented"`
	Photos            []PhotoCreateDto `json:"photos"`
	LockStatus        string           `json:"lock_status" validate:"required,oneof=locked unlocked"`
}

func (dto UpdateMotorbikeRequest) ToDBModel(m model.Motorbike) model.Motorbike {
	m.Model = dto.Model
	m.LocationLatitude = dto.LocationLatitude
	m.LocationLongitude = dto.LocationLongitude
	m.Status = model.MotorBikeStatus(dto.Status)
	m.LockStatus = model.LockStatus(dto.LockStatus)

	return m
}

// Fotoğraf modellerine dönüştürme
func (dto UpdateMotorbikeRequest) ToPhotoModels(motorbikeID int) []model.MotorbikePhoto {
	var photos []model.MotorbikePhoto
	for _, photoVM := range dto.Photos {
		photos = append(photos, model.MotorbikePhoto{
			MotorbikeID: motorbikeID,
			PhotoURL:    photoVM.PhotoURL,
		})
	}
	return photos
}

// Fotoğraf detayları için dto
type PhotoDetailDto struct {
	ID          int    `json:"id"`
	MotorbikeID int    `json:"motorbike_id"`
	PhotoURL    string `json:"photo_url"`
}

type MotorbikeResponse struct {
	ID                int64            `json:"id"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	Model             string           `json:"model"`
	LocationLatitude  float64          `json:"location_latitude"`
	LocationLongitude float64          `json:"location_longitude"`
	Status            string           `json:"status"`
	Photos            []PhotoDetailDto `json:"photos"`
	LockStatus        string           `json:"lock_status"`
}

func (dto MotorbikeResponse) ToResponseModel(m model.Motorbike) MotorbikeResponse {
	var photoDTOs []PhotoDetailDto
	for _, photo := range m.Photos {
		photoDTOs = append(photoDTOs, PhotoDetailDto{
			ID:          int(photo.ID),
			MotorbikeID: photo.MotorbikeID,
			PhotoURL:    photo.PhotoURL,
		})
	}

	dto.ID = m.ID
	dto.CreatedAt = m.CreatedAt
	dto.UpdatedAt = m.UpdatedAt
	dto.Model = m.Model
	dto.LocationLatitude = m.LocationLatitude
	dto.LocationLongitude = m.LocationLongitude
	dto.Status = string(m.Status)
	dto.LockStatus = string(m.LockStatus)
	dto.Photos = photoDTOs

	return dto
}
