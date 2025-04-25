package query

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/uptrace/bun"
	"strings"
)

// Sayfalama için varsayılan değerler
const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// Sıralama yönü
type SortDirection string

const (
	SortAsc  SortDirection = "ASC"
	SortDesc SortDirection = "DESC"
)

// Sayfalama bilgisi
type Pagination struct {
	Page       int   `json:"page" query:"page"`
	PageSize   int   `json:"page_size" query:"page_size"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

// Sıralama bilgisi
type Sort struct {
	Field     string        `json:"field" query:"sort_field"`
	Direction SortDirection `json:"direction" query:"sort_direction"`
}

// Filtre operatörleri
type FilterOperator string

const (
	Equal              FilterOperator = "eq"
	NotEqual           FilterOperator = "ne"
	GreaterThan        FilterOperator = "gt"
	GreaterThanOrEqual FilterOperator = "gte"
	LessThan           FilterOperator = "lt"
	LessThanOrEqual    FilterOperator = "lte"
	Like               FilterOperator = "like"
	ILike              FilterOperator = "ilike"
	In                 FilterOperator = "in"
	NotIn              FilterOperator = "not_in"
	IsNull             FilterOperator = "is_null"
	IsNotNull          FilterOperator = "is_not_null"
)

// Filtre yapısı
type Filter struct {
	Field    string         `json:"field" query:"filter_field"`
	Operator FilterOperator `json:"operator" query:"filter_operator"`
	Value    interface{}    `json:"value" query:"filter_value"`
}

// Query parametreleri
type Params struct {
	Pagination Pagination `json:"pagination"`
	Sort       []Sort     `json:"sort"`
	Filters    []Filter   `json:"filters"`
	Search     string     `json:"search" query:"search"`
}

// Fiber context'inden query parametrelerini okur
func ParseFromContext(c *fiber.Ctx) (*Params, error) {
	params := &Params{
		Pagination: Pagination{
			Page:     DefaultPage,
			PageSize: DefaultPageSize,
		},
	}

	// Sayfalama
	if page := c.QueryInt("page", DefaultPage); page > 0 {
		params.Pagination.Page = page
	}

	if pageSize := c.QueryInt("page_size", DefaultPageSize); pageSize > 0 && pageSize <= MaxPageSize {
		params.Pagination.PageSize = pageSize
	}

	// Sıralama
	if sortField := c.Query("sort_field"); sortField != "" {
		direction := SortDirection(strings.ToUpper(c.Query("sort_direction", string(SortAsc))))
		if direction != SortAsc && direction != SortDesc {
			direction = SortAsc
		}
		params.Sort = append(params.Sort, Sort{Field: sortField, Direction: direction})
	}

	// Arama
	if search := c.Query("search"); search != "" {
		params.Search = search
	}

	// Filtreler
	if filterField := c.Query("filter_field"); filterField != "" {
		operator := FilterOperator(c.Query("filter_operator", string(Equal)))
		value := c.Query("filter_value")
		params.Filters = append(params.Filters, Filter{
			Field:    filterField,
			Operator: operator,
			Value:    value,
		})
	}

	return params, nil
}

// Query Builder'a filtreleri uygular
func ApplyFilters(q *bun.SelectQuery, filters []Filter) *bun.SelectQuery {
	for _, filter := range filters {
		switch filter.Operator {
		case Equal:
			q = q.Where("? = ?", bun.Ident(filter.Field), filter.Value)
		case NotEqual:
			q = q.Where("? != ?", bun.Ident(filter.Field), filter.Value)
		case GreaterThan:
			q = q.Where("? > ?", bun.Ident(filter.Field), filter.Value)
		case GreaterThanOrEqual:
			q = q.Where("? >= ?", bun.Ident(filter.Field), filter.Value)
		case LessThan:
			q = q.Where("? < ?", bun.Ident(filter.Field), filter.Value)
		case LessThanOrEqual:
			q = q.Where("? <= ?", bun.Ident(filter.Field), filter.Value)
		case Like:
			q = q.Where("? LIKE ?", bun.Ident(filter.Field), fmt.Sprintf("%%%v%%", filter.Value))
		case ILike:
			q = q.Where("? ILIKE ?", bun.Ident(filter.Field), fmt.Sprintf("%%%v%%", filter.Value))
		case In:
			q = q.Where("? IN (?)", bun.Ident(filter.Field), filter.Value)
		case NotIn:
			q = q.Where("? NOT IN (?)", bun.Ident(filter.Field), filter.Value)
		case IsNull:
			q = q.Where("? IS NULL", bun.Ident(filter.Field))
		case IsNotNull:
			q = q.Where("? IS NOT NULL", bun.Ident(filter.Field))
		}
	}
	return q
}

// Query Builder'a sıralama uygular
func ApplySort(q *bun.SelectQuery, sorts []Sort) *bun.SelectQuery {
	for _, sort := range sorts {
		if sort.Direction == SortDesc {
			q = q.Order(fmt.Sprintf("%s DESC", sort.Field))
		} else {
			q = q.Order(fmt.Sprintf("%s ASC", sort.Field))
		}
	}
	return q
}

// Query Builder'a sayfalama uygular
func ApplyPagination(q *bun.SelectQuery, p Pagination) *bun.SelectQuery {
	offset := (p.Page - 1) * p.PageSize
	return q.Limit(p.PageSize).Offset(offset)
}

// Toplam kayıt sayısını hesaplar ve sayfalama bilgisini günceller
func UpdatePaginationInfo(ctx context.Context, q *bun.SelectQuery, p *Pagination) error {
	count, err := q.Count(ctx)
	if err != nil {
		return err
	}

	p.TotalRows = int64(count)
	p.TotalPages = (int(p.TotalRows) + p.PageSize - 1) / p.PageSize

	return nil
}

// Response için pagination bilgisini hazırlar
func GetPaginationResponse(p Pagination) map[string]interface{} {
	return map[string]interface{}{
		"current_page": p.Page,
		"page_size":    p.PageSize,
		"total_rows":   p.TotalRows,
		"total_pages":  p.TotalPages,
	}
}
