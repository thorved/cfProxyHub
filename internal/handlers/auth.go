package handlers

import (
	"net/http"
	"time"

	"cfPorxyHub/internal/config"

	"github.com/gin-gonic/gin"
)

// LoginHandler handles user login
type LoginHandler struct {
	config *config.Config
}

// NewLoginHandler creates a new login handler
func NewLoginHandler(cfg *config.Config) *LoginHandler {
	return &LoginHandler{
		config: cfg,
	}
}

// LoginForm displays the login form
func (h *LoginHandler) LoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

// Login processes the login form submission
func (h *LoginHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	// TODO: Implement actual authentication logic here
	// This is a placeholder - replace with your actual authentication
	if h.authenticateUser(username, password) {
		// Create session token (you should implement proper session management)
		sessionToken := h.generateSessionToken(username)

		// Set session cookie
		c.SetCookie(
			"session_token", // name
			sessionToken,    // value
			3600*24*7,       // maxAge (7 days)
			"/",             // path
			"",              // domain
			false,           // secure (set to true in production with HTTPS)
			true,            // httpOnly
		)

		// Redirect to dashboard/home page
		c.Redirect(http.StatusFound, "/")
		return
	}

	// Authentication failed
	c.HTML(http.StatusUnauthorized, "login.html", gin.H{
		"title": "Login",
		"error": "Invalid username or password",
	})
}

// Logout handles user logout
func (h *LoginHandler) Logout(c *gin.Context) {
	// Clear the session cookie
	c.SetCookie(
		"session_token",
		"",
		-1, // maxAge negative to delete
		"/",
		"",
		false,
		true,
	)

	// Redirect to login page
	c.Redirect(http.StatusFound, "/login")
}

// authenticateUser validates user credentials
func (h *LoginHandler) authenticateUser(username, password string) bool {
	// Check against credentials from .env file
	return username == h.config.AdminUsername && password == h.config.AdminPassword
}

// generateSessionToken creates a session token for the user
func (h *LoginHandler) generateSessionToken(username string) string {
	// Verify the username matches admin username from config
	if username != h.config.AdminUsername {
		// This should never happen as we validate in authenticateUser
		// but adding as an extra security measure
		return ""
	}

	// In a real application, you would:
	// 1. Generate a cryptographically secure random token
	// 2. Store it in database/redis with expiration
	// 3. Associate it with the user

	// For demo purposes, create a simple token with timestamp
	return username + "_" + time.Now().Format("20060102150405")
}
