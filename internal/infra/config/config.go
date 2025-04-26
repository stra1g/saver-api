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
			Port: GetEnvWithDefault("PORT", "8080"),
			Host: GetEnvWithDefault("HOST", "0.0.0.0"),
		},
		Database: struct {
			Host     string `validate:"required"`
			Port     string `validate:"required,numeric"`
			User     string `validate:"required"`
			Password string `validate:"required"`
			Name     string `validate:"required"`
			SSLMode  string `validate:"required,oneof=disable require"`
		}{
			Host:     GetEnvWithDefault("DB_HOST", "localhost"),
			Port:     GetEnvWithDefault("DB_PORT", "5432"),
			User:     GetEnvWithDefault("DB_USER", ""),
			Password: GetEnvWithDefault("DB_PASSWORD", ""),
			Name:     GetEnvWithDefault("DB_NAME", ""),
			SSLMode:  GetEnvWithDefault("DB_SSLMODE", "disable"),
		},
	}

	if err := ValidateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func ValidateConfig(config *Config) error {
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
