package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/algoholic/services"
)

// AuthMiddleware is a middleware for JWT authentication
func AuthMiddleware(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		tokenString := parts[1]

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Store user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("email", claims.Email)

		return c.Next()
	}
}

// OptionalAuthMiddleware is a middleware that extracts auth info if present but doesn't require it
func OptionalAuthMiddleware(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				claims, err := authService.ValidateToken(parts[1])
				if err == nil {
					c.Locals("user_id", claims.UserID)
					c.Locals("username", claims.Username)
					c.Locals("email", claims.Email)
				}
			}
		}
		return c.Next()
	}
}

// GetUserID extracts user ID from context
func GetUserID(c *fiber.Ctx) (int, bool) {
	userID, ok := c.Locals("user_id").(int)
	return userID, ok
}

// GetUsername extracts username from context
func GetUsername(c *fiber.Ctx) (string, bool) {
	username, ok := c.Locals("username").(string)
	return username, ok
}
