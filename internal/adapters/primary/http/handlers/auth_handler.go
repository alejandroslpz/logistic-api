package handlers

import (
	"net/http"

	httpDto "logistics-api/internal/adapters/primary/http/dto"
	"logistics-api/internal/core/usecases/auth"
	"logistics-api/internal/core/usecases/dto"
	appErrors "logistics-api/internal/pkg/errors"
	"logistics-api/internal/pkg/logger"
	"logistics-api/internal/pkg/validator"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	registerUC *auth.RegisterUseCase
	loginUC    *auth.LoginUseCase
	validator  *validator.Validator
	logger     logger.Logger
}

func NewAuthHandler(
	registerUC *auth.RegisterUseCase,
	loginUC *auth.LoginUseCase,
	validator *validator.Validator,
	logger logger.Logger,
) *AuthHandler {
	return &AuthHandler{
		registerUC: registerUC,
		loginUC:    loginUC,
		validator:  validator,
		logger:     logger,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

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

	response, err := h.registerUC.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	httpDto.SuccessResponse(c, http.StatusCreated, "User registered successfully", response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

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

	response, err := h.loginUC.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	httpDto.SuccessResponse(c, http.StatusOK, "Login successful", response)
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	if appErr, ok := err.(*appErrors.AppError); ok {
		httpDto.ErrorResponse(c, appErr.Code, appErr.Type, appErr.Message)
		return
	}

	h.logger.Error("Unexpected error", logger.Error(err))
	httpDto.InternalErrorResponse(c)
}
