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
	serialized_user := new(user_serializers.UserSerializer)
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

	serialized_user.Id = user.Id
	serialized_user.Email = user.Email
	serialized_user.UserName = user.UserName
	serialized_user.FirstName = user.FirstName
	serialized_user.LastName = user.LastName
	serialized_user.PhoneNumber = user.PhoneNumber
	serialized_user.IsActive = user.IsActive
	serialized_user.IsBusiness = user.IsBusiness
	serialized_user.LastLogin= user.LastLogin
	serialized_user.CreatedAt = user.CreatedAt
	serialized_user.UpdatedAt = user.UpdatedAt

	return serialized_user, nil
}

func GetUsers(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	users := []models.Xuser{}
	serialized_user := new(user_serializers.UserSerializer)
	serialized_users := []user_serializers.UserSerializer{}
	privilege := authenticated_user["privilege"]
	if privilege == "ADMIN" {
		db.Order("created_at desc").Find(&users)
	} else {
		db.Order("created_at desc").First(&users, "id = ?", authenticated_user["id"])
	}
	for _, user := range users {
		serialized_user.Id = user.Id
		serialized_user.Email = user.Email
		serialized_user.UserName = user.UserName
		serialized_user.FirstName = user.FirstName
		serialized_user.LastName = user.LastName
		serialized_user.PhoneNumber = user.PhoneNumber
		serialized_user.IsActive = user.IsActive
		serialized_user.IsBusiness = user.IsBusiness
		serialized_user.LastLogin= user.LastLogin
		serialized_user.CreatedAt = user.CreatedAt
		serialized_user.UpdatedAt = user.UpdatedAt

		serialized_users = append(serialized_users, *serialized_user)
	}
	return utils.SuccessResponse(c, serialized_users, "Successfully fetched users")
}
