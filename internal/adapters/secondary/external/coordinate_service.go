package external

import (
	"context"
	"errors"
	"math"

	"logistics-api/internal/core/domain"
)

type CoordinateService struct {
	// Google Maps
}

func NewCoordinateService() *CoordinateService {
	return &CoordinateService{}
}

func (c *CoordinateService) ValidateCoordinates(ctx context.Context, coords domain.Coordinates) error {
	if coords.Latitude < -90 || coords.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if coords.Longitude < -180 || coords.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}

	return nil
}

func (c *CoordinateService) GetAddressFromCoordinates(ctx context.Context, coords domain.Coordinates) (*domain.Address, error) {
	return nil, errors.New("reverse geocoding not implemented")
}

func (c *CoordinateService) GetDistanceBetweenPoints(ctx context.Context, origin, destination domain.Coordinates) (float64, error) {
	return haversineDistance(origin, destination), nil
}

// Haversine distance
func haversineDistance(origin, destination domain.Coordinates) float64 {
	const earthRadius = 6371

	lat1Rad := toRadians(origin.Latitude)
	lon1Rad := toRadians(origin.Longitude)
	lat2Rad := toRadians(destination.Latitude)
	lon2Rad := toRadians(destination.Longitude)

	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func toRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
