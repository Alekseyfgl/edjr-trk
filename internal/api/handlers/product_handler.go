package handlers

import (
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/service"
	"edjr-trk/pkg/http_error"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ProductHandler struct {
	service service.ProductServiceInterface
	logger  *zap.Logger
}

func NewProductHandler(service service.ProductServiceInterface, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	h.logger.Info("Received request to create a new product")

	reqInterface := c.Locals("validatedBody")
	req, ok := reqInterface.(dto.CreateProductRequest)
	if !ok {
		h.logger.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	product, err := h.service.CreateProduct(c.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create product", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create product", nil).Send(c)
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
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

	products, err := h.service.GetAllProducts(c.Context(), pageNumber, pageSize)
	if err != nil {
		h.logger.Error("Failed to fetch paginated products", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to fetch products", nil).Send(c)
	}

	h.logger.Info("Paginated articles fetched successfully",
		zap.Int("pageNumber", pageNumber),
		zap.Int("pageSize", pageSize),
		zap.Int("fetchedItems", len(products.Items)),
	)

	return c.Status(fiber.StatusOK).JSON(products)
}

func (h *ProductHandler) GetProductById(c *fiber.Ctx) error {
	h.logger.Info("Received request to fetch an product by ID")

	productIDInterface := c.Locals("productID")
	productID, ok := productIDInterface.(string)
	if !ok || productID == "" {
		h.logger.Error("Product ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Product ID is required", nil).Send(c)
	}

	article, err := h.service.GetProductById(c.Context(), productID)
	if err != nil {
		h.logger.Error("Failed to fetch product by ID", zap.String("productID", productID), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to fetch product", nil).Send(c)
	}

	if article == nil {
		return http_error.NewHTTPError(fiber.StatusNotFound, "Product not found", nil).Send(c)
	}

	h.logger.Info("Product fetched successfully", zap.String("productID", productID))
	return c.Status(fiber.StatusOK).JSON(article)
}

func (h *ProductHandler) RemoveProductById(c *fiber.Ctx) error {
	h.logger.Info("Received request to remove an product")

	productIDInterface := c.Locals("productID")
	productID, ok := productIDInterface.(string)
	if !ok || productID == "" {
		h.logger.Error("Product ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Product ID is required", nil).Send(c)
	}

	removedProduct, err := h.service.RemoveProductById(c.Context(), productID)
	if err != nil {
		h.logger.Error("Failed to remove article", zap.String("productID", productID), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to remove product", nil).Send(c)
	}

	h.logger.Info("Product removed successfully", zap.String("productID", productID))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": removedProduct})
}

func (h *ProductHandler) PatchProductById(c *fiber.Ctx) error {
	h.logger.Info("Received request to patch an product")

	productIDInterface := c.Locals("productID")
	productID, ok := productIDInterface.(string)
	if !ok || productID == "" {
		h.logger.Error("Product ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Product ID is required", nil).Send(c)
	}

	// Retrieve validated data from context.
	reqInterface := c.Locals("validatedBody")
	req, ok := reqInterface.(dto.PatchProductRequest)
	if !ok {
		h.logger.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	// Update the article via the service.
	updatedArticle, err := h.service.PatchProductById(c.Context(), req, productID)
	if err != nil {
		h.logger.Error("Failed to update article", zap.String("productID", productID), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to update product", nil).Send(c)
	}

	h.logger.Info("Product updated successfully", zap.String("productID", productID))
	return c.Status(fiber.StatusOK).JSON(updatedArticle)
}
