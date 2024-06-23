package admin_utils

import (
	"errors"
	"txrnxp/initialisers"
	"txrnxp/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateAdminUser(c *fiber.Ctx) (*models.AdminUser, error) {

	db := initialisers.ConnectDb().Db
	user := new(models.AdminUser)
	err := c.BodyParser(user)
	if err != nil {
		return nil, errors.New("invalid request body")
	}
	err = db.Create(&user).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return user, nil
}

func GetAdminUsers(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	users := []models.AdminUser{}
	db.First(&users, "id = ?", authenticated_user["id"])
	return c.Status(200).JSON(users)
}
