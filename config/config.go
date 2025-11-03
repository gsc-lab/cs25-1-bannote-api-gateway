package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort string

	// Services
	UserServiceAddr  string
	TokenServiceAddr string

	// Google OAuth
	GoogleClientID string

	// Frontend
	FrontendURL string

	// Cookie
	CookieMaxAge int
}

var AppConfig *Config

// LoadConfig 환경 변수에서 값 로드
func LoadConfig() error {
	// .env 파일 로드
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not loaded: %v\n", err)
		fmt.Println("Using environment variables or default values")
	} else {
		fmt.Println("✓ Successfully loaded .env file")
	}

	config := &Config{
		ServerPort:       getEnv("SERVER_PORT", ""),
		UserServiceAddr:  getEnv("USER_SERVICE_ADDR", ""),
		TokenServiceAddr: getEnv("TOKEN_SERVICE_ADDR", ""),
		GoogleClientID:   getEnv("GOOGLE_CLIENT_ID", ""),
		FrontendURL:      getEnv("FRONTEND_URL", ""),
		CookieMaxAge:     getEnvAsInt("COOKIE_MAX_AGE", 900), // 15분 (900초)
	}

	// 필수 환경변수 검증
	if config.GoogleClientID == "" {
		return fmt.Errorf("GOOGLE_CLIENT_ID is required")
	}

	AppConfig = config
	return nil
}

// getEnv env에서 string으로 값 가져오기
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt env에서 int로 값 가져오기
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return defaultValue
	}
	return value
}
