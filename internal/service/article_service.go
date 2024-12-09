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

type ArticleService struct {
	repo   repository.ArticleRepositoryInterface // Интерфейс репозитория
	logger *zap.Logger
}

// ArticleServiceInterface - интерфейс для работы с сервисом статей.
type ArticleServiceInterface interface {
	CreateArticle(ctx context.Context, dto dto.CreateArticleRequest) (*model.ArticleResponse, error)
	RemoveArticleById(ctx context.Context, id string) (string, error)
	PatchArticleById(ctx context.Context, dto dto.PatchArticleRequest, id string) (*model.ArticleResponse, error)
	GetArticleById(ctx context.Context, id string) (*model.ArticleResponse, error)
	GetAllArticles(ctx context.Context, pageNumber, pageSize int) (*model.Paginate[*model.ArticleResponse], error)
}

// NewArticleService - создаёт новый экземпляр ArticleService.
func NewArticleService(repo repository.ArticleRepositoryInterface, logger *zap.Logger) *ArticleService {
	return &ArticleService{repo: repo, logger: logger}
}

// GetAllArticles - получает статьи с пагинацией.
func (s *ArticleService) GetAllArticles(ctx context.Context, pageNumber, pageSize int) (*model.Paginate[*model.ArticleResponse], error) {
	articles, totalCount, err := s.repo.GetAll(ctx, pageNumber, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch all articles", zap.Error(err))
		return nil, err
	}

	// Преобразуем RowArticle в ArticleResponse
	transformedResp := make([]*model.ArticleResponse, len(articles))
	for i, article := range articles {
		transformedResp[i] = article.CreateArtResp()
	}

	// Формируем структуру Paginate с типом ArticleResponse
	result := &model.Paginate[*model.ArticleResponse]{
		PageNumber:     pageNumber,
		RowTotalCount:  totalCount,
		TotalPageCount: utils.CalculateTotalPages(totalCount, pageSize),
		PageSize:       pageSize,
		Items:          transformedResp,
	}

	s.logger.Info("All articles fetched successfully with pagination",
		zap.Int("pageNumber", pageNumber),
		zap.Int("pageSize", pageSize),
		zap.Int("totalCount", int(totalCount)),
		zap.Int("fetchedItems", len(transformedResp)),
	)

	return result, nil
}

// CreateArticle - создаёт новую статью.
func (s *ArticleService) CreateArticle(ctx context.Context, req dto.CreateArticleRequest) (*model.ArticleResponse, error) {
	// Создание новой статьи.
	newArticle := model.RowArticle{
		ID:    primitive.NewObjectID(),
		Title: req.Title,
		Text:  req.Text,
		Img:   req.Img,
		Date:  time.Now(), // Используем primitive.DateTime для MongoDB
	}

	// Сохранение статьи в репозитории.
	createdArticle, err := s.repo.Create(ctx, newArticle)
	if err != nil {
		s.logger.Error("Failed to save article", zap.Error(err))
		return nil, err
	}

	transformedResp := newArticle.CreateArtResp()

	s.logger.Info("Article created successfully", zap.String("id", createdArticle.ID.Hex()))
	return transformedResp, nil
}

func (s *ArticleService) GetArticleById(ctx context.Context, id string) (*model.ArticleResponse, error) {
	article, err := s.repo.GetArticleById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to save article", zap.Error(err))
		return nil, err
	}

	result := article.CreateArtResp()
	return result, err
}

// PatchArticleById - обновляет существующую статью частично.
func (s *ArticleService) PatchArticleById(ctx context.Context, dto dto.PatchArticleRequest, id string) (*model.ArticleResponse, error) {
	patchedArticle, err := s.repo.PatchArticleById(ctx, &dto, id)

	if err != nil {
		s.logger.Error("Failed to save article", zap.Error(err))
		return nil, err
	}
	transformedResp := patchedArticle.CreateArtResp()
	return transformedResp, nil
}

// RemoveArticleById - обновляет существующую статью частично.
func (s *ArticleService) RemoveArticleById(ctx context.Context, id string) (string, error) {
	err := s.repo.RemoveArticleById(ctx, id)

	if err != nil {
		s.logger.Error("Failed to save article", zap.Error(err))
		return "", err
	}

	return id, nil
}
