package repository

import (
	"context"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/internal/model"
	"github.com/uptrace/bun"
)

type IBluetoothConnectionRepository interface {
	Create(ctx context.Context, conn *model.BluetoothConnection) error
	GetByID(ctx context.Context, id int64) (*model.BluetoothConnection, error)
	GetByMotorbikeID(ctx context.Context, id int64) (*model.BluetoothConnection, error)
	Update(ctx context.Context, conn *model.BluetoothConnection) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]model.BluetoothConnection, error)
}

type BluetoothConnectionRepository struct {
	db *bun.DB
}

func NewBluetoothConnectionRepository(db *bun.DB) IBluetoothConnectionRepository {
	return &BluetoothConnectionRepository{db: db}
}

func (r *BluetoothConnectionRepository) Create(ctx context.Context, conn *model.BluetoothConnection) error {
	_, err := r.db.NewInsert().Model(conn).Exec(ctx)
	return err
}

func (r *BluetoothConnectionRepository) GetByID(ctx context.Context, id int64) (*model.BluetoothConnection, error) {
	var conn model.BluetoothConnection
	err := r.db.NewSelect().Model(&conn).Where("id = ?", id).Scan(ctx)
	return &conn, err
}

func (r *BluetoothConnectionRepository) GetByMotorbikeID(ctx context.Context, id int64) (*model.BluetoothConnection, error) {
	var conn model.BluetoothConnection
	err := r.db.NewSelect().Model(&conn).Where("motorbike_id = ?", id).Scan(ctx)
	return &conn, err
}

func (r *BluetoothConnectionRepository) Update(ctx context.Context, conn *model.BluetoothConnection) error {
	_, err := r.db.NewUpdate().Model(conn).WherePK().Exec(ctx)
	return err
}

func (r *BluetoothConnectionRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*model.BluetoothConnection)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *BluetoothConnectionRepository) List(ctx context.Context) ([]model.BluetoothConnection, error) {
	var conn []model.BluetoothConnection
	err := r.db.NewSelect().Model(&conn).Scan(ctx)
	return conn, err
}
