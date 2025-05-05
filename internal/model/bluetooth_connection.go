package model

import (
	"github.com/uptrace/bun"
	"time"
)

type BluetoothConnection struct {
	bun.BaseModel `bun:"table:bluetooth_connections,alias:bluetooth_connection"`

	ID             int64      `json:"id" bun:",pk,autoincrement"`
	CreatedAt      time.Time  `json:"created_at" bun:",nullzero,default:current_timestamp"`
	UpdatedAt      time.Time  `json:"updated_at" bun:",nullzero,default:current_timestamp"`
	UserID         int64      `json:"user_id"`
	MotorbikeID    int64      `json:"motorbike_id"`
	ConnectedAt    time.Time  `json:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at"`
}
