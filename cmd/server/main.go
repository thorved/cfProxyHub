package main

import (
	"cfPorxyHub/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	routes.SetupRoutes(router)
	// Start the server on port 8080
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
