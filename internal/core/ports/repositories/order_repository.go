package repositories

import (
	"context"
	"logistics-api/internal/core/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id string) (*domain.Order, error)
	GetByClientID(ctx context.Context, clientID string, limit, offset int) ([]*domain.Order, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	Delete(ctx context.Context, id string) error
	GetByStatus(ctx context.Context, status domain.OrderStatus, limit, offset int) ([]*domain.Order, error)
	UpdateStatus(ctx context.Context, orderID string, status domain.OrderStatus) error
	CountByClientID(ctx context.Context, clientID string) (int64, error)
	CountTotal(ctx context.Context) (int64, error)
	CountByStatus(ctx context.Context, status domain.OrderStatus) (int64, error)
}
