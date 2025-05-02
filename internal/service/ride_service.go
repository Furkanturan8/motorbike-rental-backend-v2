package service

import (
	"context"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/repository"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
	"strconv"
	"time"
)

type RideService struct {
	rideRepo  repository.IRideRepository
	motorRepo repository.MotorbikeRepository
}

func NewRideService(repo repository.IRideRepository) *RideService {
	return &RideService{rideRepo: repo}
}

func (s *RideService) Create(ctx context.Context, ride *model.Ride) error {
	if err := s.rideRepo.Create(ctx, ride); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}
	return nil
}

func (s *RideService) GetByID(ctx context.Context, id int64) (*model.Ride, error) {
	ride, err := s.rideRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}
	return ride, nil
}

func (s *RideService) Update(ctx context.Context, ride model.Ride) error {
	if err := s.rideRepo.Update(ctx, &ride); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}
	return nil
}

func (s *RideService) Delete(ctx context.Context, id int64) error {
	if err := s.rideRepo.Delete(ctx, id); err != nil {
		return errorx.WrapErr(errorx.ErrInternal, err)
	}
	return nil
}

func (s *RideService) List(ctx context.Context) ([]model.Ride, error) {
	rides, err := s.rideRepo.List(ctx)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}
	return rides, nil
}

func (s *RideService) GetByUserID(ctx context.Context, userID int64) ([]model.Ride, error) {
	rides, err := s.rideRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}
	return rides, nil
}

func (s *RideService) GetByMotorbikeID(ctx context.Context, motorbikeID int64) ([]model.Ride, error) {
	rides, err := s.rideRepo.ListByMotorbikeID(ctx, motorbikeID)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}
	return rides, nil
}

func (s *RideService) FinishRide(ctx context.Context, rideID int64, userID uint) (*model.Ride, error) {
	ride, err := s.rideRepo.GetByID(ctx, rideID)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrNotFound, err)
	}

	if ride.UserID != userID {
		return nil, errorx.WrapMsg(errorx.ErrUnauthorized, "Bu sürüşe erişim yetkiniz yok.")
	}

	if ride.EndTime != nil {
		return nil, errorx.WrapMsg(errorx.ErrInvalidRequest, "Sürüş zaten bitirildi!")
	}

	// Sürüş süresini hesapla// Sürüş süresini hesapla (StartTime bir pointer değilse)
	now := time.Now().UTC()
	ride.EndTime = &now
	duration := now.Sub(ride.StartTime)

	ride.Duration = strconv.Itoa(int(duration.Seconds()))
	ride.Cost = float64(int(duration.Minutes())*3) + 10 // Süreyi dakika cinsinden hesapla (her dakika için 3 TL)

	motorbike, err := s.motorRepo.GetByID(ctx, int64(ride.MotorbikeID))
	if err != nil {
		return nil, errorx.WrapMsg(errorx.ErrInternal, "Motorbike bilgileri alınamadı!")
	}
	if motorbike.LockStatus != model.Locked {
		return nil, errorx.WrapMsg(errorx.ErrInvalidRequest, "Motorbike kilitlenmedi! Lütfen önce kilitleyin!")
	}

	if err = s.rideRepo.Update(ctx, ride); err != nil {
		return nil, errorx.WrapMsg(errorx.ErrInternal, "Sürüş Bitirilemedi!")
	}

	return ride, nil
}
