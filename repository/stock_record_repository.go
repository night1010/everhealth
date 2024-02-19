package repository

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
)

type StockRecordRepository interface {
	BaseRepository[entity.StockRecord]
	FindAllStockRecord(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	MonthlyReport(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	BulkCreate(ctx context.Context, records []*entity.StockRecord) error
}

type stockRecordRepository struct {
	*baseRepository[entity.StockRecord]
	db *gorm.DB
}

func NewStockRecordRepository(db *gorm.DB) StockRecordRepository {
	return &stockRecordRepository{
		db:             db,
		baseRepository: &baseRepository[entity.StockRecord]{db: db},
	}
}

func (r *stockRecordRepository) FindAllStockRecord(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "id":
			query.WithSortBy("\"stock_records\".id ")
		case "name":
			query.WithSortBy("\"PharmacyProduct__Product\".name")
		}
		db.Joins("PharmacyProduct.Product.ProductCategory").Joins("PharmacyProduct.Pharmacy")

		isReduction := query.GetConditionValue("is_reduction")
		name := query.GetConditionValue("name")
		ProductId := query.GetConditionValue("PharmacyProductId")
		db.Where("\"PharmacyProduct__Pharmacy\".admin_id =?", ctx.Value("user_id").(uint))
		if isReduction != nil {
			db.Where("\"stock_records\".is_reduction = ?", isReduction)
		}
		if ProductId != nil {
			db.Where("\"stock_records\".pharmacy_product_id", ProductId)
		}
		if name != nil {
			db.Where("\"PharmacyProduct__Product\".name ILIKE ?", name)
		}
		return db
	})
}

func (r *stockRecordRepository) MonthlyReport(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	var totalItem int64
	result := []*dto.ReportRes{}
	newQuery := `    with calc as (SELECT sr.pharmacy_product_id,
		sr.change_at,
		CASE
			When sr.is_reduction = true then quantity
			when sr.is_reduction = false then 0
			end as reduction_qty,
		case
			when sr.is_reduction = false then quantity
			when sr.is_reduction = true then 0
			end as additon_qty
 From stock_records sr)`
	rawQuery := newQuery
	newQuery = newQuery +
		` select sum(calc.additon_qty)                           as additions,
sum(calc.reduction_qty)                         as deductions,
sum(calc.additon_qty) - sum(calc.reduction_qty) as final_stock,
p.name as product_name,
ph.name as pharmacy_name,
TO_CHAR(
  TO_DATE(extract(month from calc.change_at)::text, 'MM'), 'Month'
)                                              AS month`
	newQuery = newQuery +
		` from calc
JOIN pharmacy_products as pp ON pp.id = calc.pharmacy_product_id
JOIN public.products p on p.id = pp.product_id
JOIN pharmacies as ph on pp.pharmacy_id = ph.id`
	rawQuery = rawQuery + ` select count( Distinct  calc.pharmacy_product_id) as total_item`
	rawQuery = rawQuery + ` from calc
	JOIN pharmacy_products as pp ON pp.id = calc.pharmacy_product_id
	JOIN public.products p on p.id = pp.product_id
	JOIN pharmacies as ph on pp.pharmacy_id = ph.id`
	where := fmt.Sprintf(" Where ph.admin_id = %v", ctx.Value("user_id").(uint))
	name := query.GetConditionValue("product_name")
	month := query.GetConditionValue("month")
	sortBy := query.GetOrder()
	if name != nil {
		temp := "%" + name.(string) + "%"
		where = where + fmt.Sprintf(" and p.name ILIKE  '%v'", temp)
	}
	if month != nil {
		where = where + fmt.Sprintf(" and extract(month from calc.change_at) = %v", month)
	}
	where2 := where + " group by calc.pharmacy_product_id, p.name, ph.name, month"
	sortString := strings.Split(query.GetOrder(), " ")
	if sortBy != "" {
		switch sortString[0] {
		case "id":
			sortBy = "calc.pharmacy_product_id " + sortString[1]

		}
		where2 = where2 + fmt.Sprintf(" ORDER BY %v", sortBy)
	}
	newQuery = newQuery + where2
	rawQuery = rawQuery + where + ` group by calc.pharmacy_product_id`

	page := query.GetPage()
	limit := query.GetLimit()
	if limit != nil && page > 0 {
		offset := (page - 1) * *limit
		newQuery = newQuery + fmt.Sprintf(" Offset %v", offset)
	}
	if limit != nil {
		newQuery = newQuery + fmt.Sprintf(" Limit %v", *limit)
	}
	err := r.db.WithContext(ctx).Raw(newQuery).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	err = r.db.WithContext(ctx).Raw(rawQuery).Scan(&totalItem).Error
	if err != nil {
		return nil, err
	}
	totalPage := 0
	if limit == nil {
		totalPage = 1
	} else {
		div := int(math.Min(float64(totalItem), float64(*limit)))
		if div == 0 {
			div = 1
		}
		totalPage = int(totalItem) / div
		if int(totalItem)%div != 0 {
			totalPage++
		}
	}
	currentPage := int(math.Min(float64(query.GetPage()), float64(totalPage)))
	pageResult := valueobject.PagedResult{Data: result, CurrentItems: len(result),
		CurrentPage: currentPage, TotalPage: totalPage, TotalItem: int(totalItem)}
	return &pageResult, nil
}

func (r *stockRecordRepository) BulkCreate(ctx context.Context, records []*entity.StockRecord) error {
	return r.conn(ctx).Model(&entity.StockRecord{}).Create(records).Error
}
