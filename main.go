package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gsc-lab/cs25-1-bannote-api-gateway/routes"
)

func main() {
	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run()
}
