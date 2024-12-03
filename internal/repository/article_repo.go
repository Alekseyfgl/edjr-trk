package repository

import (
	"context"
	"edjr-trk/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// ArticleRepositoryInterface - интерфейс для работы с коллекцией статей.
type ArticleRepositoryInterface interface {
	//Create(ctx context.Context, article RowArticle) (primitive.ObjectID, error)
	//GetByID(ctx context.Context, id primitive.ObjectID) (*RowArticle, error)
	GetAll(ctx context.Context) ([]model.RowArticle, error)
}

// articleRepository - конкретная реализация интерфейса.
type articleRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewArticleRepository - создаёт новый экземпляр репозитория статей.
func NewArticleRepository(client *mongo.Client, logger *zap.Logger) ArticleRepositoryInterface {
	return &articleRepository{
		collection: client.Database("test").Collection("articles"),
		logger:     logger,
	}
}

// GetAll - получает все статьи из коллекции.
func (r *articleRepository) GetAll(ctx context.Context) ([]model.RowArticle, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		r.logger.Error("Failed to find articles", zap.Error(err))
		return nil, err
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			r.logger.Warn("Failed to close cursor", zap.Error(closeErr))
		}
	}()

	var articles []model.RowArticle
	for cursor.Next(ctx) {
		var article model.RowArticle
		if decodeErr := cursor.Decode(&article); decodeErr != nil {
			r.logger.Error("Failed to decode article", zap.Error(decodeErr))
			return nil, decodeErr
		}
		articles = append(articles, article)
	}

	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor encountered an error", zap.Error(err))
		return nil, err
	}

	r.logger.Info("Articles fetched successfully", zap.Int("count", len(articles)))
	return articles, nil
}
