package httpserver

import (
	"fmt"
	"github.com/BernsteinMondy/medods-test-task/src/internal/service"
	"github.com/gofiber/fiber/v2"
)

func MapRoutes(f *fiber.App, service *service.Service) {
	registerGroup := f.Group("/register")
	registerGroup.Post("/", registerHandler(service))

	loginGroup := f.Group("/login")
	loginGroup.Post("/login", loginHandler())

	authenticate := authenticated(service.TokenService)
	authenticatedRoutes := f.Group("/", authenticate)

	userGroup := authenticatedRoutes.Group("/users")
	userGroup.Get("/", getUserIDHandler(service))

	tokensGroup := userGroup.Group("/tokens")
	tokensGroup.Get("/:userID")
	tokensGroup.Put("/refresh", refreshTokensHandler(service))
}

func registerHandler(srvc *service.Service) fiber.Handler {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}
	return func(f *fiber.Ctx) error {
		ctx := f.UserContext()

		var req request
		err := f.BodyParser(&req)
		if err != nil {
			f.Status(fiber.StatusBadRequest)
			return nil
		}

		if len(req.Username) == 0 || len(req.Username) > 32 {
			f.Status(fiber.StatusBadRequest)
			return nil
		}

		if len(req.Password) == 0 || len(req.Password) > 32 {
			f.Status(fiber.StatusBadRequest)
			return nil
		}

		ip := f.IP()
		userAgent := f.Get("User-Agent")

		if len(userAgent) == 0 || len(ip) == 0 {
			f.Status(fiber.StatusBadRequest)
			return nil
		}

		data := &service.RegisterUserData{
			Username:  req.Username,
			Password:  req.Password,
			IP:        ip,
			UserAgent: userAgent,
		}

		refreshToken, accessToken, err := srvc.RegisterUser(ctx, data)
		if err != nil {
			return fmt.Errorf("service: register user: %w", err)
		}

		var resp = response{
			RefreshToken: refreshToken,
			AccessToken:  accessToken,
		}

		err = f.JSON(resp)
		if err != nil {
			return fmt.Errorf("write json body: %w", err)
		}

		f.Status(fiber.StatusOK)
		return nil
	}
}

func loginHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func getUserIDHandler(service *service.Service) fiber.Handler {
	type response struct {
		ID string `json:"id"`
	}
	return func(f *fiber.Ctx) error {
		ctx := f.UserContext()

		token := f.Get("Authorization")

		id, err := service.GetUserID(ctx, token)
		if err != nil {
			return fmt.Errorf("serivce: get user id: %w", err)
		}

		var resp = response{
			ID: id.String(),
		}

		err = f.JSON(resp)
		if err != nil {
			return fmt.Errorf("write json body: %w", err)
		}

		f.Status(fiber.StatusOK)
		return nil
	}
}

func refreshTokensHandler(service *service.Service) fiber.Handler {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}
	type response struct {
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}
	return func(f *fiber.Ctx) error {
		ctx := f.UserContext()

		var req request
		err := f.BodyParser(&req)
		if err != nil {
			f.Status(fiber.StatusBadRequest)
			return nil
		}

		if len(req.RefreshToken) == 0 {
			f.Status(fiber.StatusBadRequest)
			return nil
		}

		ip := f.IP()
		userAgent := f.Get("User-Agent")

		refresh, access, err := service.RefreshTokenPair(ctx, req.RefreshToken, userAgent, ip)
		if err != nil {
			return fmt.Errorf("service: refresh token pair: %w", err)
		}

		var resp = response{
			RefreshToken: refresh,
			AccessToken:  access,
		}

		err = f.JSON(resp)
		if err != nil {
			return fmt.Errorf("write json body: %w", err)
		}

		f.Status(fiber.StatusOK)
		return nil
	}
}
