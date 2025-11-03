package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/config"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/routes"
)

func main() {
	// 설정 로드
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// gRPC clients container 초기화
	container, err := client.NewContainer(config.AppConfig.UserServiceAddr, config.AppConfig.TokenServiceAddr)
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	defer container.Close()

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(container.Middleware()) // Context에 클라이언트 주입

	// 모든 경로 앞에 "/api" 설정
	api := router.Group("/api")
	routes.SetupRoutes(api)

	// 설정된 포트로 서버 시작
	router.Run(":" + config.AppConfig.ServerPort)
}
