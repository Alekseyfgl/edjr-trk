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

type ProductRepositoryInterface interface {
	CreateProduct(ctx context.Context, article model.RowProduct) (model.RowProduct, error)
	PatchProductById(ctx context.Context, dto *dto.PatchProductRequest, id string) (*model.RowProduct, error)
	GetProductById(ctx context.Context, id string) (*model.RowProduct, error)
	GetAllProducts(ctx context.Context, pageNumber, pageSize int) ([]model.RowProduct, int, error)
	RemoveProductById(ctx context.Context, id string) error
}

type productRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewProductRepository(client *mongo.Client, logger *zap.Logger) ProductRepositoryInterface {
	return &productRepository{
		collection: client.Database(env.GetEnv("MONGO_DB_NAME", "")).Collection("products"),
		logger:     logger,
	}
}

func (r *productRepository) GetAllProducts(ctx context.Context, pageNumber, pageSize int) ([]model.RowProduct, int, error) {
	if pageNumber < 1 {
		pageNumber = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	skip := utils.CalculateOffset(pageNumber, pageSize)

	//common document count
	totalCount64, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		r.logger.Error("Failed to count articles", zap.Error(err))
		return nil, 0, err
	}
	totalCount := int(totalCount64)
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

	products, decodeErr := utils.DecodeCursor[model.RowProduct](ctx, cursor, r.logger)
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
		zap.Int("totalCount", totalCount),
		zap.Int("fetchedItems", len(products)),
	)

	return products, totalCount, nil
}

func (r *productRepository) CreateProduct(ctx context.Context, product model.RowProduct) (model.RowProduct, error) {
	result, err := r.collection.InsertOne(ctx, product)
	if err != nil {
		r.logger.Error("Failed to insert product", zap.Error(err))
		return model.RowProduct{}, err
	}

	// Checking and converting InsertedID to ObjectID
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		r.logger.Error("Inserted ID is not of type ObjectID")
		return model.RowProduct{}, err
	}

	product.ID = insertedID
	r.logger.Info("Product created successfully", zap.String("id", product.ID.Hex()))

	return product, nil
}

func (r *productRepository) GetProductById(ctx context.Context, id string) (*model.RowProduct, error) {
	var product model.RowProduct

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid ID format", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	// Поиск статьи по ID в коллекции MongoDB.
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warn("Product not found", zap.String("id", id))
			return nil, mongo.ErrNoDocuments
		}
		r.logger.Error("Failed to query database", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	r.logger.Info("Product found", zap.String("id", product.ID.Hex()))
	return &product, nil
}

func (r *productRepository) PatchProductById(ctx context.Context, dto *dto.PatchProductRequest, id string) (*model.RowProduct, error) {
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
	if dto.ShortText != nil {
		update["shortText"] = *dto.ShortText
	}
	if dto.Img != nil {
		update["img"] = *dto.Img
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
		r.logger.Error("Failed to update product", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	if updateResult.MatchedCount == 0 {
		r.logger.Warn("Product not found for update", zap.String("id", id))
		return nil, mongo.ErrNoDocuments
	}

	r.logger.Info("Product updated successfully", zap.String("id", id))

	// Получение обновленной статьи
	updatedProduct, err := r.GetProductById(ctx, id)
	if err != nil {
		r.logger.Error("Failed to retrieve updated product", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return updatedProduct, nil
}

func (r *productRepository) RemoveProductById(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid ID format", zap.String("id", id), zap.Error(err))
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		r.logger.Error("Failed to delete product", zap.String("id", id), zap.Error(err))
		return err
	}

	r.logger.Info("Product successfully deleted", zap.String("id", id))
	return nil
}
