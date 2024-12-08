package handlers

import (
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/service"
	"edjr-trk/pkg/http_error"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type userHandler struct {
	service service.UserServiceInterface
	logger  *zap.Logger
}

// UserHandlerInterface - интерфейс для работы с сервисом user.
type UserHandlerInterface interface {
	CreateUser(c *fiber.Ctx) error
	GetAllUsers(c *fiber.Ctx) error
	RemoveUserById(c *fiber.Ctx) error
}

// NewUserHandler creates a new instance of UserHandler.
func NewUserHandler(service service.UserServiceInterface, logger *zap.Logger) UserHandlerInterface {
	return &userHandler{
		service: service,
		logger:  logger,
	}
}

// CreateUser handles creating a new user.
func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	h.logger.Info("Received request to create a new user")

	////Retrieve validated data from context.
	reqInterface := c.Locals("validatedBody")
	body, ok := reqInterface.(dto.CreateUserRequest)

	if !ok {
		h.logger.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	////without validation
	//var body dto.CreateUserRequest
	//if err := c.BodyParser(&body); err != nil {
	//	return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid request body", nil).Send(c)
	//}

	user, err := h.service.CreateNewAdmin(c.Context(), &body)
	if err != nil {
		h.logger.Error("Failed to create article", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create article", nil).Send(c)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetAllUsers handles fetching all users with pagination.
func (h *userHandler) GetAllUsers(c *fiber.Ctx) error {
	h.logger.Info("Received request to fetch paginated users")

	// Retrieve pagination parameters from context.
	pageNumberInterface := c.Locals("pageNumber")
	pageSizeInterface := c.Locals("pageSize")

	pageNumber, ok := pageNumberInterface.(int)
	if !ok {
		h.logger.Error("Failed to retrieve page number from context")
		pageNumber = 1
	}

	pageSize, ok := pageSizeInterface.(int)
	if !ok {
		h.logger.Error("Failed to retrieve page size from context")
		pageSize = 10
	}

	// Fetch articles via the service.
	users, err := h.service.GetAllUsers(c.Context(), pageNumber, pageSize)
	if err != nil {
		h.logger.Error("Failed to fetch paginated users", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to fetch users", nil).Send(c)
	}

	h.logger.Info("Paginated users fetched successfully",
		zap.Int("pageNumber", pageNumber),
		zap.Int("pageSize", pageSize),
		zap.Int("fetchedItems", len(users.Items)),
	)

	return c.Status(fiber.StatusOK).JSON(users)
}

// RemoveUserById handles removing a user by its ID.
func (h *userHandler) RemoveUserById(c *fiber.Ctx) error {
	h.logger.Info("Received request to remove an user")

	// Retrieve article ID from context.
	userIdInterface := c.Locals("userId")
	userId, ok := userIdInterface.(string)
	if !ok || userId == "" {
		h.logger.Error("User ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "User ID is required", nil).Send(c)
	}

	// Remove the article via the service.
	removedUser, err := h.service.RemoveUserById(c.Context(), userId)
	if err != nil {
		h.logger.Error("Failed to remove article", zap.String("userId", userId), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to remove article", nil).Send(c)
	}

	h.logger.Info("User removed successfully", zap.String("userId", userId))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": removedUser})
}
