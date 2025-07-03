package httpserver

import (
	"errors"
	"github.com/BernsteinMondy/medods-test-task/src/internal/authorization"
	"github.com/BernsteinMondy/medods-test-task/src/internal/service"
	"github.com/gofiber/fiber/v2"
)

func authenticated(tokenSvc service.TokenService) fiber.Handler {
	return func(f *fiber.Ctx) error {
		authHeader := f.Get("Authorization")
		if authHeader == "" {
			return f.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		_, err := tokenSvc.ParseToken(authHeader)
		if err != nil {
			if errors.Is(err, authorization.ErrTokenExpired) {
				return f.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Token expired",
				})
			}
			return f.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		return f.Next()
	}
}
