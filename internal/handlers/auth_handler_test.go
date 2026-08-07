package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"matrix-qr-service/internal/config"
	"matrix-qr-service/internal/models"
)

func setupAuthTestApp() (*fiber.App, *config.Config) {
	cfg := &config.Config{
		Port:        "3000",
		NodeAPIURL:  "http://localhost:4000",
		HTTPTimeout: 5 * time.Second,
		AuthEnabled: true,
		JWTSecret:   "test_secret_key",
		AuthUser:    "Talsory",
		AuthPass:    "Prueba@2026",
	}

	app := fiber.New()
	authHandler := NewAuthHandler(cfg)

	app.Post("/api/v1/auth/login", authHandler.HandleLogin)

	return app, cfg
}

func TestHandleLogin_Success(t *testing.T) {
	app, cfg := setupAuthTestApp()

	payload := models.LoginRequest{
		Username: "Talsory",
		Password: "Prueba@2026",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error realizando petición: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Código de estado esperado 200, obtenido %d", resp.StatusCode)
	}

	var loginResp models.LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("Error al deserializar respuesta de login: %v", err)
	}

	if loginResp.Token == "" {
		t.Error("Se esperaba un token JWT no vacío")
	}

	if loginResp.User.Username != "Talsory" {
		t.Errorf("Usuario esperado 'Talsory', obtenido '%s'", loginResp.User.Username)
	}

	// Verify token validity with test_secret_key
	token, err := jwt.Parse(loginResp.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		t.Errorf("Error o token JWT inválido: %v", err)
	}
}

func TestHandleLogin_InvalidCredentials(t *testing.T) {
	app, _ := setupAuthTestApp()

	payload := models.LoginRequest{
		Username: "Talsory",
		Password: "wrongpassword",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error realizando petición: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Código de estado esperado 401, obtenido %d", resp.StatusCode)
	}
}

func TestHandleLogin_MissingFields(t *testing.T) {
	app, _ := setupAuthTestApp()

	payload := models.LoginRequest{
		Username: "",
		Password: "Prueba@2026",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error realizando petición: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Código de estado esperado 400, obtenido %d", resp.StatusCode)
	}
}
