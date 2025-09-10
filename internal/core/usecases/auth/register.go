package auth

import (
	"context"
	"logistics-api/internal/core/domain"
	"logistics-api/internal/core/ports/repositories"
	"logistics-api/internal/core/ports/services"
	"logistics-api/internal/core/usecases/dto"
	appErrors "logistics-api/internal/pkg/errors"
	"logistics-api/internal/pkg/logger"
)

type RegisterUseCase struct {
	userRepo    repositories.UserRepository
	authService services.AuthService
	logger      logger.Logger
}

func NewRegisterUseCase(
	userRepo repositories.UserRepository,
	authService services.AuthService,
	logger logger.Logger,
) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo:    userRepo,
		authService: authService,
		logger:      logger,
	}
}

func (uc *RegisterUseCase) Execute(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	uc.logger.Info("Attempting to register user", logger.String("email", req.Email))

	exists, err := uc.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		uc.logger.Error("Failed to check if user exists", logger.Error(err))
		return nil, appErrors.NewInternalError()
	}

	if exists {
		uc.logger.Warn("User registration failed - email already exists", logger.String("email", req.Email))
		return nil, appErrors.NewValidationError("email already exists")
	}

	user, err := domain.NewUser(req.Email, req.Password, req.Role)
	if err != nil {
		uc.logger.Error("Failed to create user domain entity", logger.Error(err))
		return nil, appErrors.NewValidationError(err.Error())
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.logger.Error("Failed to save user", logger.Error(err))
		return nil, appErrors.NewInternalError()
	}

	token, err := uc.authService.GenerateToken(user)
	if err != nil {
		uc.logger.Error("Failed to generate token", logger.Error(err))
		return nil, appErrors.NewInternalError()
	}

	uc.logger.Info("User registered successfully", logger.String("user_id", user.ID))

	return &dto.RegisterResponse{
		User:  dto.ToUserResponse(user),
		Token: token,
	}, nil
}
