package config

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	// Server
	ServerPort string

	// Services
	UserServiceAddr      string
	TokenServiceAddr     string
	ScheduleServiceAddr  string
	StudyroomServiceAddr string

	// Google OAuth
	GoogleClientID string

	// Frontend
	FrontendURL string

	// Cookie
	CookieMaxAge int

	// JWT
	PublicKey *rsa.PublicKey
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

	publicKeyBase64 := os.Getenv("PUBLIC_KEY_BASE64")

	publicKey, err := decodePublicKey(publicKeyBase64)

	if err != nil {
		return err
	}

	config := &Config{
		ServerPort:           getEnv("SERVER_PORT", ""),
		UserServiceAddr:      getEnv("USER_SERVICE_ADDR", ""),
		TokenServiceAddr:     getEnv("TOKEN_SERVICE_ADDR", ""),
		ScheduleServiceAddr:  getEnv("SCHEDULE_SERVICE_ADDR", ""),
		StudyroomServiceAddr: getEnv("STUDYROOM_SERVICE_ADDR", ""),
		GoogleClientID:       getEnv("GOOGLE_CLIENT_ID", ""),
		FrontendURL:          getEnv("FRONTEND_URL", ""),
		CookieMaxAge:         getEnvAsInt("COOKIE_MAX_AGE", 900), // 15분 (900초)
		PublicKey:            publicKey,
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

// decodePublicKey JWT 해석에 필요한 퍼블릭 키의 타입 변환
func decodePublicKey(publicKeyBase64 string) (*rsa.PublicKey, error) {

	if publicKeyBase64 == "" {
		return nil, fmt.Errorf("PUBLIC_KEY_BASE64 environment variable is required")
	}

	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)

	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)

	if err != nil {
		return nil, err
	}

	return publicKey, nil

}
