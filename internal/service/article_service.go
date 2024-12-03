package service

import (
	"context"
	"edjr-trk/internal/model"
	"edjr-trk/internal/repository"
	"go.uber.org/zap"
)

type ArticleService struct {
	repo   repository.ArticleRepositoryInterface // Интерфейс репозитория
	logger *zap.Logger
}

// ArticleServiceInterface - интерфейс для работы с сервисом статей.
type ArticleServiceInterface interface {
	//CreateArticle(ctx context.Context, article repository.RowArticle) (primitive.ObjectID, error)
	//GetArticleByID(ctx context.Context, id string) (*repository.RowArticle, error)
	GetAllArticles(ctx context.Context) ([]model.ArticleResponse, error)
}

// NewArticleService - создаёт новый экземпляр ArticleService.
func NewArticleService(repo repository.ArticleRepositoryInterface, logger *zap.Logger) *ArticleService {
	return &ArticleService{repo: repo, logger: logger}
}

// GetAllArticles - получает все статьи.
func (s *ArticleService) GetAllArticles(ctx context.Context) ([]model.ArticleResponse, error) {
	articles, err := s.repo.GetAll(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch all articles", zap.Error(err))
		return nil, err
	}

	// transform data []RowArticle to []ArticleResponse
	transformedResp := make([]model.ArticleResponse, len(articles))
	for i, article := range articles {
		transformedResp[i] = article.CreateArtResp()
	}

	s.logger.Info("All articles fetched successfully", zap.Int("count", len(articles)))
	return transformedResp, nil
}
