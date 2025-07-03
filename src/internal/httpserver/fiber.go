package httpserver

import (
	"github.com/BernsteinMondy/medods-test-task/src/internal/service"
	"github.com/gofiber/fiber/v2"
)

func NewFiber(service *service.Service) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:               "medods-test-task",
		DisableStartupMessage: true,
		Immutable:             true,
	})
	MapRoutes(app, service)

	return app
}
