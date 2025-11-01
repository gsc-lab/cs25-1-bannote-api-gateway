package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/grpc/client"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/routes"
)

func main() {
	// gRPC clients container 초기화
	//container, err := client.NewContainer("localhost:9090", "localhost:9091")
	container, err := client.NewContainer("host.docker.internal:9090", "host.docker.internal:9091")
	if err != nil {
		log.Fatalf("Failed to initialize gRPC clients: %v", err)
	}
	defer container.Close()

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(container.Middleware()) // Context에 클라이언트 주입

	routes.SetupRoutes(router)

	router.Run()
}
