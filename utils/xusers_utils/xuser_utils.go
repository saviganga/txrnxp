package xusers_utils

import (
	"errors"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/utils"
	"txrnxp/utils/wallets_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateUser(c *fiber.Ctx) (*models.Xuser, error) {

	db := initialisers.ConnectDb().Db
	user := new(models.Xuser)
	err := c.BodyParser(user)
	if err != nil {
		return nil, errors.New("invalid request body")
	}
	err = db.Create(&user).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = wallets_utils.CreateUserWallet(user)
	if err != nil {
		return nil, errors.New("unable to create user wallet")
	}

	return user, nil
}

func GetUsers(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	users := []models.Xuser{}
	privilege := authenticated_user["privilege"]
	if privilege == "ADMIN" {
		db.Find(&users)
	} else {
		db.First(&users, "id = ?", authenticated_user["id"])
	}
	return utils.SuccessResponse(c, users, "Successfully fetched users")
}
