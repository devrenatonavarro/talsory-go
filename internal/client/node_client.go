package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"matrix-qr-service/internal/models"
)

// NodeClient interface defines the contract for sending matrix statistics to the Node.js API.
type NodeClient interface {
	SendStats(ctx context.Context, q, r [][]float64) (interface{}, error)
}

type nodeClient struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

// NewNodeClient creates a new HTTP client for the Node.js API.
func NewNodeClient(baseURL string, timeout time.Duration) NodeClient {
	// Clean baseURL trailing slash if present
	cleanURL := strings.TrimRight(baseURL, "/")
	return &nodeClient{
		baseURL: cleanURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

// ErrBadGateway represents a failure communicating with the downstream Node.js API.
type ErrBadGateway struct {
	Message string
}

func (e *ErrBadGateway) Error() string {
	return e.Message
}

// SendStats sends Q and R matrices to POST {NODE_API_URL}/api/v1/stats and returns the parsed stats.
func (c *nodeClient) SendStats(parentCtx context.Context, q, r [][]float64) (interface{}, error) {
	endpoint := fmt.Sprintf("%s/api/v1/stats", c.baseURL)

	payload := models.NodeStatsPayload{
		Matrices: models.NodeMatricesPayload{
			Q: q,
			R: r,
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error serializando payload para Node API: %w", err)
	}

	ctx, cancel := context.WithTimeout(parentCtx, c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, &ErrBadGateway{Message: fmt.Sprintf("error creando solicitud hacia Node API: %v", err)}
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, &ErrBadGateway{Message: fmt.Sprintf("timeout de conexión con Node API tras %v", c.timeout)}
		}
		return nil, &ErrBadGateway{Message: fmt.Sprintf("no se pudo conectar con Node API (%s): %v", c.baseURL, err)}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ErrBadGateway{Message: fmt.Sprintf("error leyendo respuesta de Node API: %v", err)}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &ErrBadGateway{
			Message: fmt.Sprintf("Node API devolvió un código de estado no exitoso: %d (%s)", resp.StatusCode, string(respBody)),
		}
	}

	var parsedStats interface{}
	if len(respBody) > 0 {
		if err := json.Unmarshal(respBody, &parsedStats); err != nil {
			// If not valid JSON, return raw string output inside interface
			parsedStats = string(respBody)
		}
	} else {
		parsedStats = map[string]interface{}{}
	}

	return parsedStats, nil
}
