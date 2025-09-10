package order

import (
	"context"
	"math"

	"logistics-api/internal/core/domain"
	"logistics-api/internal/core/ports/repositories"
	"logistics-api/internal/core/usecases/dto"
	appErrors "logistics-api/internal/pkg/errors"
	"logistics-api/internal/pkg/logger"
)

type GetOrdersUseCase struct {
	orderRepo repositories.OrderRepository
	userRepo  repositories.UserRepository
	logger    logger.Logger
}

func NewGetOrdersUseCase(
	orderRepo repositories.OrderRepository,
	userRepo repositories.UserRepository,
	logger logger.Logger,
) *GetOrdersUseCase {
	return &GetOrdersUseCase{
		orderRepo: orderRepo,
		userRepo:  userRepo,
		logger:    logger,
	}
}

func (uc *GetOrdersUseCase) ExecuteForClient(ctx context.Context, clientID string, req dto.ListOrdersRequest) (*dto.ListOrdersResponse, error) {
	uc.logger.Info("Getting orders for client", logger.String("client_id", clientID))

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	orders, err := uc.orderRepo.GetByClientID(ctx, clientID, req.Limit, offset)
	if err != nil {
		uc.logger.Error("Failed to get orders", logger.Error(err))
		return nil, appErrors.NewInternalError()
	}

	total, err := uc.orderRepo.CountByClientID(ctx, clientID)
	if err != nil {
		uc.logger.Error("Failed to count orders", logger.Error(err))
		return nil, appErrors.NewInternalError()
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &dto.ListOrdersResponse{
		Orders:     dto.ToOrderResponseList(orders),
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

func (uc *GetOrdersUseCase) ExecuteForAdmin(ctx context.Context, req dto.ListOrdersRequest) (*dto.ListOrdersResponse, error) {
	uc.logger.Info("Getting all orders for admin")

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}

	offset := (req.Page - 1) * req.Limit

	var orders []*domain.Order
	var total int64
	var err error

	if req.Status != "" {
		orders, err = uc.orderRepo.GetByStatus(ctx, req.Status, req.Limit, offset)
		if err != nil {
			uc.logger.Error("Failed to get orders by status", logger.Error(err))
			return nil, appErrors.NewInternalError()
		}
		total, err = uc.orderRepo.CountByStatus(ctx, req.Status)
	} else {
		orders, err = uc.orderRepo.GetAll(ctx, req.Limit, offset)
		if err != nil {
			uc.logger.Error("Failed to get all orders", logger.Error(err))
			return nil, appErrors.NewInternalError()
		}
		total, err = uc.orderRepo.CountTotal(ctx)
	}

	if err != nil {
		uc.logger.Error("Failed to count orders", logger.Error(err))
		return nil, appErrors.NewInternalError()
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	return &dto.ListOrdersResponse{
		Orders:     dto.ToOrderResponseList(orders),
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}
