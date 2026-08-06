package handlers

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"matrix-qr-service/internal/client"
	"matrix-qr-service/internal/models"
	"matrix-qr-service/internal/services"
)

// MatrixHandler handles HTTP requests for matrix operations.
type MatrixHandler struct {
	nodeClient client.NodeClient
}

// NewMatrixHandler constructs a MatrixHandler instance.
func NewMatrixHandler(nodeClient client.NodeClient) *MatrixHandler {
	return &MatrixHandler{
		nodeClient: nodeClient,
	}
}

// HandleHealth handles GET /health requests.
func (h *MatrixHandler) HandleHealth(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(models.HealthResponse{
		Status: "ok",
	})
}

// HandleQR handles POST /api/v1/matrix/qr requests.
func (h *MatrixHandler) HandleQR(c *fiber.Ctx) error {
	var req models.MatrixQRRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Cuerpo de la petición JSON no válido o mal formado",
			Details: err.Error(),
		})
	}

	// Validate matrix structure and numeric integrity
	if err := services.ValidateMatrix(req.Matrix); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Matriz de entrada no válida",
			Details: err.Error(),
		})
	}

	// Compute QR decomposition using Modified Gram-Schmidt
	q, r, err := services.ComputeQR(req.Matrix)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Error en la descomposición QR",
			Details: err.Error(),
		})
	}

	// Call downstream Node.js API with context timeout
	stats, err := h.nodeClient.SendStats(c.UserContext(), q, r)
	if err != nil {
		var badGatewayErr *client.ErrBadGateway
		if errors.As(err, &badGatewayErr) {
			return c.Status(http.StatusBadGateway).JSON(models.ErrorResponse{
				Error:   "Error de comunicación con el servicio externo Node.js (502 Bad Gateway)",
				Details: badGatewayErr.Message,
			})
		}

		return c.Status(http.StatusBadGateway).JSON(models.ErrorResponse{
			Error:   "Error llamando a la API de Node.js",
			Details: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(models.MatrixQRResponse{
		Q:     q,
		R:     r,
		Stats: stats,
	})
}

// HandleRotate handles POST /api/v1/matrix/rotate requests.
func (h *MatrixHandler) HandleRotate(c *fiber.Ctx) error {
	var req models.MatrixRotateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Cuerpo de la petición JSON no válido o mal formado",
			Details: err.Error(),
		})
	}

	if err := services.ValidateMatrix(req.Matrix); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Matriz de entrada no válida",
			Details: err.Error(),
		})
	}

	rotated, err := services.RotateMatrix(req.Matrix, req.Direction)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Error rotando la matriz",
			Details: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(models.MatrixRotateResponse{
		Rotated: rotated,
	})
}
