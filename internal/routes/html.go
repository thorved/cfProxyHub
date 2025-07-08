package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupHTMLRoutes configures HTML-related routes
func SetupHTMLRoutes(router *gin.Engine) {

	router.GET("/demo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "demo.html", gin.H{})
	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
}
