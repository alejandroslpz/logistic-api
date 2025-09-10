package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string
type PackageSize string

const (
	StatusCreated   OrderStatus = "creado"
	StatusCollected OrderStatus = "recolectado"
	StatusAtStation OrderStatus = "en_estacion"
	StatusInRoute   OrderStatus = "en_ruta"
	StatusDelivered OrderStatus = "entregado"
	StatusCancelled OrderStatus = "cancelado"

	PackageSizeS       PackageSize = "S"
	PackageSizeM       PackageSize = "M"
	PackageSizeL       PackageSize = "L"
	PackageSizeSpecial PackageSize = "SPECIAL"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90" gorm:"not null"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180" gorm:"not null"`
}

type Order struct {
	ID                 string      `json:"id" gorm:"primaryKey"`
	ClientID           string      `json:"client_id" gorm:"not null;index"`
	OriginCoords       Coordinates `json:"origin_coordinates" gorm:"embedded;embeddedPrefix:origin_"`
	DestinationCoords  Coordinates `json:"destination_coordinates" gorm:"embedded;embeddedPrefix:destination_"`
	OriginAddress      Address     `json:"origin_address" gorm:"embedded;embeddedPrefix:origin_addr_"`
	DestinationAddress Address     `json:"destination_address" gorm:"embedded;embeddedPrefix:dest_addr_"`
	ProductQuantity    int         `json:"product_quantity" gorm:"not null" validate:"required,min=1"`
	TotalWeight        float64     `json:"total_weight" gorm:"not null" validate:"required,min=0.1"`
	PackageSize        PackageSize `json:"package_size" gorm:"not null"`
	Status             OrderStatus `json:"status" gorm:"not null;default:'creado'"`
	CreatedAt          time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time   `json:"updated_at" gorm:"autoUpdateTime"`

	Client User `json:"client,omitempty" gorm:"foreignKey:ClientID"`
}

func NewOrder(
	clientID string,
	originCoords, destCoords Coordinates,
	originAddr, destAddr Address,
	productQuantity int,
	totalWeight float64,
) (*Order, error) {

	if err := validateCoordinates(originCoords); err != nil {
		return nil, errors.New("invalid origin coordinates: " + err.Error())
	}

	if err := validateCoordinates(destCoords); err != nil {
		return nil, errors.New("invalid destination coordinates: " + err.Error())
	}

	if productQuantity <= 0 {
		return nil, errors.New("product quantity must be greater than 0")
	}

	if totalWeight <= 0 {
		return nil, errors.New("total weight must be greater than 0")
	}

	packageSize := determinePackageSize(totalWeight)
	if packageSize == PackageSizeSpecial {
		return nil, errors.New("weight exceeds standard service limit. Please contact us for special arrangements")
	}

	return &Order{
		ID:                 uuid.New().String(),
		ClientID:           clientID,
		OriginCoords:       originCoords,
		DestinationCoords:  destCoords,
		OriginAddress:      originAddr,
		DestinationAddress: destAddr,
		ProductQuantity:    productQuantity,
		TotalWeight:        totalWeight,
		PackageSize:        packageSize,
		Status:             StatusCreated,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}, nil
}

func (o *Order) UpdateStatus(newStatus OrderStatus) error {
	validTransitions := map[OrderStatus][]OrderStatus{
		StatusCreated:   {StatusCollected, StatusCancelled},
		StatusCollected: {StatusAtStation, StatusCancelled},
		StatusAtStation: {StatusInRoute, StatusCancelled},
		StatusInRoute:   {StatusDelivered, StatusCancelled},
		StatusDelivered: {},
		StatusCancelled: {},
	}

	allowedStatuses, exists := validTransitions[o.Status]
	if !exists {
		return errors.New("invalid current status")
	}

	for _, allowed := range allowedStatuses {
		if allowed == newStatus {
			o.Status = newStatus
			o.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("invalid status transition from " + string(o.Status) + " to " + string(newStatus))
}

func (o *Order) CanBeModifiedBy(userRole UserRole) bool {
	return userRole == AdminRole
}

func validateCoordinates(coords Coordinates) error {
	if coords.Latitude < -90 || coords.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if coords.Longitude < -180 || coords.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	return nil
}

func determinePackageSize(weight float64) PackageSize {
	switch {
	case weight <= 5:
		return PackageSizeS
	case weight <= 15:
		return PackageSizeM
	case weight <= 25:
		return PackageSizeL
	default:
		return PackageSizeSpecial
	}
}

func GetValidStatuses() []OrderStatus {
	return []OrderStatus{
		StatusCreated, StatusCollected, StatusAtStation,
		StatusInRoute, StatusDelivered, StatusCancelled,
	}
}

func GetValidPackageSizes() []PackageSize {
	return []PackageSize{
		PackageSizeS, PackageSizeM, PackageSizeL,
	}
}
