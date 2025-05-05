package service

import (
	"context"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/repository"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
)

type BluetoothConnectionService struct {
	connRepo repository.IBluetoothConnectionRepository
}

func NewBluetoothConnectionService(repo repository.IBluetoothConnectionRepository) *BluetoothConnectionService {
	return &BluetoothConnectionService{connRepo: repo}
}

func (s *BluetoothConnectionService) Create(ctx context.Context, conn *model.BluetoothConnection) error {
	if err := s.connRepo.Create(ctx, conn); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}
	return nil
}

func (s *BluetoothConnectionService) GetByID(ctx context.Context, id int64) (*model.BluetoothConnection, error) {
	conn, err := s.connRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}
	return conn, nil
}

func (s *BluetoothConnectionService) Update(ctx context.Context, conn model.BluetoothConnection) error {
	if err := s.connRepo.Update(ctx, &conn); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}
	return nil
}

func (s *BluetoothConnectionService) Delete(ctx context.Context, id int64) error {
	if err := s.connRepo.Delete(ctx, id); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}
	return nil
}

func (s *BluetoothConnectionService) List(ctx context.Context) ([]model.BluetoothConnection, error) {
	conn, err := s.connRepo.List(ctx)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}
	return conn, nil
}
