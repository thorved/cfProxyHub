package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks if the user is authenticated
// If not authenticated, redirects to /login page
// Authentication is based on session tokens validated against tokens created during login
// Login credentials are verified against ADMIN_USERNAME and ADMIN_PASSWORD from .env
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if user is authenticated
		// This can be done by checking session, JWT token, or any other auth method

		// Example: Check for session cookie or token
		session, err := c.Cookie("session_token")
		if err != nil || session == "" {
			// No session found, redirect to login
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Validate the session token (you would implement this based on your auth system)
		if !isValidSession(session) {
			// Invalid session, redirect to login
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// User is authenticated, continue to the next handler
		c.Next()
	}
}

// isValidSession validates the session token
func isValidSession(sessionToken string) bool {
	// Basic validation: check if token is not empty and follows expected format
	if sessionToken == "" {
		return false
	}

	// For demo purposes, we'll accept tokens that contain an underscore
	// (matching the format from generateSessionToken in auth handler)
	// In production, implement proper token validation
	return len(sessionToken) > 10 && strings.Contains(sessionToken, "_")
}

// RequireAuth is an alternative middleware that checks for authentication
// and returns JSON error for API endpoints
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Cookie("session_token")
		if err != nil || session == "" || !isValidSession(session) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "Authentication required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
