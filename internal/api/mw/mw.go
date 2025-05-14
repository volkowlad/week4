package mw

import "github.com/gofiber/fiber/v2"

func Authorization(token string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// проверка токена вторизации
		return c.Next()
	}
}
