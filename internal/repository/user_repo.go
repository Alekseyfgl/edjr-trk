package repository

import (
	"context"
	"edjr-trk/configs/env"
	configMongo "edjr-trk/configs/mongo"
	"edjr-trk/internal/model"
	"edjr-trk/pkg/utils"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// UserRepositoryInterface - интерфейс для работы с коллекцией пользователей.
type UserRepositoryInterface interface {
	CreateNewAdmin(ctx context.Context, user *model.RowUser) (*model.RowUser, error)
	RemoveUserById(ctx context.Context, id string) error
	GetAll(ctx context.Context, pageNumber, pageSize int) (*[]model.RowUser, int, error)
	GetUserByEmail(ctx context.Context, email string) (*model.RowUser, error)
	GetUserById(ctx context.Context, id string) (*model.RowUser, error)
}

// userRepository - конкретная реализация интерфейса.
type userRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewUserRepository(client *mongo.Client, logger *zap.Logger) UserRepositoryInterface {
	return &userRepository{
		collection: client.Database(env.GetEnv("MONGO_DB_NAME", "")).Collection(configMongo.UsersCollection),
		logger:     logger,
	}
}

// CreateNewAdmin - create admin as  super-user
func (r *userRepository) CreateNewAdmin(ctx context.Context, user *model.RowUser) (*model.RowUser, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		r.logger.Error("Failed to insert user", zap.Error(err))
		return nil, err
	}

	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		r.logger.Error("Inserted ID is not of type ObjectID")
		return nil, err
	}

	user.ID = insertedID
	r.logger.Info("User created successfully", zap.String("id", user.ID.Hex()))

	return user, nil
}

// RemoveUserById - remove user by id
func (r *userRepository) RemoveUserById(ctx context.Context, id string) error {
	// Преобразование строки в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid ID format", zap.String("id", id), zap.Error(err))
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		r.logger.Error("Failed to delete user", zap.String("id", id), zap.Error(err))
		return err
	}

	r.logger.Info("User successfully deleted", zap.String("id", id))
	return nil
}

// GetAll - get all articles with sort(desc) and pagination
func (r *userRepository) GetAll(ctx context.Context, pageNumber, pageSize int) (*[]model.RowUser, int, error) {
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
		r.logger.Error("Failed to count users", zap.Error(err))
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
		r.logger.Error("Failed to find users", zap.Error(err))
		return nil, 0, err
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			r.logger.Warn("Failed to close cursor", zap.Error(closeErr))
		}
	}()

	users, decodeErr := utils.DecodeCursor[model.RowUser](ctx, cursor, r.logger)
	if decodeErr != nil {
		return nil, 0, decodeErr
	}

	if err := cursor.Err(); err != nil {
		r.logger.Error("Cursor encountered an error", zap.Error(err))
		return nil, 0, err
	}

	r.logger.Info("Users fetched successfully with pagination and sorting",
		zap.Int("pageNumber", pageNumber),
		zap.Int("pageSize", pageSize),
		zap.Int("totalCount", totalCount),
		zap.Int("fetchedItems", len(users)),
	)

	return &users, totalCount, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*model.RowUser, error) {
	r.logger.Info("Start GetUserByEmail", zap.String("email", email))

	var user model.RowUser

	// Поиск пользователя по email
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warn("User not found", zap.String("email", email))
			return nil, mongo.ErrNoDocuments
		}
		r.logger.Error("Failed to query database", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	r.logger.Info("User found", zap.String("email", user.Email))
	return &user, nil
}

func (r *userRepository) GetUserById(ctx context.Context, id string) (*model.RowUser, error) {
	r.logger.Info("Start GetArticleById", zap.String("id", id))

	var user model.RowUser

	// Преобразование строки в ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Error("Invalid ID format", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	// Поиск статьи по ID в коллекции MongoDB.
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			r.logger.Warn("User not found", zap.String("id", id))
			return nil, mongo.ErrNoDocuments
		}
		r.logger.Error("Failed to query database", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	r.logger.Info("User found", zap.String("id", user.ID.Hex()))
	return &user, nil
}
