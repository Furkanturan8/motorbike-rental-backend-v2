package service

import (
	"context"
	"strconv"
	"time"

	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/repository"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/errorx"
)

type RideService struct {
	rideRepo  repository.IRideRepository
	motorRepo repository.IMotorbikeRepository
}

func NewRideService(rideRepo repository.IRideRepository, motorRepo repository.IMotorbikeRepository) *RideService {
	return &RideService{
		rideRepo:  rideRepo,
		motorRepo: motorRepo,
	}
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
	return *rides, nil
}

func (s *RideService) ListByDateRange(ctx context.Context, startTime, endTime time.Time) (*[]model.Ride, error) {
	startDateStr := startTime.Format("2006-01-02")
	endDateStr := endTime.Format("2006-01-02")
	rides, err := s.rideRepo.ListByDateRange(ctx, startDateStr, endDateStr)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrInternal, err)
	}
	return &rides, nil
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

func (s *RideService) FinishRide(ctx context.Context, rideID int64, userID int64) (*model.Ride, error) {
	ride, err := s.rideRepo.GetByID(ctx, rideID)
	if err != nil {
		return nil, errorx.WrapErr(errorx.ErrNotFound, err)
	}

	if ride == nil {
		return nil, errorx.WrapMsg(errorx.ErrNotFound, "Sürüş bulunamadı")
	}
	if ride.UserID != userID {
		return nil, errorx.WrapMsg(errorx.ErrUnauthorized, "Bu sürüşe erişim yetkiniz yok.")
	}
	if ride.EndTime != nil && !ride.EndTime.IsZero() {
		return nil, errorx.WrapMsg(errorx.ErrInvalidRequest, "Sürüş zaten bitirildi!")
	}

	now := time.Now().UTC()
	startTime := ride.StartTime.UTC()

	// Süreyi hesapla ve mutlak değerini al
	duration := now.Sub(startTime)
	if duration < 0 {
		duration = -duration
	}

	minutes := int(duration.Minutes())
	if minutes < 0 {
		minutes = 0
	}
	ride.Duration = strconv.Itoa(int(duration.Seconds()))
	ride.Cost = float64(minutes*3) + 10
	ride.EndTime = &now

	motorbike, err := s.motorRepo.GetByID(ctx, ride.MotorbikeID)
	if err != nil {
		return nil, errorx.WrapMsg(errorx.ErrInternal, "Motorbike bilgileri alınamadı!")
	}
	if motorbike == nil {
		return nil, errorx.WrapMsg(errorx.ErrNotFound, "Motorbike bulunamadı")
	}

	if motorbike.LockStatus != model.Locked {
		return nil, errorx.WrapMsg(errorx.ErrInvalidRequest, "Motorbike kilitlenmedi! Lütfen önce kilitleyin!")
	}

	if err = s.rideRepo.Update(ctx, ride); err != nil {
		return nil, errorx.WrapMsg(errorx.ErrInternal, "Sürüş Bitirilemedi!")
	}

	updatedRide, _ := s.rideRepo.GetByID(ctx, rideID)
	return updatedRide, nil
}

func (s *RideService) HandleAfterPhotoUpload(ctx context.Context, rideID int) error {
	ride, err := s.rideRepo.GetByID(ctx, int64(rideID))
	if err != nil {
		return errorx.WrapMsg(errorx.ErrNotFound, "Sürüş bulunamadı")
	}

	if ride.Motorbike.LockStatus != model.Locked {
		return errorx.WrapMsg(errorx.ErrInvalidRequest, "Lütfen önce motoru kilitleyin")
	}

	// todo connectRepo yapınca hallet
	/*if err = s.connRepo.Disconnect(ctx, int(ride.MotorbikeID)); err != nil {
		return errorx.WrapMsg(errorx.ErrInternal, "Motor bağlantısı kesilemedi")
	}
	*/
	return nil
}
