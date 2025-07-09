package main

import (
	"cfProxyHub/internal/config"
	"cfProxyHub/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	router := gin.Default()

	routes.SetupRoutes(router)
	// Start the server on configured port
	if err := router.Run(":" + cfg.Port); err != nil {
		panic(err)
	}
}
