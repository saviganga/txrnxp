package auth

import (
	"time"
	"txrnxp/initialisers"
	"txrnxp/serializers/auth"
	"txrnxp/utils"
	"txrnxp/utils/auth_utils"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the United States of Ganga auth biaaattccchhhhhh!")
}

func Login(c *fiber.Ctx) error {

	db := initialisers.ConnectDb().Db

	platform := c.Get("Platform")
	if platform == "" {
		return utils.BadRequestResponse(c, "pass your platform")
	}
	user_request := new(auth.UserLoginSerializer)
	err := c.BodyParser(user_request)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	// validate user email
	user, admin, err := auth_utils.ValidateUserEmail(user_request.Email, platform)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	if user != nil {

		// validate password
		is_password := auth_utils.ValidateUserPassword(user, user_request.Password)
		if !is_password {
			return utils.BadRequestResponse(c, "Invalid credentials")
		}

		// auth token
		token, err := auth_utils.CreateUserAuthToken(user)
		if err != nil {
			return utils.BadRequestResponse(c, err.Error())
		}

		user.LastLogin = time.Now()
		db.Save(&user)

		// return authenticated user and token
		respMessage := "User successfully authenticated"
		return utils.SuccessResponse(c, token, respMessage)

	} else {

		// validate password
		is_password := auth_utils.ValidateAdminUserPassword(admin, user_request.Password)
		if !is_password {
			return utils.BadRequestResponse(c, "Invalid credentials")
		}

		// auth token
		token, err := auth_utils.CreateAdminUserAuthToken(admin)
		if err != nil {
			return utils.BadRequestResponse(c, err.Error())
		}

		admin.LastLogin = time.Now()
		db.Save(&admin)

		// return authenticated user and token
		respMessage := "User successfully authenticated"
		return utils.SuccessResponse(c, token, respMessage)

	}

}
