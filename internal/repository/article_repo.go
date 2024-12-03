package repository

import (
	"context"
	"edjr-trk/internal/model"
	"edjr-trk/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// ArticleRepositoryInterface - интерфейс для работы с коллекцией статей.
type ArticleRepositoryInterface interface {
	//Create(ctx context.Context, article RowArticle) (primitive.ObjectID, error)
	//GetByID(ctx context.Context, id primitive.ObjectID) (*RowArticle, error)
	GetAll(ctx context.Context, pageNumber, pageSize int) ([]model.RowArticle, int64, error)
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

// GetAll - get all articles with sort(desc) and pagination
func (r *articleRepository) GetAll(ctx context.Context, pageNumber, pageSize int) ([]model.RowArticle, int64, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	skip := utils.CalculateOffset(pageNumber, pageSize)

	//common document count
	totalCount, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error("Failed to count articles", zap.Error(err))
		return nil, 0, err
	}

	// Setting up search parameters with sorting
	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "date", Value: -1}}) // sort by desc

	// Retrieving data with pagination and sorting
	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		r.logger.Error("Failed to find articles", zap.Error(err))
		return nil, 0, err
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			r.logger.Warn("Failed to close cursor", zap.Error(closeErr))
		}
	}()

	articles, decodeErr := utils.DecodeCursor[model.RowArticle](ctx, cursor, r.logger)
	if decodeErr != nil {
		return nil, 0, decodeErr
	}

	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor encountered an error", zap.Error(err))
		return nil, 0, err
	}

	r.logger.Info("Articles fetched successfully with pagination and sorting",
		zap.Int("pageNumber", pageNumber),
		zap.Int("pageSize", pageSize),
		zap.Int64("totalCount", totalCount),
		zap.Int("fetchedItems", len(articles)),
	)

	return articles, totalCount, nil
}
