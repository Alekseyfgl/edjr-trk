# edjr-trk

## Project Structure

This project is organized using a modular and scalable structure. Below is an overview of the folder structure and their responsibilities:

```plaintext
edjr-trk/
├── cmd
│   └── main.go              # Main entry point to start the application
├── configs/
│   ├── env/
│   │   └── env.go               # Environment variables loader and helper functions
│   └── mongo/
│       └── mongo.go             # MongoDB connection setup and initialization
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   └── article_handler.go  # Logic for handling "article" API requests
│   │   ├── middlewares/
│   │   │   ├── error_handler.go # Middleware for centralized error handling
│   │   │   └── request_logger.go # Middleware for logging HTTP requests
│   │   └── routes/
│   │       └── article_routes.go # Route definitions for "articles" API
│   ├── ioc/
│   │   └── ioc.go               # Inversion of Control container for dependency injection
│   ├── repository/
│   │   └── article_repo.go      # Data access layer for "articles" in the database
│   └── service/
│       └── article_service.go   # Business logic layer for "articles"
├── pkg/
│   └── log/
│       └── log.go               # Logger setup and utility functions
