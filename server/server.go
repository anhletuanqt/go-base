package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Setup() *fiber.App {
	app := new()
	initMiddleware(app)

	return app
}

func new() *fiber.App {
	return fiber.New(fiber.Config{
		ErrorHandler: errorHandle,
	})
}

func initMiddleware(app *fiber.App) {
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(jwtAuth)
}
