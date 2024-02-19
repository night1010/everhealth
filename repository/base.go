package repository

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BaseRepository[T any] interface {
	Find(ctx context.Context, query *valueobject.Query) ([]*T, error)
	FindOne(ctx context.Context, query *valueobject.Query) (*T, error)
	FindById(ctx context.Context, id uint) (*T, error)
	Create(ctx context.Context, t *T) (*T, error)
	Update(ctx context.Context, t *T) (*T, error)
	Delete(ctx context.Context, t *T) error
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func (r *baseRepository[T]) conn(ctx context.Context) *gorm.DB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *baseRepository[T]) Find(ctx context.Context, q *valueobject.Query) ([]*T, error) {
	var ts []*T
	query := r.conn(ctx).Model(ts)

	for _, s := range q.GetAssociations() {
		if s.Type == valueobject.AssociationTypeJoin {
			query.Joins(s.Entity)
		} else if s.Type == valueobject.AssociationTypePreload {
			query.Preload(s.Entity)
		}
	}

	for _, condition := range q.GetConditions() {
		sql := fmt.Sprintf("%s %s ?", condition.Field, condition.Operation)
		query.Where(sql, condition.Value)
	}

	if q.GetLimit() != nil {
		limit := *q.GetLimit()
		offset := countOffset(q.GetPage(), limit)
		query.Limit(limit).Offset(offset)
	}

	err := query.
		Order(q.GetOrder()).
		Find(&ts).
		Error
	if err != nil {
		return nil, err
	}
	return ts, nil
}

func (r *baseRepository[T]) FindOne(ctx context.Context, q *valueobject.Query) (*T, error) {
	conditions := q.GetConditions()
	var t *T
	query := r.conn(ctx).Model(t)
	if q.IsLocked() {
		query.Clauses(clause.Locking{Strength: "UPDATE"})
	}

	for _, s := range q.GetAssociations() {
		if s.Type == valueobject.AssociationTypeJoin {
			query.Joins(s.Entity)
		} else if s.Type == valueobject.AssociationTypePreload {
			query.Preload(s.Entity)
		}
	}

	for _, condition := range conditions {
		sql := fmt.Sprintf("%s %s ?", condition.Field, condition.Operation)
		query.Where(sql, condition.Value)
	}
	err := query.First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

func (r *baseRepository[T]) FindById(ctx context.Context, id uint) (*T, error) {
	var t *T
	err := r.conn(ctx).Where("id", id).First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return t, nil
}

func (r *baseRepository[T]) Create(ctx context.Context, t *T) (*T, error) {
	result := r.conn(ctx).Create(t)
	if result.Error != nil {
		return nil, result.Error
	}
	return t, nil
}

func (r *baseRepository[T]) Update(ctx context.Context, t *T) (*T, error) {
	result := r.conn(ctx).Model(t).Clauses(clause.Returning{}).Select("*").Updates(t)
	if result.Error != nil {
		return nil, result.Error
	}
	return t, nil
}
func (r *baseRepository[T]) Delete(ctx context.Context, t *T) error {
	result := r.conn(ctx).Delete(t)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *baseRepository[T]) paginate(ctx context.Context, query *valueobject.Query, modifier func(*gorm.DB) *gorm.DB) (*valueobject.PagedResult, error) {
	var ts []*T
	var totalItem int64

	itemQuery := r.conn(ctx).Model(ts)
	itemQuery = modifier(itemQuery)

	if query.GetLimit() != nil {
		limit := *query.GetLimit()
		offset := countOffset(query.GetPage(), limit)
		itemQuery.Limit(limit).Offset(offset)
	}

	err := itemQuery.Order(query.GetOrder()).Find(&ts).Error
	if err != nil {
		return nil, err
	}

	countQuery := r.conn(ctx).Model(ts)
	countQuery = modifier(countQuery)

	err = countQuery.Count(&totalItem).Error
	if err != nil {
		return nil, err
	}

	totalPage := countTotalPage(totalItem, query.GetLimit())
	currentPage := int(math.Min(float64(query.GetPage()), float64(totalPage)))

	result := &valueobject.PagedResult{
		Data:         ts,
		CurrentPage:  currentPage,
		CurrentItems: len(ts),
		TotalPage:    totalPage,
		TotalItem:    int(totalItem),
	}

	return result, nil
}
