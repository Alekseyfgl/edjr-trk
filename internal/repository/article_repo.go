package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"time"
)

// Article - структура для хранения данных статьи.
type Article struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Text  string             `json:"text" bson:"text"`
	Title string             `json:"title" bson:"title"`
	Img   string             `json:"img,omitempty" bson:"img,omitempty"`
	Date  time.Time          `json:"date" bson:"date"`
}

// ArticleRepositoryInterface - интерфейс для работы с коллекцией статей.
type ArticleRepositoryInterface interface {
	Create(ctx context.Context, article Article) (primitive.ObjectID, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*Article, error)
	GetAll(ctx context.Context) ([]Article, error)
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

// Create - добавляет новую статью в коллекцию.
func (r *articleRepository) Create(ctx context.Context, article Article) (primitive.ObjectID, error) {
	article.ID = primitive.NewObjectID()
	article.Date = time.Now() // Устанавливаем текущую дату

	_, err := r.collection.InsertOne(ctx, article)
	if err != nil {
		r.logger.Error("Failed to insert article", zap.Error(err))
		return primitive.NilObjectID, err
	}
	r.logger.Info("Article created successfully", zap.String("id", article.ID.Hex()))
	return article.ID, nil
}

// GetByID - получает статью по ID.
func (r *articleRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*Article, error) {
	var article Article
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&article)
	if err != nil {
		r.logger.Error("Failed to find article by ID", zap.Error(err), zap.String("id", id.Hex()))
		return nil, err
	}
	r.logger.Info("Article fetched successfully", zap.String("id", id.Hex()))
	return &article, nil
}

// GetAll - получает все статьи из коллекции.
func (r *articleRepository) GetAll(ctx context.Context) ([]Article, error) {
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

	var articles []Article
	for cursor.Next(ctx) {
		var article Article
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
