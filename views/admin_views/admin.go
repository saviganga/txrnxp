package admin_views

import (
	"txrnxp/utils/admin_utils"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the United States of Ganga admins biaaattccchhhhhh!")
}


func CreateAdminUsers(c *fiber.Ctx) error {

	user, err := admin_utils.CreateAdminUser(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(user)
}

func GetAdminUsers(c *fiber.Ctx) error {
	return admin_utils.GetAdminUsers(c)
}
