package repository

import (
	"context"
	"edjr-trk/configs/env"
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/model"
	"edjr-trk/pkg/utils"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// ArticleRepositoryInterface - интерфейс для работы с коллекцией статей.
type ArticleRepositoryInterface interface {
	Create(ctx context.Context, article model.RowArticle) (model.RowArticle, error)
	PatchArticleById(ctx context.Context, dto *dto.PatchArticleRequest, id string) (*model.RowArticle, error)
	GetArticleById(ctx context.Context, id string) (*model.RowArticle, error)
	GetAll(ctx context.Context, pageNumber, pageSize int) ([]model.RowArticle, int64, error)
	RemoveArticleById(ctx context.Context, id string) error
}

// articleRepository - конкретная реализация интерфейса.
type articleRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

// NewArticleRepository - создаёт новый экземпляр репозитория статей.
func NewArticleRepository(client *mongo.Client, logger *zap.Logger) ArticleRepositoryInterface {
	return &articleRepository{
		//collection: client.Database("test").Collection("articles"),
		collection: client.Database(env.GetEnv("MONGO_DB_NAME", "")).Collection("articles"),
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

// Create -create new article
func (r *articleRepository) Create(ctx context.Context, article model.RowArticle) (model.RowArticle, error) {
	// Вставка статьи в коллекцию.
	result, err := r.collection.InsertOne(ctx, article)
	if err != nil {
		r.logger.Error("Failed to insert article", zap.Error(err))
		return model.RowArticle{}, err
	}

	// Checking and converting InsertedID to ObjectID
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		r.logger.Error("Inserted ID is not of type ObjectID")
		return model.RowArticle{}, err
	}

	article.ID = insertedID
	r.logger.Info("Article created successfully", zap.String("id", article.ID.Hex()))

	return article, nil
}

// GetArticleById - находит статью по ObjectID.
func (r *articleRepository) GetArticleById(ctx context.Context, id string) (*model.RowArticle, error) {
	r.logger.Info("Start GetArticleById", zap.String("id", id))

	var article model.RowArticle

	// Преобразование строки в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid ID format", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	// Поиск статьи по ID в коллекции MongoDB.
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&article)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warn("Article not found", zap.String("id", id))
			return nil, mongo.ErrNoDocuments
		}
		r.logger.Error("Failed to query database", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	r.logger.Info("Article found", zap.String("id", article.ID.Hex()))
	return &article, nil
}

func (r *articleRepository) PatchArticleById(ctx context.Context, dto *dto.PatchArticleRequest, id string) (*model.RowArticle, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid ID format", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	// Формирование обновления
	update := bson.M{}
	if dto.Title != nil {
		update["title"] = *dto.Title
	}
	if dto.Text != nil {
		update["text"] = *dto.Text
	}
	if dto.Img != nil {
		update["img"] = *dto.Img
	} else {
		update["img"] = nil
	}

	if len(update) == 0 {
		r.logger.Warn("No fields to update", zap.String("id", id))
		return nil, errors.New("no fields to update")
	}

	// Обновление статьи в MongoDB
	updateResult, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		bson.M{"$set": update},
	)
	if err != nil {
		r.logger.Error("Failed to update article", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		r.logger.Warn("Article not found for update", zap.String("id", id))
		return nil, mongo.ErrNoDocuments
	}

	r.logger.Info("Article updated successfully", zap.String("id", id))

	// Получение обновленной статьи
	updatedArticle, err := r.GetArticleById(ctx, id)
	if err != nil {
		r.logger.Error("Failed to retrieve updated article", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return updatedArticle, nil
}

// RemoveArticleById - находит статью по ObjectID.
func (r *articleRepository) RemoveArticleById(ctx context.Context, id string) error {
	// Преобразование строки в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid ID format", zap.String("id", id), zap.Error(err))
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		r.logger.Error("Failed to delete article", zap.String("id", id), zap.Error(err))
		return err
	}

	r.logger.Info("Article successfully deleted", zap.String("id", id))
	return nil
}
