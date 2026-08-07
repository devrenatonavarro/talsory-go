package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"matrix-qr-service/internal/client"
	"matrix-qr-service/internal/config"
	"matrix-qr-service/internal/handlers"
	"matrix-qr-service/internal/middleware"
)

func main() {
	cfg := config.LoadConfig()

	app := fiber.New(fiber.Config{
		AppName: "Matrix QR Factorization & Rotation API v1.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Global Middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, OPTIONS",
	}))
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.LoggerMiddleware())

	// Initialize Clients & Handlers
	nodeClient := client.NewNodeClient(cfg.NodeAPIURL, cfg.HTTPTimeout)
	matrixHandler := handlers.NewMatrixHandler(nodeClient)
	authHandler := handlers.NewAuthHandler(cfg)

	// Register Public Routes
	app.Get("/health", matrixHandler.HandleHealth)
	app.Post("/api/v1/auth/login", authHandler.HandleLogin)

	// API Group /api/v1
	api := app.Group("/api/v1")

	// Apply JWT Authentication Middleware if enabled
	api.Use(middleware.JWTMiddleware(cfg.AuthEnabled, cfg.JWTSecret))

	// Register Matrix Routes
	api.Post("/matrix/qr", matrixHandler.HandleQR)
	api.Post("/matrix/rotate", matrixHandler.HandleRotate)

	// Graceful Shutdown Channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Servidor escuchando en el puerto %s (AUTH_ENABLED=%v)...", cfg.Port, cfg.AuthEnabled)
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	<-stop
	log.Println("Apagando el servidor elegantemente (Graceful Shutdown)...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Error durante la detención del servidor: %v", err)
	}
	log.Println("Servidor detenido correctamente.")
}
