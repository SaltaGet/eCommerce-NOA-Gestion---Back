package config

import (
	"os"
	"strconv"
	"time"
)

// Config contiene la configuración de la aplicación
type Config struct {
	// Server
	Port int
	Env  string

	// Backend gRPC
	BackendAddr     string
	BackendInsecure bool
	BackendTimeout  time.Duration

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Cache
	CacheTTL time.Duration

	// Logging
	LogLevel string
}

// Load carga la configuración desde variables de entorno
func Load() *Config {
	return &Config{
		Port:            getIntEnv("PORT", 3030),
		Env:             getEnv("APP_ENV", "development"),
		BackendAddr:     getEnv("BACKEND_ADDR", "localhost:3000"),
		BackendInsecure: getBoolEnv("BACKEND_INSECURE", true),
		BackendTimeout:  getDurationEnv("BACKEND_TIMEOUT", 10*time.Second),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6480"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getIntEnv("REDIS_DB", 0),
		CacheTTL:        getDurationEnv("CACHE_TTL", 5*time.Minute),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
	}
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if val := os.Getenv(key); val != "" {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
	}
	return defaultValue
}
