package ioc

import (
	"edjr-trk/configs/env"
	"edjr-trk/configs/mongo"
	"edjr-trk/internal/api/handlers"
	"edjr-trk/internal/repository"
	"edjr-trk/internal/service"
	"edjr-trk/pkg/log"
	mongodb "go.mongodb.org/mongo-driver/mongo" // Импортируем MongoDB драйвер правильн
	"go.uber.org/zap"
)

// Container - структура для хранения зависимостей.
type Container struct {
	Logger         *zap.Logger
	MongoClient    *mongodb.Client
	ArticleRepo    repository.ArticleRepositoryInterface
	UserRepo       repository.UserRepositoryInterface
	ArticleService service.ArticleServiceInterface
	UserService    service.UserServiceInterface
	JwtService     service.JWTServiceInterface
	AuthService    service.AuthServiceInterface
	ArticleHandler *handlers.ArticleHandler
	UserHandler    handlers.UserHandlerInterface
	AuthHandler    handlers.AuthHandlerInterface
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

	// Create repositories
	articleRepo := repository.NewArticleRepository(clientDB, logger)
	userRepo := repository.NewUserRepository(clientDB, logger)
	// Create services
	articleService := service.NewArticleService(articleRepo, logger)
	userService := service.NewUserService(userRepo, logger)
	jwtService := service.NewJWTService(env.GetEnv("JWT_TOKEN", ""), logger)
	authService := service.NewAuthService(userRepo, jwtService, logger)
	// Create handlers
	articleHandler := handlers.NewArticleHandler(articleService, logger)
	userHandler := handlers.NewUserHandler(userService, logger)
	authHandler := handlers.NewAuthHandler(authService, logger)

	// Return the container with all dependencies
	return &Container{
		Logger:         logger,
		MongoClient:    clientDB,
		ArticleRepo:    articleRepo,
		UserRepo:       userRepo,
		ArticleService: articleService,
		UserService:    userService,
		JwtService:     jwtService,
		ArticleHandler: articleHandler,
		UserHandler:    userHandler,
		AuthHandler:    authHandler,
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
