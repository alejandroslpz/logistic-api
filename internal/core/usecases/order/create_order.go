package order

import (
	"context"
	"logistics-api/internal/core/domain"
	"logistics-api/internal/core/ports/repositories"
	"logistics-api/internal/core/ports/services"
	"logistics-api/internal/core/usecases/dto"
	appErrors "logistics-api/internal/pkg/errors"
	"logistics-api/internal/pkg/logger"
)

type CreateOrderUseCase struct {
	orderRepo    repositories.OrderRepository
	userRepo     repositories.UserRepository
	coordService services.CoordinateService
	logger       logger.Logger
}

func NewCreateOrderUseCase(
	orderRepo repositories.OrderRepository,
	userRepo repositories.UserRepository,
	coordService services.CoordinateService,
	logger logger.Logger,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:    orderRepo,
		userRepo:     userRepo,
		coordService: coordService,
		logger:       logger,
	}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, clientID string, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	uc.logger.Info("Creating new order", logger.String("client_id", clientID))

	client, err := uc.userRepo.GetByID(ctx, clientID)
	if err != nil {
		uc.logger.Error("Client not found", logger.String("client_id", clientID))
		return nil, appErrors.NewNotFoundError("client")
	}

	if err := uc.coordService.ValidateCoordinates(ctx, req.OriginCoordinates); err != nil {
		uc.logger.Warn("Invalid origin coordinates", logger.Error(err))
		return nil, appErrors.NewValidationError("invalid origin coordinates")
	}

	if err := uc.coordService.ValidateCoordinates(ctx, req.DestinationCoordinates); err != nil {
		uc.logger.Warn("Invalid destination coordinates", logger.Error(err))
		return nil, appErrors.NewValidationError("invalid destination coordinates")
	}

	order, err := domain.NewOrder(
		clientID,
		req.OriginCoordinates,
		req.DestinationCoordinates,
		req.OriginAddress,
		req.DestinationAddress,
		req.ProductQuantity,
		req.TotalWeight,
	)
	if err != nil {
		uc.logger.Error("Failed to create order entity", logger.Error(err))
		return nil, appErrors.NewValidationError(err.Error())
	}

	if err := uc.orderRepo.Create(ctx, order); err != nil {
		uc.logger.Error("Failed to save order", logger.Error(err))
		return nil, appErrors.NewInternalError()
	}

	order.Client = *client

	uc.logger.Info("Order created successfully",
		logger.String("order_id", order.ID),
		logger.String("client_id", clientID),
	)

	return dto.ToOrderResponse(order), nil
}
