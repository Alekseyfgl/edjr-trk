package ioc

import (
	"edjr-trk/configs/mongo"
	"edjr-trk/internal/api/handlers"
	"edjr-trk/internal/repository"
	"edjr-trk/internal/service"
	"edjr-trk/pkg/log"
	mongodb "go.mongodb.org/mongo-driver/mongo" // Импортируем MongoDB драйвер правильно
	"go.uber.org/zap"
)

// Container - структура для хранения зависимостей.
type Container struct {
	Logger         *zap.Logger
	MongoClient    *mongodb.Client
	ArticleRepo    repository.ArticleRepositoryInterface
	ArticleService service.ArticleServiceInterface
	ArticleHandler *handlers.ArticleHandler
}

// NewContainer - создаем контейнер с зависимостями.
func NewContainer() *Container {
	// Initialize logger
	log.InitLogger()

	// Initialize MongoDB client (singleton)
	mongo.InitMongoSingleton()

	// Get MongoDB client
	clientDB := mongo.GetClient()

	// Get global logger
	logger := log.GetLogger()

	// Create dependencies
	articleRepo := repository.NewArticleRepository(clientDB, logger)
	articleService := service.NewArticleService(articleRepo, logger) // Передаём логгер
	articleHandler := handlers.NewArticleHandler(articleService, logger)

	// Return the container with all dependencies
	return &Container{
		Logger:         logger,
		MongoClient:    clientDB,
		ArticleRepo:    articleRepo, // Инъекция интерфейса
		ArticleService: articleService,
		ArticleHandler: articleHandler,
	}
}

// Close - закрываем все ресурсы.
func (c *Container) Close() {
	// Close MongoDB client
	mongo.CloseMongoClient()

	// Sync logger before exiting
	log.SyncLogger()

	// Log closing info
	c.Logger.Info("All resources have been closed successfully.")
}
