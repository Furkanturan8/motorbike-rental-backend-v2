package model

import "time"

type BaseModel struct {
	ID        int64      `json:"id" bun:",pk,autoincrement"`
	CreatedAt time.Time  `json:"created_at" bun:",nullzero,default:current_timestamp"`
	UpdatedAt time.Time  `json:"updated_at" bun:",nullzero,default:current_timestamp"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bun:",soft_delete,nullzero,default:null"` // Soft delete
}
