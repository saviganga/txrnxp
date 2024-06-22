package wallets_utils

import (
	"errors"
	"fmt"
	"txrnxp/initialisers"
	"txrnxp/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateUserWallet(user *models.Xuser) error {
	// create user wallet
	db := initialisers.ConnectDb().Db
	userwallet_query := models.UserWallet{UserId: user.Id}
	dbError := db.Create(&userwallet_query).Error
	if dbError != nil {
		fmt.Println(dbError)
		return errors.New("oops! error creating user wallet")
	}
	return nil
}

func GetUserWallets(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	userwallets := []models.UserWallet{}
	db.First(&userwallets, "user_id = ?", authenticated_user["id"])
	return c.Status(200).JSON(userwallets)

}
