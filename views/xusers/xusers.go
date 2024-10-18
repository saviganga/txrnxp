package xusers

import (
	"txrnxp/initialisers"
	"txrnxp/serializers/user_serializers"
	"txrnxp/utils"
	"txrnxp/utils/auth_utils"
	"txrnxp/utils/xusers_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
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

func GetUser(c *fiber.Ctx) error {
	return xusers_utils.GetUser(c)
}

func UpdateUser(c *fiber.Ctx) error {
	return xusers_utils.UpdateUser(c)
}

func UploadUserImage(c *fiber.Ctx) error {
	return xusers_utils.UploadUserImage(c)
}

func ChangePassword(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	user_request := new(user_serializers.ChangePasswordSerializer)
	err := c.BodyParser(user_request)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	privilege := authenticated_user["privilege"]

	// validate user email
	user, admin, err := auth_utils.ValidateUserEmail(authenticated_user["email"].(string), privilege.(string))
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	if user_request.NewPassword != user_request.ConfirmPassword {
		return utils.BadRequestResponse(c, "Password and Confirm Password mismatch")
	}

	if user != nil {
		is_password := auth_utils.ValidateUserPassword(user, user_request.OldPassword)
		if !is_password {
			return utils.BadRequestResponse(c, "Invalid credentials")
		}
		pass := []byte(user_request.NewPassword)
		hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		user.Password = string(hash)
		db.Save(&user)

		// return authenticated user and token
		respMessage := "User Password successfully updated"
		return utils.NoDataSuccessResponse(c, respMessage)
	} else {
		// validate password
		is_password := auth_utils.ValidateAdminUserPassword(admin, user_request.NewPassword)
		if !is_password {
			return utils.BadRequestResponse(c, "Invalid credentials")
		}
		pass := []byte(user_request.NewPassword)
		hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		admin.Password = string(hash)

		db.Save(&admin)

		/// return authenticated user and token
		respMessage := "Admin Password successfully updated"
		return utils.NoDataSuccessResponse(c, respMessage)

	}
}
