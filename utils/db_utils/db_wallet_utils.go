package db_utils

import (
	"txrnxp/initialisers"
	"txrnxp/models"
)

func UpdateWallet(wallet *models.UserWallet) (bool, string) {
	db := initialisers.ConnectDb().Db
	err := db.Save(&models.UserWallet{Id: wallet.Id, AvailableBalance: wallet.AvailableBalance, LedgerBalance: wallet.LedgerBalance, UserId: wallet.UserId, CreatedAt: wallet.CreatedAt}).Error
	if err != nil {
		return false, "unable to save wallet"
	}
	return true, "successfully updated wallet"

}
