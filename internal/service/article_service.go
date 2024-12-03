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
	GetAllArticles(ctx context.Context, pageNumber, pageSize int) (*model.Paginate[model.ArticleResponse], error)
}

// NewArticleService - создаёт новый экземпляр ArticleService.
func NewArticleService(repo repository.ArticleRepositoryInterface, logger *zap.Logger) *ArticleService {
	return &ArticleService{repo: repo, logger: logger}
}

// GetAllArticles - получает статьи с пагинацией.
func (s *ArticleService) GetAllArticles(ctx context.Context, pageNumber, pageSize int) (*model.Paginate[model.ArticleResponse], error) {
	articles, totalCount, err := s.repo.GetAll(ctx, pageNumber, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch all articles", zap.Error(err))
		return nil, err
	}

	// Преобразуем RowArticle в ArticleResponse
	transformedResp := make([]model.ArticleResponse, len(articles))
	for i, article := range articles {
		transformedResp[i] = article.CreateArtResp()
	}

	// Формируем структуру Paginate с типом ArticleResponse
	result := &model.Paginate[model.ArticleResponse]{
		PageNumber:    pageNumber,
		RowTotalCount: int(totalCount),
		CurrentPage:   pageNumber,
		PageSize:      pageSize,
		Items:         transformedResp,
	}

	s.logger.Info("All articles fetched successfully with pagination",
		zap.Int("pageNumber", pageNumber),
		zap.Int("pageSize", pageSize),
		zap.Int("totalCount", int(totalCount)),
		zap.Int("fetchedItems", len(transformedResp)),
	)

	return result, nil
}
