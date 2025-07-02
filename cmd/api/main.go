package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"log/slog"

	"backend.go.characters.api/internal/adapters/logger"
	"backend.go.characters.api/internal/adapters/primary/http"
	"backend.go.characters.api/internal/adapters/secondary/db/postgres"
	"backend.go.characters.api/internal/adapters/secondary/dragonballapi"
	"backend.go.characters.api/internal/core/services"
	"backend.go.characters.api/internal/infrastructure/config"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	appLogger := logger.NewSlogLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		appLogger.Error("Failed to load configuration", slog.String("error", err.Error()))
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		appLogger.Error("Failed to connect to database", slog.String("error", err.Error()))
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			appLogger.Error("Failed to close database connection", slog.String("error", err.Error()))
		}
	}()

	// Ping database to ensure connection is established
	for i := 0; i < 5; i++ { // Retry connection
		err = db.Ping()
		if err == nil {
			appLogger.Info("Successfully connected to database!")
			break
		}
		appLogger.Error("Failed to ping database, retrying...", slog.String("error", err.Error()), slog.Int("attempt", i+1))
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		appLogger.Error("Could not connect to database after multiple retries", slog.String("error", err.Error()))
		log.Fatalf("Could not connect to database: %v", err)
	}

	// Apply migrations (simple for demo, consider using a proper migration tool)
	if err := applyMigrations(db, appLogger); err != nil {
		appLogger.Error("Failed to apply database migrations", slog.String("error", err.Error()))
		log.Fatalf("Failed to apply database migrations: %v", err)
	}

	// Initialize adapters
	characterRepository := postgres.NewCharacterRepository(db, appLogger)
	dragonBallAPIClient := dragonballapi.NewDragonBallAPIClient(appLogger)

	// Initialize core service
	characterService := services.NewCharacterService(characterRepository, dragonBallAPIClient, appLogger)

	// Initialize HTTP handler
	characterHandler := http.NewCharacterHandler(characterService, appLogger)

	// Set up Gin router
	router := gin.Default()
	router.POST("/characters", characterHandler.CreateCharacter)

	appLogger.Info(fmt.Sprintf("Starting server on :%s", cfg.Port))
	if err := router.Run(":" + cfg.Port); err != nil {
		appLogger.Error("Failed to start server", slog.String("error", err.Error()))
		log.Fatalf("Failed to run server: %v", err)
	}
}

// applyMigrations is a simple function to apply schema.
// For production, consider using a dedicated migration library like golang-migrate/migrate.
func applyMigrations(db *sql.DB, logger *slog.Logger) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS characters (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		ki VARCHAR(255),
		race VARCHAR(255),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		logger.Error("Error creating characters table", slog.String("error", err.Error()))
		return fmt.Errorf("error creating characters table: %w", err)
	}
	logger.Info("Characters table ensured to exist or created.")
	return nil
}
