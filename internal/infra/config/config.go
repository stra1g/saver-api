package config

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"os"
)

type Config struct {
	Server struct {
		Port string `validate:"required,numeric"`
		Host string `validate:"required"`
	}
	Database struct {
		Host     string `validate:"required"`
		Port     string `validate:"required,numeric"`
		User     string `validate:"required"`
		Password string `validate:"required"`
		Name     string `validate:"required"`
		SSLMode  string `validate:"required,oneof=disable require"`
	}
}

func NewConfig() (*Config, error) {
	config := &Config{
		Server: struct {
			Port string `validate:"required,numeric"`
			Host string `validate:"required"`
		}{
			Port: getEnvWithDefault("PORT", "8080"),
			Host: getEnvWithDefault("HOST", "0.0.0.0"),
		},
		Database: struct {
			Host     string `validate:"required"`
			Port     string `validate:"required,numeric"`
			User     string `validate:"required"`
			Password string `validate:"required"`
			Name     string `validate:"required"`
			SSLMode  string `validate:"required,oneof=disable require"`
		}{
			Host:     getEnvWithDefault("DB_HOST", "localhost"),
			Port:     getEnvWithDefault("DB_PORT", "5432"),
			User:     getEnvWithDefault("DB_USER", ""),
			Password: getEnvWithDefault("DB_PASSWORD", ""),
			Name:     getEnvWithDefault("DB_NAME", ""),
			SSLMode:  getEnvWithDefault("DB_SSLMODE", "disable"),
		},
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func validateConfig(config *Config) error {
	validate := validator.New()

	err := validate.Struct(config)
	if err != nil {
		var errMessages []string
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()
			tag := err.Tag()

			switch tag {
			case "required":
				errMessages = append(errMessages, fmt.Sprintf("%s is required", fieldName))
			case "numeric":
				errMessages = append(errMessages, fmt.Sprintf("%s must be a numeric value", fieldName))
			default:
				errMessages = append(errMessages, fmt.Sprintf("%s failed validation: %s", fieldName, tag))
			}
		}

		return errors.New(fmt.Sprintf("Configuration validation failed: %v", errMessages))
	}

	return nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name, c.Database.SSLMode)
}
