package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   ServerConfig
	DB       DBConfig
	LogLevel string
}

type ServerConfig struct {
	Port             int
	ReadTimeoutSecs  int
	WriteTimeoutSecs int
	IdleTimeoutSecs  int
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	Timeout  time.Duration
}

func Load() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:             getEnvAsInt("SERVER_PORT", 4000),
			WriteTimeoutSecs: getEnvAsInt("SERVER_WRITE_TIMEOUT", 10),
			ReadTimeoutSecs:  getEnvAsInt("SERVER_READ_TIMEOUT", 5),
			IdleTimeoutSecs:  getEnvAsInt("SERVER_IDLE_TIMEOUT", 120),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "Qwerty1,0*"),
			DBName:   getEnv("DB_NAME", "payments_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			MaxConns: getEnvAsInt("DB_MAX_CONNS", 10),
			Timeout: time.Duration(
				getEnvAsInt("DB_TIMEOUT_SECS", 10) * int(time.Second)),
		},
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		value, err := strconv.Atoi(valueStr)
		if err == nil {
			return value
		}
	}
	return defaultValue
}

func (dbc *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbc.Host, dbc.Port, dbc.User, dbc.Password, dbc.DBName, dbc.SSLMode)
}
