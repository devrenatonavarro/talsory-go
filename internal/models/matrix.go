package models

// MatrixQRRequest defines the request body for the QR factorization endpoint.
type MatrixQRRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

// MatrixQRResponse defines the response body for the QR factorization endpoint.
type MatrixQRResponse struct {
	Q     [][]float64 `json:"q"`
	R     [][]float64 `json:"r"`
	Stats interface{} `json:"stats"`
}

// NodeStatsPayload defines the body sent to the Node.js API endpoint.
type NodeStatsPayload struct {
	Matrices NodeMatricesPayload `json:"matrices"`
}

// NodeMatricesPayload holds Q and R matrices sent to the Node.js API.
type NodeMatricesPayload struct {
	Q [][]float64 `json:"q"`
	R [][]float64 `json:"r"`
}

// MatrixRotateRequest defines the request body for the matrix rotation endpoint.
type MatrixRotateRequest struct {
	Matrix    [][]float64 `json:"matrix"`
	Direction string      `json:"direction,omitempty"` // "clockwise" or "counterclockwise", defaults to "clockwise"
}

// MatrixRotateResponse defines the response body for matrix rotation.
type MatrixRotateResponse struct {
	Rotated [][]float64 `json:"rotated"`
}

// HealthResponse defines the body returned by the health check endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

// ErrorResponse defines standard structured error response for API clients.
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
