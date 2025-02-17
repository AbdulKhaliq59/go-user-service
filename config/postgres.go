package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}
	port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		return nil, fmt.Errorf("invalid DATABASE_PORT: %v", err)
	}
	return &DatabaseConfig{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     port,
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		DBName:   os.Getenv("DATABASE_DB"),
	}, nil
}

func ConnectDatabase() (*gorm.DB, error) {
	config, err := LoadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.Host, config.User, config.Password, config.DBName, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return db, nil
}
