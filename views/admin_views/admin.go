package admin_views

import (
	"txrnxp/utils"
	"txrnxp/utils/admin_utils"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the United States of Ganga admins biaaattccchhhhhh!")
}


func CreateAdminUsers(c *fiber.Ctx) error {

	user, err := admin_utils.CreateAdminUser(c)

	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	return utils.CreatedResponse(c, user, "Successfully created admin user")
}

func GetAdminUsers(c *fiber.Ctx) error {
	return admin_utils.GetAdminUsers(c)
}
