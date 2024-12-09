package service

import (
	"context"
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/model"
	"edjr-trk/internal/repository"
	"edjr-trk/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"time"
)

type UserServiceInterface interface {
	CreateNewAdmin(ctx context.Context, dto *dto.CreateUserRequest) (*model.UserResponse, error)
	RemoveUserById(ctx context.Context, id string) (string, error)
	GetAllUsers(ctx context.Context, pageNumber, pageSize int) (*model.Paginate[*model.UserResponse], error)
}

type userService struct {
	repo   repository.UserRepositoryInterface
	logger *zap.Logger
}

// NewUserService - создаёт новый экземпляр UserService.
func NewUserService(repo repository.UserRepositoryInterface, logger *zap.Logger) UserServiceInterface {
	return &userService{repo: repo, logger: logger}
}

// CreateNewAdmin - creat admin as user-user
func (s *userService) CreateNewAdmin(ctx context.Context, dto *dto.CreateUserRequest) (*model.UserResponse, error) {

	newAdmin := model.RowUser{
		ID:        primitive.NewObjectID(),
		Email:     dto.Email,
		Phone:     dto.Phone,
		IsAdmin:   true,
		Password:  dto.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := newAdmin.HashPassword()
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, err
	}

	createdAdmin, err := s.repo.CreateNewAdmin(ctx, &newAdmin)
	if err != nil {
		s.logger.Error("Failed to save admin", zap.Error(err))
		return nil, err
	}

	transformedResp := createdAdmin.CreateUserResp()
	s.logger.Info("Admin created successfully", zap.String("id", transformedResp.ID.Hex()))
	return transformedResp, nil
}

// RemoveUserById - remove user by id
func (s *userService) RemoveUserById(ctx context.Context, id string) (string, error) {
	err := s.repo.RemoveUserById(ctx, id)

	if err != nil {
		s.logger.Error("Failed to remove user", zap.Error(err))
		return "", err
	}

	return id, nil
}

func (s *userService) GetAllUsers(ctx context.Context, pageNumber, pageSize int) (*model.Paginate[*model.UserResponse], error) {
	users, totalCount, err := s.repo.GetAll(ctx, pageNumber, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch all users", zap.Error(err))
		return nil, err
	}

	// Преобразуем RowArticle в ArticleResponse
	transformedResp := make([]*model.UserResponse, len(users))
	for i, user := range users {
		transformedResp[i] = user.CreateUserResp()
	}

	// Формируем структуру Paginate с типом ArticleResponse
	result := &model.Paginate[*model.UserResponse]{
		PageNumber:     pageNumber,
		RowTotalCount:  int(totalCount),
		TotalPageCount: utils.CalculateTotalPages(totalCount, pageSize),
		PageSize:       pageSize,
		Items:          transformedResp,
	}

	s.logger.Info("All users fetched successfully with pagination",
		zap.Int("pageNumber", pageNumber),
		zap.Int("pageSize", pageSize),
		zap.Int("totalCount", int(totalCount)),
		zap.Int("fetchedItems", len(transformedResp)),
	)

	return result, nil
}
