package wallets_utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/user_serializers"
	"txrnxp/serializers/wallet_serializers"
	"txrnxp/utils"
	"txrnxp/utils/db_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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
	serialized_wallet := new(wallet_serializers.WalletSerializer)
	serialized_wallets := []wallet_serializers.WalletSerializer{}
	serialized_user := new(user_serializers.UserSerializer)
	privilege := authenticated_user["privilege"]
	message := ""
	if privilege == "ADMIN" {
		db.Model(&models.UserWallet{}).Joins("User").Order("created_at desc").Find(&userwallets)
		message = "Succesfully fetched wallets"
	} else {
		db.Model(&models.UserWallet{}).Joins("User").Order("created_at desc").First(&userwallets, "user_wallets.user_id = ?", authenticated_user["id"])
		message = "Succesfully fetched wallet"
	}
	for _, wallet := range userwallets {

		serialized_user.Id = wallet.User.Id
		serialized_user.Email = wallet.User.Email
		serialized_user.UserName = wallet.User.UserName
		serialized_user.FirstName = wallet.User.FirstName
		serialized_user.LastName = wallet.User.LastName
		serialized_user.PhoneNumber = wallet.User.PhoneNumber
		serialized_user.IsActive = wallet.User.IsActive
		serialized_user.IsBusiness = wallet.User.IsBusiness
		serialized_user.LastLogin = wallet.User.LastLogin
		serialized_user.CreatedAt = wallet.User.CreatedAt
		serialized_user.UpdatedAt = wallet.User.UpdatedAt

		serialized_wallet.Id = wallet.Id
		serialized_wallet.User = *serialized_user
		serialized_wallet.AvailableBalance = wallet.AvailableBalance
		serialized_wallet.LedgerBalance = wallet.LedgerBalance
		serialized_wallet.CreatedAt = wallet.CreatedAt
		serialized_wallet.UpdatedAt = wallet.UpdatedAt

		serialized_wallets = append(serialized_wallets, *serialized_wallet)

	}
	return utils.SuccessResponse(c, serialized_wallets, message)

}

func GetUserWalletTransactions(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	wallet_tx := []models.TransactionEntries{}
	privilege := authenticated_user["privilege"]
	if privilege == "ADMIN" {
		db.Model(&models.TransactionEntries{}).Joins("User").Find(&wallet_tx).Order("created_at DESC")
	} else {
		db.Where("user_id = ?", authenticated_user["id"]).Joins("User").Order("created_at DESC").Find(&wallet_tx)
	}
	return utils.SuccessResponse(c, wallet_tx, "Successfully fetched wallet transactions")

}

func AdminWalletManualEntry(c *fiber.Ctx) (bool, string) {
	// mobilize men !
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	entry_request := new(wallet_serializers.WalletManualEntrySerializer)
	privilege := authenticated_user["privilege"].(string)

	if strings.ToUpper(privilege) != "ADMIN" {
		return false, "Oops! this feature is only available for admins"
	}

	err := c.BodyParser(entry_request)
	if err != nil {
		return false, err.Error()
	}

	user_id_uuid, err := utils.ConvertStringToUUID(entry_request.UserId)
	if err != nil {
		return false, err.Error()
	}

	if strings.ToUpper(entry_request.EntryType) == "CREDIT" {

		entry_description := "Manual entry - credit"
		is_credited, credited_wallet := CreditUserWallet(user_id_uuid, entry_request.Amount, entry_description)
		if !is_credited {
			return false, credited_wallet
		}

		return true, credited_wallet

	} else {

		amount_float, err := utils.ConvertStringToFloat(entry_request.Amount)
		if err != nil || amount_float == 0.0 {
			return false, err.Error()
		}

		entry_description := "Manual entry - debit"
		is_debited, debited_wallet := DebitUserWallet(user_id_uuid, amount_float, entry_description)
		if !is_debited {
			return false, debited_wallet
		}

		return true, debited_wallet
	}

}

func DebitUserWallet(user_id uuid.UUID, amount float64, description string) (bool, string) {

	db := initialisers.ConnectDb().Db
	userwallets := []models.UserWallet{}

	// get user wallet
	err := db.Model(&models.UserWallet{}).Joins("User").First(&userwallets, "user_wallets.user_id = ?", user_id).Error
	if err != nil {
		return false, "Oops! Unable to find user wallet"
	}

	// convert wallet balance to float
	userWallet_available_balance, err := strconv.ParseFloat(userwallets[0].AvailableBalance, 64)
	if err != nil {
		return false, "error converting wallet available balance"
	}
	userWallet_ledger_balance, err := strconv.ParseFloat(userwallets[0].LedgerBalance, 64)
	if err != nil {
		return false, "error converting wallet available balance"
	}

	if userWallet_available_balance < amount {
		return false, "oops! insufficient wallet funds"
	}

	// debit user wallet
	userwallets[0].AvailableBalance = strconv.FormatFloat(userWallet_available_balance-amount, 'f', -1, 64)
	userwallets[0].LedgerBalance = strconv.FormatFloat(userWallet_ledger_balance-amount, 'f', -1, 64)
	is_debited, debited_wallet := db_utils.UpdateWallet(&userwallets[0])
	if !is_debited {
		return false, debited_wallet
	}

	amount_str := fmt.Sprintf("%.2f", amount)
	// update wallet transaction
	wallet_tx := models.TransactionEntries{UserId: user_id, Amount: amount_str, Description: description, EntryType: "DEBIT"}
	dbError := db.Create(&wallet_tx).Error
	if dbError != nil {
		return false, dbError.Error()
	}
	return true, debited_wallet

}

func CreditUserWallet(user_id uuid.UUID, amount string, description string) (bool, string) {

	db := initialisers.ConnectDb().Db
	userwallets := []models.UserWallet{}

	// convert amount to float
	amount_float, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		return false, "error converting event ticket price"
	}

	// get user wallet
	err = db.Model(&models.UserWallet{}).Joins("User").First(&userwallets, "user_wallets.user_id = ?", user_id).Error
	if err != nil {
		return false, "Oops! Unable to find user wallet"
	}

	// convert wallet balance to float
	userWallet_available_balance, err := strconv.ParseFloat(userwallets[0].AvailableBalance, 64)
	if err != nil {
		return false, "error converting wallet available balance"
	}
	userWallet_ledger_balance, err := strconv.ParseFloat(userwallets[0].LedgerBalance, 64)
	if err != nil {
		return false, "error converting wallet available balance"
	}

	// credit user wallet
	userwallets[0].AvailableBalance = strconv.FormatFloat(userWallet_available_balance+amount_float, 'f', -1, 64)
	userwallets[0].LedgerBalance = strconv.FormatFloat(userWallet_ledger_balance+amount_float, 'f', -1, 64)
	is_credited, credited_wallet := db_utils.UpdateWallet(&userwallets[0])
	if !is_credited {
		return false, credited_wallet
	}
	// update wallet transaction
	wallet_tx := models.TransactionEntries{UserId: user_id, Amount: amount, Description: description, EntryType: "CREDIT"}
	dbError := db.Create(&wallet_tx).Error
	if dbError != nil {
		return false, dbError.Error()
	}
	return true, credited_wallet

}

func WalletTransfer(c *fiber.Ctx) (bool, string) {
	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	transfer_request := new(wallet_serializers.WalletTransferSerializer)
	privilege := authenticated_user["privilege"].(string)
	users := []models.Xuser{}
	sender_email := authenticated_user["email"]

	if strings.ToUpper(privilege) != "USER" {
		return false, "Oops! this feature is only available for users"
	}

	err := c.BodyParser(transfer_request)
	if err != nil {
		return false, err.Error()
	}

	if sender_email == strings.ToLower(transfer_request.ReceiverEmail) {
		return false, "Oops! You cannot transfer funds to youeself"
	}

	sender_id_string := authenticated_user["id"].(string)
	sender_id_uuid, err := utils.ConvertStringToUUID(sender_id_string)
	if err != nil {
		return false, err.Error()
	}

	db.First(&users, "email = ?", transfer_request.ReceiverEmail)
	receiver_id_uuid := users[0].Id

	amount_float, err := utils.ConvertStringToFloat(transfer_request.Amount)
	if err != nil || amount_float == 0.0 {
		return false, err.Error()
	}

	debit_entry_description := fmt.Sprintf("Wallet transfer of %.2f to %s", amount_float, transfer_request.ReceiverEmail)
	is_debited, debited_wallet := DebitUserWallet(sender_id_uuid, amount_float, debit_entry_description)
	if !is_debited {
		return false, debited_wallet
	}

	credit_entry_description := fmt.Sprintf("Wallet transfer of %.2f from %s", amount_float, users[0].Email)
	is_credited, credited_wallet := CreditUserWallet(receiver_id_uuid, transfer_request.Amount, credit_entry_description)
	if !is_credited {
		return false, credited_wallet
	}

	return true, "Wallet transfer successful"

}
