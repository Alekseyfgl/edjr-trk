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

type productService struct {
	repo   repository.ProductRepositoryInterface
	logger *zap.Logger
}

type ProductServiceInterface interface {
	CreateProduct(ctx context.Context, dto dto.CreateProductRequest) (*model.ProductResponse, error)
	RemoveProductById(ctx context.Context, id string) (string, error)
	PatchProductById(ctx context.Context, dto dto.PatchProductRequest, id string) (*model.ProductResponse, error)
	GetProductById(ctx context.Context, id string) (*model.ProductResponse, error)
	GetAllProducts(ctx context.Context, pageNumber, pageSize int) (*model.Paginate[*model.ProductResponse], error)
}

func NewProductService(repo repository.ProductRepositoryInterface, logger *zap.Logger) ProductServiceInterface {
	return &productService{repo: repo, logger: logger}
}

func (s *productService) GetAllProducts(ctx context.Context, pageNumber, pageSize int) (*model.Paginate[*model.ProductResponse], error) {
	products, totalCount, err := s.repo.GetAllProducts(ctx, pageNumber, pageSize)
	if err != nil {
		s.logger.Error("Failed to fetch all products", zap.Error(err))
		return nil, err
	}

	transformedResp := make([]*model.ProductResponse, len(products))
	for i, product := range products {
		transformedResp[i] = product.CreateProductResp()
	}

	result := &model.Paginate[*model.ProductResponse]{
		PageNumber:     pageNumber,
		RowTotalCount:  totalCount,
		TotalPageCount: utils.CalculateTotalPages(totalCount, pageSize),
		PageSize:       pageSize,
		Items:          transformedResp,
	}

	s.logger.Info("All products fetched successfully with pagination",
		zap.Int("pageNumber", pageNumber),
		zap.Int("pageSize", pageSize),
		zap.Int("totalCount", totalCount),
		zap.Int("fetchedItems", len(transformedResp)),
	)

	return result, nil
}

func (s *productService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*model.ProductResponse, error) {
	newArticle := model.RowProduct{
		ID:    primitive.NewObjectID(),
		Title: req.Title,
		Text:  req.Text,
		Img:   req.Img,
		Date:  time.Now(), // Используем primitive.DateTime для MongoDB
	}

	createdArticle, err := s.repo.CreateProduct(ctx, newArticle)
	if err != nil {
		s.logger.Error("Failed to save product", zap.Error(err))
		return nil, err
	}

	transformedResp := newArticle.CreateProductResp()

	s.logger.Info("Product created successfully", zap.String("id", createdArticle.ID.Hex()))
	return transformedResp, nil
}

func (s *productService) GetProductById(ctx context.Context, id string) (*model.ProductResponse, error) {
	product, err := s.repo.GetProductById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to save product", zap.Error(err))
		return nil, err
	}

	result := product.CreateProductResp()
	return result, err
}

func (s *productService) PatchProductById(ctx context.Context, dto dto.PatchProductRequest, id string) (*model.ProductResponse, error) {
	patchedProduct, err := s.repo.PatchProductById(ctx, &dto, id)

	if err != nil {
		s.logger.Error("Failed to save product", zap.Error(err))
		return nil, err
	}
	transformedResp := patchedProduct.CreateProductResp()
	return transformedResp, nil
}

func (s *productService) RemoveProductById(ctx context.Context, id string) (string, error) {
	err := s.repo.RemoveProductById(ctx, id)

	if err != nil {
		s.logger.Error("Failed to save product", zap.Error(err))
		return "", err
	}

	return id, nil
}
