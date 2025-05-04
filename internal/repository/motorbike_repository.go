package repository

import (
	"context"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/uptrace/bun"
)

type IMotorbikeRepository interface {
	Create(ctx context.Context, motorbike *model.Motorbike) error
	GetByID(ctx context.Context, id int64) (*model.Motorbike, error)
	Update(ctx context.Context, motorbike *model.Motorbike) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]model.Motorbike, error)
	GetMotorsForStatus(ctx context.Context, status string) ([]model.Motorbike, error)
}

type MotorbikeRepository struct {
	db *bun.DB
}

func NewMotorbikeRepository(db *bun.DB) IMotorbikeRepository {
	return &MotorbikeRepository{db: db}
}

func (r *MotorbikeRepository) Create(ctx context.Context, motorbike *model.Motorbike) error {
	_, err := r.db.NewInsert().Model(motorbike).Exec(ctx)
	return err
}

func (r *MotorbikeRepository) GetByID(ctx context.Context, id int64) (*model.Motorbike, error) {
	var motorbike model.Motorbike
	err := r.db.NewSelect().Model(&motorbike).Where("id = ?", id).Scan(ctx)
	return &motorbike, err
}

func (r *MotorbikeRepository) Update(ctx context.Context, motorbike *model.Motorbike) error {
	_, err := r.db.NewUpdate().Model(motorbike).WherePK().Exec(ctx)
	return err
}

func (r *MotorbikeRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*model.Motorbike)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *MotorbikeRepository) List(ctx context.Context) ([]model.Motorbike, error) {
	var motorbikes []model.Motorbike
	err := r.db.NewSelect().Model(&motorbikes).Scan(ctx)
	return motorbikes, err
}

func (r *MotorbikeRepository) GetMotorsForStatus(ctx context.Context, status string) ([]model.Motorbike, error) {
	var motorbikes []model.Motorbike
	if err := r.db.NewSelect().Model(&motorbikes).Where("status = ?", status).Scan(ctx); err != nil {
		return nil, err
	}
	return motorbikes, nil
}
