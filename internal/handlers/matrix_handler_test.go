package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"matrix-qr-service/internal/client"
	"matrix-qr-service/internal/models"
)

func setupTestApp(nodeURL string) *fiber.App {
	app := fiber.New()
	nodeClient := client.NewNodeClient(nodeURL, 200*time.Millisecond)
	matrixHandler := NewMatrixHandler(nodeClient)

	app.Get("/health", matrixHandler.HandleHealth)
	api := app.Group("/api/v1")
	api.Post("/matrix/qr", matrixHandler.HandleQR)
	api.Post("/matrix/rotate", matrixHandler.HandleRotate)

	return app
}

func TestHandleHealth(t *testing.T) {
	app := setupTestApp("http://localhost:9999")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error ejecutando solicitud a /health: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Código de estado esperado 200, obtenido %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var res models.HealthResponse
	if err := json.Unmarshal(body, &res); err != nil {
		t.Fatalf("Error deserializando respuesta: %v", err)
	}

	if res.Status != "ok" {
		t.Errorf("Estado esperado 'ok', obtenido '%s'", res.Status)
	}
}

func TestHandleRotate(t *testing.T) {
	app := setupTestApp("http://localhost:9999")

	t.Run("Rotación exitosa clockwise", func(t *testing.T) {
		reqPayload := models.MatrixRotateRequest{
			Matrix: [][]float64{
				{1, 2},
				{3, 4},
			},
			Direction: "clockwise",
		}
		bodyBytes, _ := json.Marshal(reqPayload)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/matrix/rotate", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Error al probar endpoint rotate: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Código de estado esperado 200, obtenido %d", resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		var res models.MatrixRotateResponse
		_ = json.Unmarshal(body, &res)

		expected := [][]float64{
			{3, 1},
			{4, 2},
		}
		if len(res.Rotated) != 2 || res.Rotated[0][0] != expected[0][0] || res.Rotated[1][1] != expected[1][1] {
			t.Errorf("Resultado de rotación inesperado: %v", res.Rotated)
		}
	})

	t.Run("Error 400 por matriz no rectangular", func(t *testing.T) {
		invalidBody := `{"matrix": [[1, 2], [3]]}`
		req := httptest.NewRequest(http.MethodPost, "/api/v1/matrix/rotate", bytes.NewReader([]byte(invalidBody)))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Error al probar endpoint rotate: %v", err)
		}

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Código de estado esperado 400 Bad Request, obtenido %d", resp.StatusCode)
		}
	})
}

func TestHandleQR_Success(t *testing.T) {
	// Mock Node.js API server
	mockNodeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/stats" && r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"norm_q": 1.0, "mean_r": 3.5, "status": "processed"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockNodeServer.Close()

	app := setupTestApp(mockNodeServer.URL)

	reqPayload := models.MatrixQRRequest{
		Matrix: [][]float64{
			{1, 2},
			{3, 4},
			{5, 6},
		},
	}
	bodyBytes, _ := json.Marshal(reqPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/matrix/qr", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error en petición /api/v1/matrix/qr: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Código de estado esperado 200, obtenido %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var qrResp models.MatrixQRResponse
	if err := json.Unmarshal(body, &qrResp); err != nil {
		t.Fatalf("Error al deserializar respuesta QR: %v", err)
	}

	if len(qrResp.Q) != 3 || len(qrResp.R) != 2 {
		t.Errorf("Dimensiones de Q o R incorrectas en la respuesta HTTP: Q len=%d, R len=%d", len(qrResp.Q), len(qrResp.R))
	}
	if qrResp.Stats == nil {
		t.Error("Se esperaban las estadísticas devueltas por Node.js API, pero se obtuvo nil")
	}
}

func TestHandleQR_NodeAPIDown_Returns502(t *testing.T) {
	// Point node client to a closed/unreachable server address
	app := setupTestApp("http://127.0.0.1:59999")

	reqPayload := models.MatrixQRRequest{
		Matrix: [][]float64{
			{1, 2},
			{3, 4},
		},
	}
	bodyBytes, _ := json.Marshal(reqPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/matrix/qr", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Error en petición /api/v1/matrix/qr: %v", err)
	}

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Código de estado esperado 503 Service Unavailable cuando Node API está inalcanzable, obtenido %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var errResp models.ErrorResponse
	_ = json.Unmarshal(body, &errResp)

	if errResp.Error == "" {
		t.Error("Se esperaba mensaje de error descriptivo en la respuesta 502")
	}
}
