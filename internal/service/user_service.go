package service

import (
	"context"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/repository"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
)

type UserService struct {
	userRepo repository.IUserRepository
}

func NewUserService(u repository.IUserRepository) *UserService {
	return &UserService{
		userRepo: u,
	}
}

func (s *UserService) Create(ctx context.Context, user model.User) error {
	// Email kontrolü
	exists, err := s.userRepo.ExistsByEmail(ctx, user.Email)
	if err != nil {
		return errorx.Wrap(errorx.ErrDatabaseOperation, err)
	}
	if exists {
		return errorx.WithDetails(errorx.ErrUserAlreadyExists, "Bu e-posta adresi zaten kullanımda")
	}

	if err = s.userRepo.Create(ctx, &user); err != nil {
		return errorx.Wrap(errorx.ErrDatabaseOperation, err)
	}

	return nil
}

func (s *UserService) List(ctx context.Context) ([]model.User, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrDatabaseOperation, err)
	}

	return users, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errorx.WithDetails(errorx.ErrUserNotFound, "Belirtilen ID ile kullanıcı bulunamadı")
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, id int64, updatedUser model.User) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.WithDetails(errorx.ErrUserNotFound, "Güncellenecek kullanıcı bulunamadı")
	}

	if updatedUser.Email != "" {
		// Email değişiyorsa, yeni email'in başka bir kullanıcıda olmadığından emin ol
		if updatedUser.Email != user.Email {
			exists, err := s.userRepo.ExistsByEmail(ctx, updatedUser.Email)
			if err != nil {
				return errorx.Wrap(errorx.ErrDatabaseOperation, err)
			}
			if exists {
				return errorx.WithDetails(errorx.ErrUserAlreadyExists, "Bu e-posta adresi başka bir kullanıcı tarafından kullanılıyor")
			}
		}
	}

	if err = s.userRepo.Update(ctx, &updatedUser); err != nil {
		return errorx.Wrap(errorx.ErrDatabaseOperation, err)
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	// Önce kullanıcının var olup olmadığını kontrol et
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return errorx.WithDetails(errorx.ErrUserNotFound, "Silinecek kullanıcı bulunamadı")
	}

	if err = s.userRepo.Delete(ctx, id); err != nil {
		return errorx.Wrap(errorx.ErrDatabaseOperation, err)
	}

	return nil
}
