package repository

import (
	"context"

	"github.com/night1010/everhealth/dto"
	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
)

type OrderItemRepository interface {
	BaseRepository[entity.OrderItem]
	BulkCreate(context.Context, []*entity.OrderItem) error
	MonthlyReport(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	MonthlyReportAdminPharmacy(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error)
	ListOfOrderItem(ctx context.Context, orderId uint, userId uint) ([]*entity.OrderItem, error)
}

type orderItemRepository struct {
	*baseRepository[entity.OrderItem]
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) OrderItemRepository {
	return &orderItemRepository{
		db:             db,
		baseRepository: &baseRepository[entity.OrderItem]{db: db},
	}
}

func (r *orderItemRepository) BulkCreate(ctx context.Context, items []*entity.OrderItem) error {
	return r.conn(ctx).Model(&entity.OrderItem{}).Create(items).Error
}

func (r *orderItemRepository) MonthlyReport(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	result := dto.DataGraphReport{}
	pharmacy := query.GetConditionValue("pharmacy")
	product := query.GetConditionValue("product")
	productCategory := query.GetConditionValue("product_category")
	pharmacyGraph := []*dto.MonthlySalesReport{}
	productGraph := []*dto.MonthlySalesReport{}
	productCategoryGraph := []*dto.MonthlySalesReport{}
	err := r.db.WithContext(ctx).Raw(`select ph.name as pharmacy_name,sum(oi.quantity) as total_item,
	sum(po.total_payment) as total_sales,
	TO_CHAR(
			TO_DATE(extract(month from po.created_at)::text, 'MM'), 'Month'
	)       AS month
from order_items as oi
  JOIN pharmacy_products as pp ON pp.id = oi.pharmacy_product_id
  JOIN public.products p ON p.id = pp.product_id
  JOIN pharmacies as ph ON pp.pharmacy_id = ph.id
  JOIN product_orders as po ON po.id = oi.order_id
  JOIN public.product_categories pc on pc.id = p.product_category_id
WHERE ph.id = ?
group by extract(month from po.created_at), ph.name
ORDER BY extract(month from po.created_at)`, pharmacy).Scan(&pharmacyGraph).Error
	if err != nil {
		return nil, err
	}
	err = r.db.WithContext(ctx).Raw(`select p.name as product_name,sum(oi.quantity) as total_item,
	sum(po.total_payment) as total_sales,
	TO_CHAR(
			TO_DATE(extract(month from po.created_at)::text, 'MM'), 'Month'
	)       AS month
from order_items as oi
	  JOIN pharmacy_products as pp ON pp.id = oi.pharmacy_product_id
	  JOIN public.products p ON p.id = pp.product_id
	  JOIN pharmacies as ph ON pp.pharmacy_id = ph.id
	  JOIN product_orders as po ON po.id = oi.order_id
	  JOIN public.product_categories pc on pc.id = p.product_category_id
WHERE ph.id = ?
AND pp.id = ?
group by extract(month from po.created_at), p.name
ORDER BY extract(month from po.created_at)`, pharmacy, product).Scan(&productGraph).Error
	if err != nil {
		return nil, err
	}
	err = r.db.WithContext(ctx).Raw(`select pc.id as product_category_id,sum(oi.quantity) as total_item,
	sum(po.total_payment) as total_sales,
	TO_CHAR(
			TO_DATE(extract(month from po.created_at)::text, 'MM'), 'Month'
	)       AS month
from order_items as oi
	  JOIN pharmacy_products as pp ON pp.id = oi.pharmacy_product_id
	  JOIN public.products p ON p.id = pp.product_id
	  JOIN pharmacies as ph ON pp.pharmacy_id = ph.id
	  JOIN product_orders as po ON po.id = oi.order_id
	  JOIN public.product_categories pc on pc.id = p.product_category_id
WHERE ph.id = ?
AND  pc.id = ?
group by extract(month from po.created_at), pc.id
ORDER BY extract(month from po.created_at)`, pharmacy, productCategory).Scan(&productCategoryGraph).Error
	if err != nil {
		return nil, err
	}
	result.PharmacyGraph = pharmacyGraph
	result.ProductGraph = productGraph
	result.ProductCategoryGraph = productCategoryGraph
	return &valueobject.PagedResult{Data: result}, nil
}

func (r *orderItemRepository) MonthlyReportAdminPharmacy(ctx context.Context, query *valueobject.Query) (*valueobject.PagedResult, error) {
	result := dto.DataGraphReport{}
	pharmacy := query.GetConditionValue("pharmacy")
	product := query.GetConditionValue("product")
	productCategory := query.GetConditionValue("product_category")
	pharmacyGraph := []*dto.MonthlySalesReport{}
	productGraph := []*dto.MonthlySalesReport{}
	productCategoryGraph := []*dto.MonthlySalesReport{}
	err := r.db.WithContext(ctx).Raw(`select ph.name as pharmacy_name,sum(oi.quantity) as total_item,
	sum(po.total_payment) as total_sales,
	TO_CHAR(
			TO_DATE(extract(month from po.created_at)::text, 'MM'), 'Month'
	)       AS month
from order_items as oi
  JOIN pharmacy_products as pp ON pp.id = oi.pharmacy_product_id
  JOIN public.products p ON p.id = pp.product_id
  JOIN pharmacies as ph ON pp.pharmacy_id = ph.id
  JOIN product_orders as po ON po.id = oi.order_id
  JOIN public.product_categories pc on pc.id = p.product_category_id
WHERE ph.id = ?
AND ph.admin_id = ?
group by month, ph.name`, pharmacy, ctx.Value("user_id").(uint)).Scan(&pharmacyGraph).Error
	if err != nil {
		return nil, err
	}
	err = r.db.WithContext(ctx).Raw(`select p.name as product_name,sum(oi.quantity) as total_item,
	sum(po.total_payment) as total_sales,
	TO_CHAR(
			TO_DATE(extract(month from po.created_at)::text, 'MM'), 'Month'
	)       AS month
from order_items as oi
	  JOIN pharmacy_products as pp ON pp.id = oi.pharmacy_product_id
	  JOIN public.products p ON p.id = pp.product_id
	  JOIN pharmacies as ph ON pp.pharmacy_id = ph.id
	  JOIN product_orders as po ON po.id = oi.order_id
	  JOIN public.product_categories pc on pc.id = p.product_category_id
WHERE ph.id = ?
AND pp.id = ?
AND ph.admin_id = ?
group by month, p.name`, pharmacy, product, ctx.Value("user_id").(uint)).Scan(&productGraph).Error
	if err != nil {
		return nil, err
	}
	err = r.db.WithContext(ctx).Raw(`select pc.id as product_category_id,sum(oi.quantity) as total_item,
	sum(po.total_payment) as total_sales,
	TO_CHAR(
			TO_DATE(extract(month from po.created_at)::text, 'MM'), 'Month'
	)       AS month
from order_items as oi
	  JOIN pharmacy_products as pp ON pp.id = oi.pharmacy_product_id
	  JOIN public.products p ON p.id = pp.product_id
	  JOIN pharmacies as ph ON pp.pharmacy_id = ph.id
	  JOIN product_orders as po ON po.id = oi.order_id
	  JOIN public.product_categories pc on pc.id = p.product_category_id
WHERE ph.id = ?
AND  pc.id = ?
AND ph.admin_id = ?
group by month, pc.id`, pharmacy, productCategory, ctx.Value("user_id").(uint)).Scan(&productCategoryGraph).Error
	if err != nil {
		return nil, err
	}
	result.PharmacyGraph = pharmacyGraph
	result.ProductGraph = productGraph
	result.ProductCategoryGraph = productCategoryGraph
	return &valueobject.PagedResult{Data: result}, nil
}

func (r *orderItemRepository) ListOfOrderItem(ctx context.Context, orderId uint, userId uint) ([]*entity.OrderItem, error) {
	var orderItems []*entity.OrderItem
	err := r.conn(ctx).
		Model(&entity.OrderItem{}).
		Joins("JOIN pharmacy_products ON pharmacy_products.id = order_items.pharmacy_product_id ").
		Joins("JOIN pharmacies ON pharmacies.id = pharmacy_products.pharmacy_id").
		Joins("JOIN product_orders ON product_orders.id = order_items.order_id").
		Where("order_id = ?", orderId).
		Where("pharmacies.admin_id = ?", userId).
		Preload("PharmacyProduct").
		Preload("PharmacyProduct.Pharmacy").
		Find(&orderItems).Error
	if err != nil {
		return nil, err
	}
	return orderItems, nil
}
