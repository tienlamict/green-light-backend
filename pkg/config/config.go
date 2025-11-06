package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
	Upload   UploadConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port string
	ENV  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	DSN      string
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type CORSConfig struct {
	Origins []string
}

type UploadConfig struct {
	Dir         string
	MaxFileSize int64
}

type LogConfig struct {
	Level string
}

func Load() (*Config, error) {
	// Load .env file if exists (ignore error in production)
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			ENV:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", "root"),
			DBName:   getEnv("DB_NAME", "greenlight_db"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "change-me-in-production"),
			ExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),
		},
		CORS: CORSConfig{
			Origins: strings.Split(getEnv("CORS_ORIGINS", "*"), ","),
		},
		Upload: UploadConfig{
			Dir:         getEnv("UPLOAD_DIR", "./uploads"),
			MaxFileSize: getEnvAsInt64("MAX_UPLOAD_SIZE", 10485760), // 10MB
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	// Build DSN
	cfg.Database.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}

func getEnvAsInt64(key string, fallback int64) int64 {
	strValue := getEnv(key, "")
	if value, err := strconv.ParseInt(strValue, 10, 64); err == nil {
		return value
	}
	return fallback
}
