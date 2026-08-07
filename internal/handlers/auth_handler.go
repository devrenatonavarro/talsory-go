package handlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"matrix-qr-service/internal/config"
	"matrix-qr-service/internal/models"
)

// AuthHandler handles HTTP requests for user authentication.
type AuthHandler struct {
	cfg *config.Config
}

// NewAuthHandler constructs a new AuthHandler instance.
func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		cfg: cfg,
	}
}

// HandleLogin handles POST /api/v1/auth/login requests.
func (h *AuthHandler) HandleLogin(c *fiber.Ctx) error {
	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "Cuerpo de la petición JSON no válido o mal formado",
			Details: err.Error(),
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(models.ErrorResponse{
			Error: "Se requieren los campos 'username' y 'password'",
		})
	}

	if req.Username != h.cfg.AuthUser || req.Password != h.cfg.AuthPass {
		return c.Status(http.StatusUnauthorized).JSON(models.ErrorResponse{
			Error: "Credenciales de usuario o contraseña incorrectas",
		})
	}

	// Create JWT token with 24h expiration
	claims := jwt.MapClaims{
		"sub": req.Username,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "Error al generar el token JWT",
			Details: err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(models.LoginResponse{
		Token: tokenString,
		User: models.UserInfo{
			Username: req.Username,
		},
	})
}
