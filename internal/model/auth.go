package model

import (
	"time"
)

// Token modeli
type Token struct {
	ID           int64     `json:"id" bun:",pk,autoincrement"`
	UserID       int64     `json:"user_id" bun:",notnull"`
	AccessToken  string    `json:"access_token" bun:",notnull"`
	RefreshToken string    `json:"refresh_token" bun:",notnull"`
	ExpiresAt    time.Time `json:"expires_at" bun:",notnull"`
	CreatedAt    time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
	RevokedAt    time.Time `json:"revoked_at,omitempty" bun:",nullzero"`
	User         *User     `json:"user,omitempty" bun:"rel:belongs-to,join:user_id=id"`
}

// Token durumları için yardımcı metodlar
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func (t *Token) IsRevoked() bool {
	return !t.RevokedAt.IsZero()
}

func (t *Token) IsValid() bool {
	return !t.IsExpired() && !t.IsRevoked()
}

// Blacklist modeli (geçersiz kılınan tokenlar için)
type TokenBlacklist struct {
	ID        int64     `json:"id" bun:",pk,autoincrement"`
	Token     string    `json:"token" bun:",notnull"`
	ExpiresAt time.Time `json:"expires_at" bun:",notnull"`
	CreatedAt time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
}

// Oturum modeli (aktif kullanıcı oturumları için)
type Session struct {
	ID           int64     `json:"id" bun:",pk,autoincrement"`
	UserID       int64     `json:"user_id" bun:",notnull"`
	RefreshToken string    `json:"refresh_token" bun:",notnull"`
	UserAgent    string    `json:"user_agent" bun:",notnull"`
	ClientIP     string    `json:"client_ip" bun:",notnull"`
	IsBlocked    bool      `json:"is_blocked" bun:",notnull,default:false"`
	ExpiresAt    time.Time `json:"expires_at" bun:",notnull"`
	CreatedAt    time.Time `json:"created_at" bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time `json:"updated_at" bun:",nullzero,notnull,default:current_timestamp"`
	User         *User     `json:"user,omitempty" bun:"rel:belongs-to,join:user_id=id"`
}

// Oturum durumu için yardımcı metodlar
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) IsValid() bool {
	return !s.IsExpired() && !s.IsBlocked
}
