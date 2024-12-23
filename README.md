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

```

## Building the Image and container

Run the following command in the root directory of the project (where the Dockerfile is located):

```bash

docker build -t go .
```
### Running the Container with Required Environment Variables
```bash

docker run \
  -p 3000:3000 \
  -e SERV_PORT=3000 \
  -e MONGO_DB_NAME="default_db" \
  -e MONGO_URI="" \
  -e JWT_KEY="vkldfgklfd" \
  -e SUPER_ADMIN_LOGIN="admin" \
  -e SUPER_ADMIN_PASSWORD="admin" \
  --name my-go-app-cnt \
  go
```

*  ```-p``` 3000:3000 maps the container's port 3000 to the host's port 3000.
* The ```-e``` flags specify the environment variables for the container.
* The ```-d``` flag runs the container in detached mode (in the background).
* ```--name``` my-go-app-cnt assigns a custom name to the container for easier management.

### Restarting a Stopped or Crashed Container
```bash

docker start my-go-app-cnt
```