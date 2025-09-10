package auth

import (
	"context"
	"logistics-api/internal/core/ports/repositories"
	"logistics-api/internal/core/ports/services"
	"logistics-api/internal/core/usecases/dto"
	appErrors "logistics-api/internal/pkg/errors"
	"logistics-api/internal/pkg/logger"
)

type LoginUseCase struct {
	userRepo    repositories.UserRepository
	authService services.AuthService
	logger      logger.Logger
}

func NewLoginUseCase(
	userRepo repositories.UserRepository,
	authService services.AuthService,
	logger logger.Logger,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:    userRepo,
		authService: authService,
		logger:      logger,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	uc.logger.Info("Attempting to login user", logger.String("email", req.Email))

	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		uc.logger.Warn("Login failed - user not found", logger.String("email", req.Email))
		return nil, appErrors.NewValidationError("invalid credentials")
	}

	if !user.ValidatePassword(req.Password) {
		uc.logger.Warn("Login failed - invalid password", logger.String("email", req.Email))
		return nil, appErrors.NewValidationError("invalid credentials")
	}

	token, err := uc.authService.GenerateToken(user)
	if err != nil {
		uc.logger.Error("Failed to generate token", logger.Error(err))
		return nil, appErrors.NewInternalError()
	}

	uc.logger.Info("User logged in successfully", logger.String("user_id", user.ID))

	return &dto.LoginResponse{
		User:  dto.ToUserResponse(user),
		Token: token,
	}, nil
}
