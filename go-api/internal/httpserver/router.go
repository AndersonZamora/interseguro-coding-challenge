// Package httpserver ensambla la aplicacion Fiber: middlewares y rutas.
package httpserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/interseguro/coding-challenge/go-api/internal/auth"
	"github.com/interseguro/coding-challenge/go-api/internal/config"
	"github.com/interseguro/coding-challenge/go-api/internal/handler"
)

// New construye la aplicacion Fiber con todos los middlewares y rutas
// registrados: recover -> logging -> CORS -> (JWT solo en /api/v1).
func New(cfg config.Config, statsFetcher handler.StatsFetcher) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "interseguro-go-api",
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.AllowedOrigin,
		AllowHeaders: "Origin, Content-Type, Authorization",
		AllowMethods: "GET, POST, OPTIONS",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	h := handler.NewHandler(statsFetcher, cfg.NodeAPITimeout)

	api := app.Group("/api/v1", auth.Middleware(cfg.JWTSecret))
	api.Post("/qr", h.PostQR)

	return app
}
