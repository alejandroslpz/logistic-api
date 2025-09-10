package dto

import (
	"logistics-api/internal/core/domain"
	"time"
)

type CreateOrderRequest struct {
	OriginCoordinates      domain.Coordinates `json:"origin_coordinates" validate:"required"`
	DestinationCoordinates domain.Coordinates `json:"destination_coordinates" validate:"required"`
	OriginAddress          domain.Address     `json:"origin_address" validate:"required"`
	DestinationAddress     domain.Address     `json:"destination_address" validate:"required"`
	ProductQuantity        int                `json:"product_quantity" validate:"required,min=1"`
	TotalWeight            float64            `json:"total_weight" validate:"required,min=0.1"`
}

type UpdateOrderStatusRequest struct {
	Status domain.OrderStatus `json:"status" validate:"required,oneof=creado recolectado en_estacion en_ruta entregado cancelado"`
}

type OrderResponse struct {
	ID                     string             `json:"id"`
	ClientID               string             `json:"client_id"`
	OriginCoordinates      domain.Coordinates `json:"origin_coordinates"`
	DestinationCoordinates domain.Coordinates `json:"destination_coordinates"`
	OriginAddress          domain.Address     `json:"origin_address"`
	DestinationAddress     domain.Address     `json:"destination_address"`
	ProductQuantity        int                `json:"product_quantity"`
	TotalWeight            float64            `json:"total_weight"`
	PackageSize            domain.PackageSize `json:"package_size"`
	Status                 domain.OrderStatus `json:"status"`
	CreatedAt              string             `json:"created_at"`
	UpdatedAt              string             `json:"updated_at"`
	Client                 *UserResponse      `json:"client,omitempty"`
}

type ListOrdersRequest struct {
	Page   int                `json:"page" validate:"min=1"`
	Limit  int                `json:"limit" validate:"min=1,max=100"`
	Status domain.OrderStatus `json:"status,omitempty"`
}

type ListOrdersResponse struct {
	Orders     []*OrderResponse `json:"orders"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

func ToOrderResponse(order *domain.Order) *OrderResponse {
	response := &OrderResponse{
		ID:                     order.ID,
		ClientID:               order.ClientID,
		OriginCoordinates:      order.OriginCoords,
		DestinationCoordinates: order.DestinationCoords,
		OriginAddress:          order.OriginAddress,
		DestinationAddress:     order.DestinationAddress,
		ProductQuantity:        order.ProductQuantity,
		TotalWeight:            order.TotalWeight,
		PackageSize:            order.PackageSize,
		Status:                 order.Status,
		CreatedAt:              order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:              order.UpdatedAt.Format(time.RFC3339),
	}

	if order.Client.ID != "" {
		response.Client = ToUserResponse(&order.Client)
	}

	return response
}

func ToOrderResponseList(orders []*domain.Order) []*OrderResponse {
	responses := make([]*OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = ToOrderResponse(order)
	}
	return responses
}
