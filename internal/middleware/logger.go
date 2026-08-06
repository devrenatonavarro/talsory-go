package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

// RequestIDMiddleware ensures every request has a unique Request ID.
func RequestIDMiddleware() fiber.Handler {
	return requestid.New()
}

// LoggerMiddleware logs Request ID, HTTP Method, Path, Status Code, and Duration.
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)
		reqID := c.GetRespHeader(fiber.HeaderXRequestID)
		if reqID == "" {
			reqID = c.IP()
		}

		status := c.Response().StatusCode()
		log.Printf("[HTTP] RequestID: %s | Method: %s | Path: %s | Status: %d | Latency: %v",
			reqID, c.Method(), c.Path(), status, duration)

		return err
	}
}
