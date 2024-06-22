package wallets

import "github.com/gofiber/fiber/v2"

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the United States of Ganga wallets nigguhhhhh!")
}
