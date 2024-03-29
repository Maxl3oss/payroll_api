package middleware

import (
	"errors"
	"maxl3oss/pkg/response"
	"maxl3oss/pkg/utils"
	"os"

	"github.com/gofiber/fiber/v2"

	jwtMiddleware "github.com/gofiber/contrib/jwt"
)

// JWTProtected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/contrib/jwt
func JWTProtected() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	// existing middleware config

	config := jwtMiddleware.Config{
		SigningKey:   jwtMiddleware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET_KEY"))},
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

// JWTProtectedAdmin is a middleware to protect routes with JWT authentication for admin roles
func JWTProtectedAdmin() func(*fiber.Ctx) error {
	// Middleware configuration
	config := jwtMiddleware.Config{
		SigningKey:   jwtMiddleware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET_KEY"))},
		ContextKey:   "jwt",
		ErrorHandler: jwtError,
	}

	// Create the JWT middleware instance with the provided configuration
	jwtMiddleware := jwtMiddleware.New(config)

	// Return the middleware function
	return func(c *fiber.Ctx) error {
		// Extract token metadata (You need to implement this function)
		extractToken, err := utils.ExtractTokenMetadata(c)
		if err != nil {
			// Handle error
			return jwtError(c, errors.New("ไม่สามารถเข้าถึงข้อมูล Token ได้"))
		}

		// Log extractToken
		// userId := extractToken.UserID
		isAdmin := extractToken.Credentials["admin"]

		// Check if the user has admin role
		if !isAdmin {
			// Return forbidden error if the user is not an admin
			return jwtError(c, errors.New("สิทธิ์เข้าถึงไม่เพียงพอ"))
		}

		// If the user has admin role, proceed to the next middleware/handler
		return jwtMiddleware(c)
	}
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "401 Unauthorized" {
		return response.Message(c, fiber.StatusBadRequest, false, err.Error())
	}

	// Return status 401 and failed authentication error.
	return response.Message(c, fiber.StatusUnauthorized, false, err.Error())
}
