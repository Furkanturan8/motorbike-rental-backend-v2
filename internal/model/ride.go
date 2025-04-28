package model

import (
	"time"
)

type Ride struct {
	BaseModel `bun:"table:rides"`

	UserID      uint       `json:"user_id" bun:"not null"`
	MotorbikeID uint       `json:"motorbike_id" bun:"not null"`
	StartTime   time.Time  `json:"start_time" bun:"default:current_timestamp"`
	EndTime     *time.Time `json:"end_time"`
	Duration    string     `json:"duration"`
	Cost        float64    `json:"cost"`

	User      User      `bun:"rel:belongs-to,join:user_id=id"`
	Motorbike Motorbike `bun:"rel:belongs-to,join:motorbike_id=id"`
}
