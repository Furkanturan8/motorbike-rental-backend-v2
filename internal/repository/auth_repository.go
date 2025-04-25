package repository

import (
	"context"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/uptrace/bun"
	"time"
)

type IAuthRepository interface {
	SaveToken(ctx context.Context, token *model.Token) error
	GetTokenByRefresh(ctx context.Context, refreshToken string) (*model.Token, error)
	RevokeToken(ctx context.Context, tokenID int64) error
	CreateSession(ctx context.Context, session *model.Session) error
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error)
	UpdateSession(ctx context.Context, session *model.Session) error
	DeleteSession(ctx context.Context, sessionID int64) error
	BlockSession(ctx context.Context, sessionID int64) error
	GetSessionsByUserID(ctx context.Context, userID int64) ([]*model.Session, error)
	AddToBlacklist(ctx context.Context, blacklist *model.TokenBlacklist) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	CleanupExpiredTokens(ctx context.Context) error
	CleanupExpiredSessions(ctx context.Context) error
	CreateUser(ctx context.Context, user *model.User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}

type AuthRepository struct {
	db *bun.DB
}

func NewAuthRepository(db *bun.DB) IAuthRepository {
	return &AuthRepository{db: db}
}

// Token işlemleri
func (r *AuthRepository) SaveToken(ctx context.Context, token *model.Token) error {
	_, err := r.db.NewInsert().Model(token).Exec(ctx)
	return err
}

func (r *AuthRepository) GetTokenByRefresh(ctx context.Context, refreshToken string) (*model.Token, error) {
	token := new(model.Token)
	err := r.db.NewSelect().
		Model(token).
		Where("refresh_token = ? AND revoked_at IS NULL", refreshToken).
		Relation("User").
		Scan(ctx)
	return token, err
}

func (r *AuthRepository) RevokeToken(ctx context.Context, tokenID int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Token)(nil)).
		Set("revoked_at = ?", time.Now()).
		Where("id = ?", tokenID).
		Exec(ctx)
	return err
}

// Session işlemleri
func (r *AuthRepository) CreateSession(ctx context.Context, session *model.Session) error {
	_, err := r.db.NewInsert().Model(session).Exec(ctx)
	return err
}

func (r *AuthRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error) {
	session := new(model.Session)
	err := r.db.NewSelect().
		Model(session).
		Where("refresh_token = ? AND is_blocked = false", refreshToken).
		Relation("User").
		Scan(ctx)
	return session, err
}

func (r *AuthRepository) UpdateSession(ctx context.Context, session *model.Session) error {
	_, err := r.db.NewUpdate().
		Model(session).
		WherePK().
		Exec(ctx)
	return err
}

func (r *AuthRepository) DeleteSession(ctx context.Context, sessionID int64) error {
	_, err := r.db.NewDelete().
		Model((*model.Session)(nil)).
		Where("id = ?", sessionID).
		Exec(ctx)
	return err
}

func (r *AuthRepository) BlockSession(ctx context.Context, sessionID int64) error {
	_, err := r.db.NewUpdate().
		Model((*model.Session)(nil)).
		Set("is_blocked = true").
		Where("id = ?", sessionID).
		Exec(ctx)
	return err
}

func (r *AuthRepository) GetSessionsByUserID(ctx context.Context, userID int64) ([]*model.Session, error) {
	var sessions []*model.Session
	err := r.db.NewSelect().
		Model(&sessions).
		Where("user_id = ? AND is_blocked = false", userID).
		Scan(ctx)
	return sessions, err
}

// Token Blacklist işlemleri
func (r *AuthRepository) AddToBlacklist(ctx context.Context, blacklist *model.TokenBlacklist) error {
	_, err := r.db.NewInsert().Model(blacklist).Exec(ctx)
	return err
}

func (r *AuthRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	exists, err := r.db.NewSelect().
		Model((*model.TokenBlacklist)(nil)).
		Where("token = ? AND expires_at > ?", token, time.Now()).
		Exists(ctx)
	return exists, err
}

// Temizlik işlemleri
func (r *AuthRepository) CleanupExpiredTokens(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*model.TokenBlacklist)(nil)).
		Where("expires_at < ?", time.Now()).
		Exec(ctx)
	return err
}

func (r *AuthRepository) CleanupExpiredSessions(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*model.Session)(nil)).
		Where("expires_at < ?", time.Now()).
		Exec(ctx)
	return err
}

// User işlemleri
func (r *AuthRepository) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *AuthRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.db.NewSelect().
		Model((*model.User)(nil)).
		Where("email = ?", email).
		Exists(ctx)
	return exists, err
}

func (r *AuthRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)
	err := r.db.NewSelect().
		Model(user).
		Where("email = ?", email).
		Scan(ctx)
	return user, err
}

func (r *AuthRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user := new(model.User)
	err := r.db.NewSelect().
		Model(user).
		Where("id = ?", id).
		Scan(ctx)
	return user, err
}

func (r *AuthRepository) Update(ctx context.Context, user *model.User) error {
	_, err := r.db.NewUpdate().
		Model(user).
		WherePK().
		Exec(ctx)
	return err
}
