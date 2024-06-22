package xusers

import (
	"errors"
	"fmt"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/utils/wallets"

	"github.com/gofiber/fiber/v2"
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
		fmt.Println(err.Error())
		return nil, errors.New(err.Error())
	}

	err = wallets.CreateUserWallet(user)
	if err != nil {
		return nil, errors.New("unable to create user wallet")
	}

	return user, nil
}
