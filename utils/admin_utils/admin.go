package admin_utils

import (
	"errors"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/user_serializers"
	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateAdminUser(c *fiber.Ctx) (*models.AdminUser, error) {

	db := initialisers.ConnectDb().Db
	user := new(models.AdminUser)
	serialized_user := new(user_serializers.UserSerializer) 
	err := c.BodyParser(user)
	if err != nil {
		return nil, errors.New("invalid request body")
	}
	err = db.Create(&user).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	serialized_user.Id = user.Id
	serialized_user.Email = user.Email
	serialized_user.UserName = user.UserName
	serialized_user.FirstName = user.FirstName
	serialized_user.LastName = user.LastName
	serialized_user.PhoneNumber = user.PhoneNumber
	serialized_user.IsActive = user.IsActive
	serialized_user.LastLogin= user.LastLogin
	serialized_user.CreatedAt = user.CreatedAt
	serialized_user.UpdatedAt = user.UpdatedAt

	return user, nil
}

func GetAdminUsers(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	users := []models.AdminUser{}
	serialized_user := new(user_serializers.UserSerializer) 
	db.First(&users, "id = ?", authenticated_user["id"])
	
	serialized_user.Id = users[0].Id
	serialized_user.Email = users[0].Email
	serialized_user.UserName = users[0].UserName
	serialized_user.FirstName = users[0].FirstName
	serialized_user.LastName = users[0].LastName
	serialized_user.PhoneNumber = users[0].PhoneNumber
	serialized_user.IsActive = users[0].IsActive
	serialized_user.LastLogin= users[0].LastLogin
	serialized_user.CreatedAt = users[0].CreatedAt
	serialized_user.UpdatedAt = users[0].UpdatedAt
	
	return utils.SuccessResponse(c, serialized_user, "Successfully fetched users")
}
