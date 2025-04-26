package repository

import (
	"context"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/uptrace/bun"
)

type IRideRepository interface {
	Create(ctx context.Context, ride *model.Ride) error
	GetByID(ctx context.Context, id int64) (*model.Ride, error)
	Update(ctx context.Context, ride *model.Ride) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]model.Ride, error)
}

type RideRepository struct {
	db *bun.DB
}

func NewRideRepository(db *bun.DB) IRideRepository {
	return &RideRepository{db: db}
}

func (r *RideRepository) Create(ctx context.Context, ride *model.Ride) error {
	_, err := r.db.NewInsert().Model(ride).Exec(ctx)
	return err
}

func (r *RideRepository) GetByID(ctx context.Context, id int64) (*model.Ride, error) {
	var ride model.Ride
	err := r.db.NewSelect().Model(&ride).Where("id = ?", id).Scan(ctx)
	return &ride, err
}

func (r *RideRepository) Update(ctx context.Context, ride *model.Ride) error {
	_, err := r.db.NewUpdate().Model(ride).WherePK().Exec(ctx)
	return err
}

func (r *RideRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*model.Ride)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *RideRepository) List(ctx context.Context) ([]model.Ride, error) {
	var rides []model.Ride
	err := r.db.NewSelect().Model(&rides).Scan(ctx)
	return rides, err
}
