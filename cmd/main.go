package main

import (
	"context"
	"edjr-trk/configs/env"
	"edjr-trk/internal/api/middlewares"
	"edjr-trk/internal/api/routes"
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load environment variables
	env.LoadEnv()

	// Initialize dependencies via IoC container
	container := ioc.NewContainer()

	// Setup Fiber application
	app := fiber.New()

	// Middleware: CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins, consider limiting this for production
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
	}))

	// Middleware: Global error handling
	app.Use(middlewares.ErrorHandlerMiddleware(container.Logger))
	// Middleware: Request logging
	app.Use(middlewares.RequestLoggerMiddleware(container.Logger))

	// Create a group for API routes with prefix `/api`
	api := app.Group("/api")

	// Register routes within the `/api` group
	routes.RegisterArticleRoutes(api, container)
	routes.RegisterUserRoutes(api, container)
	routes.RegisterAuthRoutes(api, container)
	routes.RegisterEmailRoutes(api, container)

	// Start the server
	port := env.GetEnv("SERV_PORT", "3000")
	container.Logger.Info("Starting server", zap.String("port", port))

	// Graceful shutdown handling in a goroutine
	go func() {
		if err := app.Listen(":" + port); err != nil {
			container.Logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Call graceful shutdown handler
	handleGracefulShutdown(app, container.Logger)
}

// handleGracefulShutdown handles signal-based graceful shutdown
func handleGracefulShutdown(app *fiber.App, logger *zap.Logger) {
	// Create a channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	<-signalChan

	// Gracefully shutdown the server
	shutdownTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	logger.Info("Shutting down server...")
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Failed to gracefully shutdown server", zap.Error(err))
	} else {
		logger.Info("Server shut down gracefully")
	}
}
