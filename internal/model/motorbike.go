package model

import (
	"github.com/uptrace/bun"
	"time"
)

type MotorBikeStatus string

const (
	BikeAvailable     MotorBikeStatus = "available"
	BikeInMaintenance MotorBikeStatus = "maintenance"
	BikeRented        MotorBikeStatus = "rented"
)

type LockStatus string

const (
	Locked   LockStatus = "locked"
	Unlocked LockStatus = "unlocked"
)

// Motorbike modeli
type Motorbike struct {
	bun.BaseModel `bun:"table:motorbikes,alias:motorbike"`

	ID                int64            `json:"id" bun:",pk,autoincrement"`
	CreatedAt         time.Time        `json:"created_at" bun:",nullzero,default:current_timestamp"`
	UpdatedAt         time.Time        `json:"updated_at" bun:",nullzero,default:current_timestamp"`
	Model             string           `json:"model"`
	LocationLatitude  float64          `json:"location_latitude"`
	LocationLongitude float64          `json:"location_longitude"`
	Photos            []MotorbikePhoto `json:"photos" bun:"rel:has-many,join:ID=MotorbikeID"`
	Status            MotorBikeStatus  `json:"status" bun:"type:varchar(20)"`
	LockStatus        LockStatus       `json:"lock_status" bun:"type:varchar(10)"`
}

type MotorbikePhoto struct {
	bun.BaseModel `bun:"table:motorbike_photos,alias:photo"`
	ID            int64     `json:"id" bun:",pk,autoincrement"`
	CreatedAt     time.Time `json:"created_at" bun:",nullzero,default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" bun:",nullzero,default:current_timestamp"`

	MotorbikeID int    `gorm:"not null"`
	PhotoURL    string `gorm:"type:varchar(255);not null"`
}

func (Motorbike) TableName() string {
	return "motorbike"
}

func (r MotorBikeStatus) String() string {
	switch r {
	case BikeAvailable:
		return "available"
	case BikeInMaintenance:
		return "maintenance"
	case BikeRented:
		return "rented"
	default:
		return "unknown"
	}
}

func (r LockStatus) String() string {
	switch r {
	case Locked:
		return "locked"
	case Unlocked:
		return "unlocked"
	default:
		return "unknown"
	}
}
