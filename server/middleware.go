package server

import (
	"base/utils"
	"fmt"
	"strings"

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

func jwtAuth(c *fiber.Ctx) error {
	noAuth := []string{"/auth/sign-in", "/auth/sign-up"}
	currentPath := c.Path()

	for _, v := range noAuth {
		v = "/api/v1" + v
		if v == currentPath {
			return c.Next()
		}
	}

	bearToken := c.Get("Authorization", "")
	if bearToken == "" {
		return c.Status(401).JSON(fiber.Map{
			"isSuccess": false,
			"message":   "Not Authorized",
		})
	}

	splitted := strings.Split(bearToken, " ")
	if len(splitted) != 2 {
		return c.Status(401).JSON(fiber.Map{
			"isSuccess": false,
			"message":   "Invalid token",
		})
	}

	tokenPart := splitted[1]
	claims, err := utils.VerifyJwt(tokenPart)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"isSuccess": false,
			"message":   err,
		})
	}

	c.Locals("user", map[string]interface{}(claims))
	return c.Next()

}
