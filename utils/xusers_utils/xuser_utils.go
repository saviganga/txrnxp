package xusers_utils

import (
	"errors"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/user_serializers"
	"txrnxp/utils"
	"txrnxp/utils/wallets_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateUser(c *fiber.Ctx) (*user_serializers.UserSerializer, error) {

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

	serialized_user := user_serializers.SerializeUserSerializer(*user)

	return &serialized_user, nil
}

func GetUsers(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	users := []models.Xuser{}
	privilege := authenticated_user["privilege"]
	if privilege == "ADMIN" {
		db.Order("created_at desc").Find(&users)
	} else {
		db.Order("created_at desc").First(&users, "id = ?", authenticated_user["id"])
	}
	serialized_users := user_serializers.SerializeUsers(users)
	return utils.SuccessResponse(c, serialized_users, "Successfully fetched users")
}
