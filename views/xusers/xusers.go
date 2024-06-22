package xusers

import (
	"txrnxp/utils/xusers"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the United States of Ganga users biaaattccchhhhhh!")
}

func CreateUsers(c *fiber.Ctx) error {

	user, err := xusers.CreateUser(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(user)
}
