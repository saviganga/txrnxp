package auth

import (
	"txrnxp/serializers/auth"
	"txrnxp/utils/auth_utils"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the United States of Ganga auth biaaattccchhhhhh!")
}

func Login(c *fiber.Ctx) error {

	user_request := new(auth.UserLoginSerializer)
	err := c.BodyParser(user_request)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// validate user email
	user, err := auth_utils.ValidateUserEmail(user_request.Email)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// validate password
	is_password := auth_utils.ValidateUserPassword(user, user_request.Password)
	if !is_password {
		return c.Status(400).JSON(fiber.Map{
			"message": "invalid credentials",
		})
	}

	// auth token
	token, err := auth_utils.CreateUserAuthToken(user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// return authenticated user and token
	respMessage := "User successfully authenticated"
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": respMessage,
		"token":   token,
	})
}
