package service

import (
	"context"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/repository"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
)

type RideService struct {
	rideRepo repository.IRideRepository
}

func NewRideService(repo repository.IRideRepository) *RideService {
	return &RideService{rideRepo: repo}
}

func (s *RideService) Create(ctx context.Context, ride *model.Ride) error {
	if err := s.rideRepo.Create(ctx, ride); err != nil {
		return errorx.Wrap(errorx.ErrDatabaseOperation, err)
	}
	return nil
}

func (s *RideService) GetByID(ctx context.Context, id int64) (*model.Ride, error) {
	ride, err := s.rideRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrDatabaseOperation, err)
	}
	return ride, nil
}

func (s *RideService) Update(ctx context.Context, ride model.Ride) error {
	if err := s.rideRepo.Update(ctx, &ride); err != nil {
		return errorx.Wrap(errorx.ErrDatabaseOperation, err)
	}
	return nil
}

func (s *RideService) Delete(ctx context.Context, id int64) error {
	if err := s.rideRepo.Delete(ctx, id); err != nil {
		return errorx.Wrap(errorx.ErrDatabaseOperation, err)
	}
	return nil
}

func (s *RideService) List(ctx context.Context) ([]model.Ride, error) {
	rides, err := s.rideRepo.List(ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrDatabaseOperation, err)
	}
	return rides, nil
}

func (s *RideService) GetByUserID(ctx context.Context, userID int64) ([]model.Ride, error) {
	rides, err := s.rideRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrDatabaseOperation, err)
	}
	return rides, nil
}

func (s *RideService) GetByMotorbikeID(ctx context.Context, motorbikeID int64) ([]model.Ride, error) {
	rides, err := s.rideRepo.ListByMotorbikeID(ctx, motorbikeID)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrDatabaseOperation, err)
	}
	return rides, nil
}
