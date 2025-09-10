package services

import (
	"context"
	"logistics-api/internal/core/domain"
)

type CoordinateService interface {
	ValidateCoordinates(ctx context.Context, coords domain.Coordinates) error
	GetAddressFromCoordinates(ctx context.Context, coords domain.Coordinates) (*domain.Address, error)
	GetDistanceBetweenPoints(ctx context.Context, origin, destination domain.Coordinates) (float64, error)
}
