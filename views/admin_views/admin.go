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


func CreateAdminCommissionConfig(c *fiber.Ctx) error {

	config, err := admin_utils.CreateAdminCommission(c)

	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	return utils.CreatedResponse(c, config, "Successfully created admin commission config")
}


func GetAdminCommissionConfig(c *fiber.Ctx) error {
	config, err := admin_utils.GetAdminCommission(c)

	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	return utils.CreatedResponse(c, config, "Successfully fetched admin commission config")
}


func UpdateAdminCommissionConfig(c *fiber.Ctx) error {

	config, err := admin_utils.UpdateAdminCommission(c)

	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	return utils.CreatedResponse(c, config, "Successfully created admin commission config")
}



