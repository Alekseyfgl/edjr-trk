package service

import (
	"context"
	"edjr-trk/internal/model"
	"edjr-trk/internal/repository"
	"edjr-trk/pkg/utils"
	"errors"
	"go.uber.org/zap"
	"time"
)

type AuthServiceInterface interface {
	Login(ctx context.Context, email, password string) (*model.LoginResponse, error)
}

type authService struct {
	repo       repository.UserRepositoryInterface
	jwtService JWTServiceInterface
	logger     *zap.Logger
}

// NewAuthService - создаёт новый экземпляр AuthService
func NewAuthService(repo repository.UserRepositoryInterface, jwtService JWTServiceInterface, logger *zap.Logger) AuthServiceInterface {
	return &authService{repo: repo, jwtService: jwtService, logger: logger}
}

func (s *authService) Login(ctx context.Context, email, password string) (*model.LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Failed to fetch user by email", zap.Error(err))
		return nil, err
	}

	// Проверяем пароль
	isValidPassword, err := utils.CompareHashes(password, user.Password)
	if err != nil {
		s.logger.Error("Failed to compare password hashes", zap.Error(err))
		return nil, err
	}

	if !isValidPassword {
		s.logger.Warn("Invalid password provided")
		return nil, errors.New("invalid email or password")
	}

	userIdStr := user.ID.Hex()
	// Генерируем access токен
	accessToken, err := s.jwtService.GenerateAccessToken(userIdStr, 24*time.Hour)
	if err != nil {
		s.logger.Error("Failed to generate access token", zap.Error(err))
		return nil, err
	}

	result := model.LoginResponse{AccessToken: accessToken}
	// Логируем успешный вход
	s.logger.Info("User logged in successfully", zap.String("userId", userIdStr))
	return &result, nil
}
