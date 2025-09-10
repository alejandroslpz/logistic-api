package handlers

import (
	"net/http"
	"strconv"

	httpDto "logistics-api/internal/adapters/primary/http/dto"
	"logistics-api/internal/core/domain"
	"logistics-api/internal/core/usecases/dto"
	"logistics-api/internal/core/usecases/order"
	appErrors "logistics-api/internal/pkg/errors"
	"logistics-api/internal/pkg/logger"
	"logistics-api/internal/pkg/validator"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	createOrderUC  *order.CreateOrderUseCase
	getOrdersUC    *order.GetOrdersUseCase
	updateStatusUC *order.UpdateOrderStatusUseCase
	validator      *validator.Validator
	logger         logger.Logger
}

func NewOrderHandler(
	createOrderUC *order.CreateOrderUseCase,
	getOrdersUC *order.GetOrdersUseCase,
	updateStatusUC *order.UpdateOrderStatusUseCase,
	validator *validator.Validator,
	logger logger.Logger,
) *OrderHandler {
	return &OrderHandler{
		createOrderUC:  createOrderUC,
		getOrdersUC:    getOrdersUC,
		updateStatusUC: updateStatusUC,
		validator:      validator,
		logger:         logger,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request format", logger.Error(err))
		httpDto.ValidationErrorResponse(c, "Invalid request format")
		return
	}

	if err := h.validator.Validate(req); err != nil {
		h.logger.Warn("Validation failed", logger.Error(err))
		httpDto.ValidationErrorResponse(c, err.Error())
		return
	}

	clientID := c.GetString("user_id")
	if clientID == "" {
		h.logger.Error("Client ID not found in context")
		httpDto.UnauthorizedResponse(c)
		return
	}

	response, err := h.createOrderUC.Execute(c.Request.Context(), clientID, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	httpDto.SuccessResponse(c, http.StatusCreated, "Order created successfully", response)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	listReq := dto.ListOrdersRequest{
		Page:   page,
		Limit:  limit,
		Status: domain.OrderStatus(status),
	}

	if err := h.validator.Validate(listReq); err != nil {
		httpDto.ValidationErrorResponse(c, err.Error())
		return
	}

	userRole, _ := c.Get("user_role")
	role := userRole.(domain.UserRole)
	var response *dto.ListOrdersResponse
	var err error

	if role == domain.AdminRole {
		response, err = h.getOrdersUC.ExecuteForAdmin(c.Request.Context(), listReq)
	} else {
		clientID := c.GetString("user_id")
		response, err = h.getOrdersUC.ExecuteForClient(c.Request.Context(), clientID, listReq)
	}

	if err != nil {
		h.handleError(c, err)
		return
	}

	meta := &httpDto.PaginationMeta{
		Total:      response.Total,
		Page:       response.Page,
		Limit:      response.Limit,
		TotalPages: response.TotalPages,
	}

	httpDto.PaginatedSuccessResponse(c, response.Orders, meta)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		httpDto.ValidationErrorResponse(c, "Order ID is required")
		return
	}

	httpDto.ErrorResponse(c, http.StatusNotImplemented, "not_implemented", "Feature not implemented yet")
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		httpDto.ValidationErrorResponse(c, "Order ID is required")
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request format", logger.Error(err))
		httpDto.ValidationErrorResponse(c, "Invalid request format")
		return
	}

	if err := h.validator.Validate(req); err != nil {
		h.logger.Warn("Validation failed", logger.Error(err))
		httpDto.ValidationErrorResponse(c, err.Error())
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		httpDto.UnauthorizedResponse(c)
		return
	}

	role := userRole.(domain.UserRole)
	response, err := h.updateStatusUC.Execute(c.Request.Context(), orderID, role, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	httpDto.SuccessResponse(c, http.StatusOK, "Order status updated successfully", response)
}

func (h *OrderHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*appErrors.AppError); ok {
		httpDto.ErrorResponse(c, appErr.Code, appErr.Type, appErr.Message)
		return
	}

	h.logger.Error("Unexpected error", logger.Error(err))
	httpDto.InternalErrorResponse(c)
}
