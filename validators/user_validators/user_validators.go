package user_validators

import (

	"txrnxp/serializers/user_serializers"
	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
)

func ValidateUpdateUserRequestBody(c *fiber.Ctx) error {

	// validate the request body
	body := new(user_serializers.UpdateUserSerializer)

	if err := c.BodyParser(&body); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	c.Locals("body", body)
	return c.Next()


}
