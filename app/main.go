package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/lelouchhh/friendly-basketball-media/internal/delivery/http"
	"github.com/lelouchhh/friendly-basketball-media/internal/delivery/http/middleware"
	logger2 "github.com/lelouchhh/friendly-basketball-media/internal/logger"
	"github.com/lelouchhh/friendly-basketball-media/internal/repostitory/postgres"
	"github.com/lelouchhh/friendly-basketball-media/video"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	defaultTimeout = 30
	defaultAddress = ":9090"
)

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables.")
	}
}

func getEnvOrDefault(envKey, defaultValue string) string {
	value := os.Getenv(envKey)
	if value == "" {
		return defaultValue
	}
	return value
}

func getTimeout() time.Duration {
	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Printf("Invalid CONTEXT_TIMEOUT value: %s. Using default: %d seconds\n", timeoutStr, defaultTimeout)
		timeout = defaultTimeout
	}
	return time.Duration(timeout) * time.Second
}

func connectToDatabase() *sqlx.DB {
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := os.Getenv("DB_PASSWORD") // Требуется установить в окружении
	dbName := getEnvOrDefault("DB_NAME", "friendly_basketball")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Connected to the database successfully.")
	return db
}

func startServer(e *echo.Echo, address string) {
	certFile := "/etc/letsencrypt/live/www.969975-cv27771.tmweb.ru/fullchain.pem"
	keyFile := "/etc/letsencrypt/live/www.969975-cv27771.tmweb.ru/privkey.pem"

	// Attempt to start HTTPS server first
	if _, certErr := os.Stat(certFile); os.IsNotExist(certErr) {
		// If certificates are not found, log it and fallback to HTTP
		log.Println("SSL certificates not found, starting HTTP server.")
		log.Fatal(e.Start(address)) // Start HTTP server
	} else {
		// Try to start HTTPS server
		log.Printf("SSL certificates found, starting HTTPS server on %s...\n", address)
		log.Fatal(e.StartTLS(address, certFile, keyFile)) // Start HTTPS server
	}
}

func main() {
	logger, err := logger2.NewZapLogger("info")

	loadEnv()
	logger.Info("Starting server...")

	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.CORS)

	timeoutContext := getTimeout()
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			logger.Info("Incoming HTTP request",
				zap.String("path", c.Request().URL.Path),
				zap.String("method", c.Request().Method),
				zap.Time("time", time.Now()),
			)
			return next(c)
		}
	})
	db := connectToDatabase()
	err = db.Ping()
	e.Static("/videos", "./uploads")
	logger.Info("address", zap.String("address", e.Server.Addr))
	videoRepo := postgres.NewVideo(db, logger)
	videoService := video.NewService(videoRepo, logger)
	http.NewVideoHandler(e, os.Getenv("JWT_SECRET"), videoService, logger)

	if err != nil {
		logger.Error("Can't connect to database", zap.Error(err))
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing the database connection: %v", err)
		}
	}()

	address := getEnvOrDefault("SERVER_ADDRESS", defaultAddress)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		logger.Info("Starting server...")
		startServer(e, address)
	}()
	<-stop
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
