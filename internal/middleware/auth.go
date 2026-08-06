package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"matrix-qr-service/internal/models"
)

// JWTMiddleware creates a fiber Handler that checks JWT tokens if auth is enabled.
func JWTMiddleware(authEnabled bool, secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if !authEnabled {
			return c.Next()
		}

		authHeader := c.Get(fiber.HeaderAuthorization)
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "Se requiere token de autenticación (Header Authorization missing)",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error: "Formato de token no válido. Debe ser 'Bearer <token>'",
			})
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma no esperado: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
				Error:   "Token JWT inválido o expirado",
				Details: err.Error(),
			})
		}

		c.Locals("user", token)
		return c.Next()
	}
}
