package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func errorHandle(c *fiber.Ctx, err error) error {
	code := 400
	message := err.Error()
	fmt.Printf("%+v\n", err)
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"isSuccess": false,
		"message":   message,
	})
}
