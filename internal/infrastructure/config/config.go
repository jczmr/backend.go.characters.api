package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	DBUser      string
	DBPassword  string
	DBName      string
	DBHost      string
	DBPort      string
	LogLevel    string
}

func LoadConfig() (*Config, error) {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}

	cfg := &Config{
		Port:       os.Getenv("PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
	}

	if cfg.Port == "" {
		cfg.Port = "8080" // Default port
	}

	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" || cfg.DBHost == "" || cfg.DBPort == "" {
		return nil, fmt.Errorf("database environment variables (DB_USER, DB_PASSWORD, DB_NAME, DB_HOST, DB_PORT) must be set")
	}

	cfg.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	return cfg, nil
}
