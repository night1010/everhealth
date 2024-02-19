package repository

import (
	"context"
	"strings"
	"time"

	"github.com/night1010/everhealth/entity"
	"github.com/night1010/everhealth/valueobject"
	"gorm.io/gorm"
)

type ProductOrderRepository interface {
	BaseRepository[entity.ProductOrder]
	CancelOrder(context.Context) error
	ConfirmOrder(context.Context) error
	FindAllOrders(context.Context, *valueobject.Query, uint, entity.RoleId) (*valueobject.PagedResult, error)
	FindOrderDetail(context.Context, uint, uint, entity.RoleId) (*entity.ProductOrder, error)
}

type productOrderRepository struct {
	*baseRepository[entity.ProductOrder]
	db *gorm.DB
}

func NewProductOrderRepository(db *gorm.DB) ProductOrderRepository {
	return &productOrderRepository{
		db:             db,
		baseRepository: &baseRepository[entity.ProductOrder]{db: db},
	}
}

func (r *productOrderRepository) CancelOrder(ctx context.Context) error {
	err := r.conn(ctx).
		Model(&entity.ProductOrder{}).
		Where("expired_at < ?", time.Now()).
		Update("order_status_id", 6).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *productOrderRepository) ConfirmOrder(ctx context.Context) error {
	err := r.conn(ctx).
		Model(&entity.ProductOrder{}).
		Where("ordered_at < ?", time.Now().Add(time.Hour*-7*24)).
		Where("order_status_id = ?", 4).
		Update("order_status_id", 5).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *productOrderRepository) FindAllOrders(ctx context.Context, query *valueobject.Query, userId uint, roleId entity.RoleId) (*valueobject.PagedResult, error) {
	return r.paginate(ctx, query, func(db *gorm.DB) *gorm.DB {
		orderStatus := query.GetConditionValue("order_status")
		name := query.GetConditionValue("name")
		switch strings.Split(query.GetOrder(), " ")[0] {
		case "price":
			query.WithSortBy("\"total_payment\"")
		case "order_date":
			query.WithSortBy("\"ordered_at\"")
		}
		db.Joins("LEFT JOIN order_items ON product_orders.id = order_items.order_id").
			Joins("LEFT JOIN pharmacy_products ON pharmacy_products.id = order_items.pharmacy_product_id").
			Joins("LEFT JOIN products ON pharmacy_products.product_id = products.id").
			Joins("LEFT JOIN order_statuses ON order_statuses.id = product_orders.order_status_id")

		if orderStatus != nil {
			db.Where("order_status_id = ?", orderStatus)
		}

		switch roleId {
		case entity.RoleUser:
			db.Where("profile_id = ?", userId)
			if name != nil {
				db.Where("products.name ILIKE ? ", name)
			}
		case entity.RoleAdmin:
			db.Joins("LEFT JOIN pharmacies ON pharmacies.id = pharmacy_products.pharmacy_id").
				Where("pharmacies.admin_id = ?", userId)
		}
		if roleId != entity.RoleUser {
			db.Joins("LEFT JOIN profiles ON profiles.user_id = product_orders.profile_id")
			if name != nil {
				db.Where("profiles.name ILIKE ? ", name)
			}
			db.Preload("Profile")
		}
		db.Group("product_orders.id")
		db.Preload("OrderItems.PharmacyProduct.Product")
		db.Preload("OrderStatus")
		return db
	})
}

func (r *productOrderRepository) FindOrderDetail(ctx context.Context, orderId, userId uint, roleId entity.RoleId) (*entity.ProductOrder, error) {
	var order entity.ProductOrder
	query := r.conn(ctx).
		Model(&entity.ProductOrder{}).
		Joins("LEFT JOIN order_items ON order_items.order_id = product_orders.id").
		Joins("LEFT JOIN pharmacy_products ON pharmacy_products.id = order_items.pharmacy_product_id").
		Joins("LEFT JOIN pharmacies ON pharmacies.id = pharmacy_products.pharmacy_id").
		Joins("LEFT JOIN admin_contacts ON admin_contacts.user_id = pharmacies.admin_id").
		Where("order_id = ? ", orderId).
		Preload("OrderItems.PharmacyProduct.Pharmacy").
		Preload("OrderStatus")
	switch roleId {
	case entity.RoleUser:
		query = query.Where("profile_id = ?", userId)
	case entity.RoleAdmin:
		query = query.Where("pharmacies.admin_id = ?", userId)
	case entity.RoleSuperAdmin:
		query = query.Preload("OrderItems.PharmacyProduct.Pharmacy.Admin.AdminContact.User")
	}
	err := query.Find(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil

}
