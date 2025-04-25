package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Role string
type Status string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusBanned   Status = "banned"
)

type User struct {
	ID        int64     `json:"id" bun:",pk,autoincrement"`
	Email     string    `json:"email" bun:",unique,notnull"`
	Password  string    `json:"-" bun:"password_hash,notnull"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      Role      `json:"role" bun:"type:user_role,notnull,default:'user'"`
	Status    Status    `json:"status" bun:"type:user_status,notnull,default:'active'"`
	LastLogin time.Time `json:"last_login" bun:",nullzero"`
	CreatedAt time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`

	tableName struct{} `bun:"users"`
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) GetStatus() Status {
	return u.Status
}
