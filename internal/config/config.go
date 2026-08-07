package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the application configuration parameters.
type Config struct {
	Port        string
	NodeAPIURL  string
	HTTPTimeout time.Duration
	AuthEnabled bool
	JWTSecret   string
	AuthUser    string
	AuthPass    string
}

// LoadConfig loads configuration from environment variables or .env file.
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, using system environment variables")
	}

	port := getEnv("PORT", "3000")
	nodeURL := getEnv("NODE_API_URL", "http://localhost:4000")
	timeoutMsStr := getEnv("HTTP_TIMEOUT_MS", "5000")

	timeoutMs, err := strconv.Atoi(timeoutMsStr)
	if err != nil || timeoutMs <= 0 {
		timeoutMs = 5000
	}

	authEnabledStr := getEnv("AUTH_ENABLED", "false")
	authEnabled, _ := strconv.ParseBool(authEnabledStr)

	jwtSecret := getEnv("JWT_SECRET", "supersecretkey")
	authUser := getEnv("AUTH_USER", "Talsory")
	authPass := getEnv("AUTH_PASS", "Prueba@2026")

	return &Config{
		Port:        port,
		NodeAPIURL:  nodeURL,
		HTTPTimeout: time.Duration(timeoutMs) * time.Millisecond,
		AuthEnabled: authEnabled,
		JWTSecret:   jwtSecret,
		AuthUser:    authUser,
		AuthPass:    authPass,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
