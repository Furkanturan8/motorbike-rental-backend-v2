package service

import (
	"context"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/repository"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
)

type MotorbikeService struct {
	motorbikeRepo repository.IMotorbikeRepository
}

func NewMotorbikeService(repo repository.IMotorbikeRepository) *MotorbikeService {
	return &MotorbikeService{motorbikeRepo: repo}
}

func (s *MotorbikeService) Create(ctx context.Context, motorbike *model.Motorbike) error {
	if err := s.motorbikeRepo.Create(ctx, motorbike); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}
	return nil
}

func (s *MotorbikeService) GetByID(ctx context.Context, id int64) (*model.Motorbike, error) {
	motorbike, err := s.motorbikeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}
	return motorbike, nil
}

func (s *MotorbikeService) Update(ctx context.Context, motorbike model.Motorbike) error {
	if err := s.motorbikeRepo.Update(ctx, &motorbike); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}
	return nil
}

func (s *MotorbikeService) Delete(ctx context.Context, id int64) error {
	if err := s.motorbikeRepo.Delete(ctx, id); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}
	return nil
}

func (s *MotorbikeService) List(ctx context.Context) ([]model.Motorbike, error) {
	motorbikes, err := s.motorbikeRepo.List(ctx)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}
	return motorbikes, nil
}
