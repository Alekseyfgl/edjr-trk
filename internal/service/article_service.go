package service

import (
	"context"
	"edjr-trk/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type ArticleService struct {
	repo   repository.ArticleRepositoryInterface // Интерфейс репозитория
	logger *zap.Logger
}

// ArticleServiceInterface - интерфейс для работы с сервисом статей.
type ArticleServiceInterface interface {
	CreateArticle(ctx context.Context, article repository.Article) (primitive.ObjectID, error)
	GetArticleByID(ctx context.Context, id string) (*repository.Article, error)
	GetAllArticles(ctx context.Context) ([]repository.Article, error)
}

// NewArticleService - создаёт новый экземпляр ArticleService.
func NewArticleService(repo repository.ArticleRepositoryInterface, logger *zap.Logger) *ArticleService {
	return &ArticleService{repo: repo, logger: logger}
}

// CreateArticle - создаёт новую статью.
func (s *ArticleService) CreateArticle(ctx context.Context, article repository.Article) (primitive.ObjectID, error) {
	id, err := s.repo.Create(ctx, article) // Вызов метода Create через интерфейс
	if err != nil {
		s.logger.Error("Failed to create article", zap.Error(err))
		return primitive.NilObjectID, err
	}

	s.logger.Info("Article created successfully", zap.String("id", id.Hex()))
	return id, nil
}

// GetArticleByID - получает статью по ID.
func (s *ArticleService) GetArticleByID(ctx context.Context, id string) (*repository.Article, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		s.logger.Error("Invalid article ID format", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	article, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		s.logger.Error("Failed to fetch article by ID", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	s.logger.Info("Article fetched successfully", zap.String("id", id))
	return article, nil
}

// GetAllArticles - получает все статьи.
func (s *ArticleService) GetAllArticles(ctx context.Context) ([]repository.Article, error) {
	articles, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch all articles", zap.Error(err))
		return nil, err
	}

	s.logger.Info("All articles fetched successfully", zap.Int("count", len(articles)))
	return articles, nil
}
