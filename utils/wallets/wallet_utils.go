package wallets

import (
	"errors"
	"fmt"
	"txrnxp/initialisers"
	"txrnxp/models"
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
