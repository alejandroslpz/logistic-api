package order

import (
	"context"
	"logistics-api/internal/core/domain"
	"logistics-api/internal/core/ports/repositories"
	"logistics-api/internal/core/usecases/dto"
	appErrors "logistics-api/internal/pkg/errors"
	"logistics-api/internal/pkg/logger"
)

type UpdateOrderStatusUseCase struct {
	orderRepo repositories.OrderRepository
	logger    logger.Logger
}

func NewUpdateOrderStatusUseCase(
	orderRepo repositories.OrderRepository,
	logger logger.Logger,
) *UpdateOrderStatusUseCase {
	return &UpdateOrderStatusUseCase{
		orderRepo: orderRepo,
		logger:    logger,
	}
}

func (uc *UpdateOrderStatusUseCase) Execute(ctx context.Context, orderID string, userRole domain.UserRole, req dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error) {
	uc.logger.Info("Updating order status",
		logger.String("order_id", orderID),
		logger.String("new_status", string(req.Status)),
		logger.String("user_role", string(userRole)),
	)

	if userRole != domain.AdminRole {
		uc.logger.Warn("Unauthorized attempt to update order status",
			logger.String("order_id", orderID),
			logger.String("user_role", string(userRole)),
		)
		return nil, appErrors.NewForbiddenError()
	}

	order, err := uc.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		uc.logger.Error("Order not found", logger.String("order_id", orderID))
		return nil, appErrors.NewNotFoundError("order")
	}

	if err := order.UpdateStatus(req.Status); err != nil {
		uc.logger.Warn("Invalid status transition",
			logger.String("order_id", orderID),
			logger.String("current_status", string(order.Status)),
			logger.String("new_status", string(req.Status)),
			logger.Error(err),
		)
		return nil, appErrors.NewValidationError(err.Error())
	}

	if err := uc.orderRepo.Update(ctx, order); err != nil {
		uc.logger.Error("Failed to update order", logger.Error(err))
		return nil, appErrors.NewInternalError()
	}

	uc.logger.Info("Order status updated successfully",
		logger.String("order_id", orderID),
		logger.String("new_status", string(req.Status)),
	)

	return dto.ToOrderResponse(order), nil
}
