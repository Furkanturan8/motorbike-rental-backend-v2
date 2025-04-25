package repository

import (
	"context"
	"fmt"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/cache"

	"github.com/uptrace/bun"
	"time"
)

const (
	userCacheKeyPrefix = "user:"
	userListCacheKey   = "users:list"
	userCacheDuration  = 24 * time.Hour
)

type IUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
	UpdateLastLogin(ctx context.Context, id int64) error
	List(ctx context.Context) ([]model.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("veritabanı insert hatası: %v", err)
	}

	// Yeni kullanıcıyı cache'e ekle
	cacheKey := fmt.Sprintf("%s%d", userCacheKeyPrefix, user.ID)
	if err = cache.Set(ctx, cacheKey, user, userCacheDuration); err != nil {
		// Cache hatası loglansın ama işlemi engellemeyelim
	}

	// Liste cache'ini temizle çünkü yeni kullanıcı eklendi
	cache.Delete(ctx, userListCacheKey)
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user := model.User{}
	err := r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()
	// Sadece değişen alanları güncelle
	_, err := r.db.NewUpdate().
		Model(user).
		WherePK().
		Column("email", "first_name", "last_name", "password_hash", "role", "status", "updated_at").
		Exec(ctx)
	if err != nil {
		return err
	}

	// Güncellenen kullanıcıyı cache'e ekle
	cacheKey := fmt.Sprintf("%s%d", userCacheKeyPrefix, user.ID)
	if err = cache.Set(ctx, cacheKey, user, userCacheDuration); err != nil {
		// Cache hatası loglansın ama işlemi engellemeyelim
		fmt.Println("Cache güncelleme hatası:", err)
	}

	// Liste cache'ini temizle çünkü kullanıcı güncellendi
	cache.Delete(ctx, userListCacheKey)
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*model.User)(nil)).Where("id = ?", id).Exec(ctx)
	if err != nil {
		return err
	}

	// Silinen kullanıcıyı cache'den sil
	cacheKey := fmt.Sprintf("%s%d", userCacheKeyPrefix, id)
	cache.Delete(ctx, cacheKey)

	// Liste cache'ini temizle çünkü kullanıcı silindi
	cache.Delete(ctx, userListCacheKey)
	return nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id int64) error {
	user := &model.User{ID: id}
	_, err := r.db.NewUpdate().
		Model(user).
		Column("last_login").
		WherePK().
		Exec(ctx)

	if err != nil {
		return err
	}

	// Son giriş tarihini güncelle
	user.LastLogin = time.Now()

	// Güncellenen kullanıcıyı cache'e ekle
	cacheKey := fmt.Sprintf("%s%d", userCacheKeyPrefix, id)
	if err = cache.Set(ctx, cacheKey, user, userCacheDuration); err != nil {
		// Cache hatası loglansın ama işlemi engellemeyelim
	}

	return nil
}

func (r *UserRepository) List(ctx context.Context) ([]model.User, error) {
	// Önce cache'den kontrol et
	var users []model.User
	err := cache.Get(ctx, userListCacheKey, &users)
	if err == nil {
		fmt.Printf("Kullanıcılar cache'den alındı\n")
		return users, nil
	}

	err = r.db.NewSelect().Model(&users).Scan(ctx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Kullanıcılar veritabanından alındı\n")

	// Cache'e kaydet
	if err = cache.Set(ctx, userListCacheKey, &users, userCacheDuration); err != nil {
		// Cache hatası loglansın ama işlemi engellemeyecek
		return users, nil
	}
	return users, nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	exists, err := r.db.NewSelect().
		Model((*model.User)(nil)).
		Where("email = ?", email).
		Exists(ctx)

	if err != nil {
		return false, err
	}

	return exists, nil
}
