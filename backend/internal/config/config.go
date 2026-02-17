package config

import (
	"os"
	"strings"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	UploadDir   string
	CORSOrigins []string
	DataDir     string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "./trackside.db"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
		CORSOrigins: strings.Split(getEnv("CORS_ORIGINS", "http://localhost:3000"), ","),
		DataDir:     getEnv("DATA_DIR", "../trackside/data"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
