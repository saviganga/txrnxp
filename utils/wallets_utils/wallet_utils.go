package wallets_utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/wallet_serializers"
	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateUserWallet(user *models.Xuser) error {
	// create user wallet
	db := initialisers.ConnectDb().Db
	userwallet_query := models.UserWallet{UserId: user.Id}
	dbError := db.Create(&userwallet_query).Error
	if dbError != nil {
		return errors.New("oops! error creating user wallet")
	}
	return nil
}

func GetUserWallets(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	userwallets := []models.UserWallet{}
	privilege := authenticated_user["privilege"]
	if privilege == "ADMIN" {
		db.Model(&models.UserWallet{}).Joins("User").Find(&userwallets)
	} else {
		db.Model(&models.UserWallet{}).Joins("User").First(&userwallets, "user_wallets.user_id = ?", authenticated_user["id"])
	}
	return c.Status(200).JSON(userwallets)

}

func GetUserWalletTransactions(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	wallet_tx := []models.TransactionEntries{}
	privilege := authenticated_user["privilege"]
	if privilege == "ADMIN" {
		db.Model(&models.TransactionEntries{}).Joins("User").Find(&wallet_tx)
	} else {
		db.Model(&models.TransactionEntries{}).Joins("User").First(&wallet_tx, "user_wallets.user_id = ?", authenticated_user["id"])
	}
	return c.Status(200).JSON(wallet_tx)

}

func AdminWalletManualEntry(c *fiber.Ctx) (bool, string) {
	// initialise niggas
	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	userwallets := []models.UserWallet{}

	privilege := authenticated_user["privilege"].(string)

	if strings.ToUpper(privilege) != "ADMIN" {
		return false, "Oops! this feature is only available for admins"
	}

	entry_request := new(wallet_serializers.WalletManualEntrySerializer)
	err := c.BodyParser(entry_request)
	if err != nil {
		return false, err.Error()
	}

	// get the user wallet
	db.Model(&models.UserWallet{}).Joins("User").First(&userwallets, "user_wallets.user_id = ?", entry_request.UserId)
	userwallet := userwallets[0]
	fmt.Println(userwallet.UserId)
	// fmt.Println(userwallet.UserId)

	old_available_balance, err := strconv.ParseFloat(userwallets[0].AvailableBalance, 64)
	if err != nil {
		return false, err.Error()
	}
	old_ledger_balance, err := strconv.ParseFloat(userwallets[0].LedgerBalance, 64)
	if err != nil {
		return false, err.Error()
	}

	entry_request_amount_float, err := strconv.ParseFloat(entry_request.Amount, 64)
	if err != nil {
		return false, err.Error()
	}

	user_id_uuid, err := utils.ConvertStringToUUID(entry_request.UserId)
	if err != nil {
		return false, err.Error()
	}

	if strings.ToUpper(entry_request.EntryType) == "CREDIT" {

		new_available_balance := strconv.FormatFloat(old_available_balance + entry_request_amount_float, 'f', -1, 64)
		new_ledger_balance := strconv.FormatFloat(old_ledger_balance + entry_request_amount_float, 'f', -1, 64)
		entry_description := "Manual entry - credit"
		err = db.Save(&models.UserWallet{Id: userwallet.Id, AvailableBalance: new_available_balance, LedgerBalance: new_ledger_balance, UserId: user_id_uuid}).Error
		if err != nil {
			return false, err.Error()
		}

		// update wallet transaction
		wallet_tx := models.TransactionEntries{UserId: user_id_uuid, Amount: entry_request.Amount, Description: entry_description}
		dbError := db.Create(&wallet_tx).Error
		if dbError != nil {
			return false, dbError.Error()
		}

		return true, "Successful"

	} else {

		new_available_balance := strconv.FormatFloat(old_available_balance - entry_request_amount_float, 'f', -1, 64)
		new_ledger_balance := strconv.FormatFloat(old_ledger_balance - entry_request_amount_float, 'f', -1, 64)
		entry_description := "Manual entry - debit"
		err = db.Save(&models.UserWallet{Id: userwallet.Id, AvailableBalance: new_available_balance, LedgerBalance: new_ledger_balance, UserId: user_id_uuid}).Error
		if err != nil {
			return false, err.Error()
		}

		// update wallet transaction
		wallet_tx := models.TransactionEntries{UserId: user_id_uuid, Amount: entry_request.Amount, Description: entry_description}
		dbError := db.Create(&wallet_tx).Error
		if dbError != nil {
			return false, dbError.Error()
		}
		
		return true, "Successful"
	}

}
