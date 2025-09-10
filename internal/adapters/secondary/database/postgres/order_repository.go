package postgres

import (
	"context"
	"errors"

	"logistics-api/internal/core/domain"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	err := r.db.WithContext(ctx).Preload("Client").Where("id = ?", id).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetByClientID(ctx context.Context, clientID string, limit, offset int) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.WithContext(ctx).
		Preload("Client").
		Where("client_id = ?", clientID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.WithContext(ctx).
		Preload("Client").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	return r.db.WithContext(ctx).Save(order).Error
}

func (r *OrderRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Order{}, "id = ?", id).Error
}

func (r *OrderRepository) GetByStatus(ctx context.Context, status domain.OrderStatus, limit, offset int) ([]*domain.Order, error) {
	var orders []*domain.Order
	err := r.db.WithContext(ctx).
		Preload("Client").
		Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID string, status domain.OrderStatus) error {
	return r.db.WithContext(ctx).
		Model(&domain.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}

func (r *OrderRepository) CountByClientID(ctx context.Context, clientID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Order{}).
		Where("client_id = ?", clientID).
		Count(&count).Error
	return count, err
}

func (r *OrderRepository) CountTotal(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.Order{}).Count(&count).Error
	return count, err
}

func (r *OrderRepository) CountByStatus(ctx context.Context, status domain.OrderStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Order{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}
