package repository

import (
	"context"
	"github.com/Furkanturan8/motorbike-rental-backend-v2/pkg/query"
	"github.com/uptrace/bun"
)

type IBaseRepository interface {
	List(ctx context.Context, model interface{}, params *query.Params) error
	Get(ctx context.Context, model interface{}, id int64) error
	Create(ctx context.Context, model interface{}) error
	Update(ctx context.Context, model interface{}) error
	Delete(ctx context.Context, model interface{}, id int64) error
}

type BaseRepository struct {
	db *bun.DB
}

func NewBaseRepository(db *bun.DB) BaseRepository {
	return BaseRepository{db: db}
}

// Liste sorgularını hazırlar ve çalıştırır
func (r *BaseRepository) List(ctx context.Context, model interface{}, params *query.Params) error {
	q := r.db.NewSelect().Model(model)

	// Filtreleri uygula
	if len(params.Filters) > 0 {
		q = query.ApplyFilters(q, params.Filters)
	}

	// Sıralamayı uygula
	if len(params.Sort) > 0 {
		q = query.ApplySort(q, params.Sort)
	}

	// Toplam kayıt sayısını hesapla
	if err := query.UpdatePaginationInfo(ctx, q, &params.Pagination); err != nil {
		return err
	}

	// Sayfalamayı uygula
	q = query.ApplyPagination(q, params.Pagination)

	// Sorguyu çalıştır
	return q.Scan(ctx)
}

// Tekil kayıt sorgularını hazırlar ve çalıştırır
func (r *BaseRepository) Get(ctx context.Context, model interface{}, id int64) error {
	return r.db.NewSelect().
		Model(model).
		Where("id = ?", id).
		Scan(ctx)
}

// Kayıt oluşturur
func (r *BaseRepository) Create(ctx context.Context, model interface{}) error {
	_, err := r.db.NewInsert().
		Model(model).
		Exec(ctx)
	return err
}

// Kayıt günceller
func (r *BaseRepository) Update(ctx context.Context, model interface{}) error {
	_, err := r.db.NewUpdate().
		Model(model).
		WherePK().
		Exec(ctx)
	return err
}

// Kayıt siler
func (r *BaseRepository) Delete(ctx context.Context, model interface{}, id int64) error {
	_, err := r.db.NewDelete().
		Model(model).
		Where("id = ?", id).
		Exec(ctx)
	return err
}
