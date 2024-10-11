package xusers

import (
	"txrnxp/utils"
	"txrnxp/utils/xusers_utils"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the United States of Ganga users biaaattccchhhhhh!")
}

func CreateUsers(c *fiber.Ctx) error {

	user, err := xusers_utils.CreateUser(c)

	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	return utils.CreatedResponse(c, user, "Successfully created user")
}

func GetUsers(c *fiber.Ctx) error {
	return xusers_utils.GetUsers(c)
}

func UploadUserImage(c *fiber.Ctx) error {
	return xusers_utils.UploadUserImage(c)
}
