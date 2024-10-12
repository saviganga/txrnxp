package wallet_serializers

import (
	"time"
	"txrnxp/models"
	"txrnxp/serializers/user_serializers"

	"github.com/google/uuid"
)

type WalletSerializer struct {
	Id               uuid.UUID                       `json:"id" validate:"required"`
	User             user_serializers.UserSerializer `json:"user" validate:"required"`
	AvailableBalance string                          `json:"available_balance" validate:"required"`
	LedgerBalance    string                          `json:"ledger_balance" validate:"required"`
	CreatedAt        time.Time                       `json:"created_at" validate:"required"`
	UpdatedAt        time.Time                       `json:"updated_at" validate:"required"`
}

type ReadWalletEntrySerializer struct {
	Id          uuid.UUID                       `json:"id" validate:"required"`
	User        user_serializers.UserSerializer `json:"user" validate:"required"`
	Amount      string                          `json:"amount" validate:"required"`
	Reference   string                          `json:"reference" validate:"required"`
	EntryType   string                          `json:"entry_type" validate:"required"`
	Description string                          `json:"description" validate:"required"`
	CreatedAt   time.Time                       `json:"created_at" validate:"required"`
	UpdatedAt   time.Time                       `json:"updated_at" validate:"required"`
}

type WalletManualEntrySerializer struct {
	UserId    string `json:"user_id" validate:"required"`
	Amount    string `json:"amount" validate:"required"`
	EntryType string `json:"entry_type" validate:"required"`
}

type WalletTransferSerializer struct {
	ReceiverEmail string `json:"receiver_email" validate:"required"`
	Amount        string `json:"amount" validate:"required"`
}

type AdminWalletSerializer struct {
	Id               uuid.UUID                       `json:"id" validate:"required"`
	AvailableBalance string                          `json:"available_balance" validate:"required"`
	LedgerBalance    string                          `json:"ledger_balance" validate:"required"`
	CreatedAt        time.Time                       `json:"created_at" validate:"required"`
	UpdatedAt        time.Time                       `json:"updated_at" validate:"required"`
}

type ReadAdminWalletEntrySerializer struct {
	Id          uuid.UUID                       `json:"id" validate:"required"`
	Amount      string                          `json:"amount" validate:"required"`
	Reference   string                          `json:"reference" validate:"required"`
	EntryType   string                          `json:"entry_type" validate:"required"`
	Description string                          `json:"description" validate:"required"`
	CreatedAt   time.Time                       `json:"created_at" validate:"required"`
	UpdatedAt   time.Time                       `json:"updated_at" validate:"required"`
}







func SerializeGetWallets(userwallets []models.UserWallet) []WalletSerializer {

	serialized_wallet := new(WalletSerializer)
	serialized_wallets := []WalletSerializer{}
	serialized_user := new(user_serializers.UserSerializer)

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

	return serialized_wallets

}

func SerializeGetWalletEntries(wallet_tx []models.TransactionEntries) []ReadWalletEntrySerializer {

	serialized_wallet_tx := new(ReadWalletEntrySerializer)
	serialized_wallet_txs := []ReadWalletEntrySerializer{}
	serialized_user := new(user_serializers.UserSerializer)

	for _, tx := range wallet_tx {

		serialized_user.Id = tx.User.Id
		serialized_user.Email = tx.User.Email
		serialized_user.UserName = tx.User.UserName
		serialized_user.FirstName = tx.User.FirstName
		serialized_user.LastName = tx.User.LastName
		serialized_user.PhoneNumber = tx.User.PhoneNumber
		serialized_user.IsActive = tx.User.IsActive
		serialized_user.IsBusiness = tx.User.IsBusiness
		serialized_user.LastLogin = tx.User.LastLogin
		serialized_user.CreatedAt = tx.User.CreatedAt
		serialized_user.UpdatedAt = tx.User.UpdatedAt

		serialized_wallet_tx.Id = tx.Id
		serialized_wallet_tx.User = *serialized_user
		serialized_wallet_tx.Amount = tx.Amount
		serialized_wallet_tx.Reference = tx.Reference
		serialized_wallet_tx.EntryType = tx.EntryType
		serialized_wallet_tx.Description = tx.Description
		serialized_wallet_tx.CreatedAt = tx.CreatedAt
		serialized_wallet_tx.UpdatedAt = tx.UpdatedAt

		serialized_wallet_txs = append(serialized_wallet_txs, *serialized_wallet_tx)
	}

	return serialized_wallet_txs

}

func SerializeGetAdminWallet(wallet *models.AdminWallet) AdminWalletSerializer {

	serialized_wallet := new(AdminWalletSerializer)

	serialized_wallet.Id = wallet.Id
	serialized_wallet.AvailableBalance = wallet.AvailableBalance
	serialized_wallet.LedgerBalance = wallet.LedgerBalance
	serialized_wallet.CreatedAt = wallet.CreatedAt
	serialized_wallet.UpdatedAt = wallet.UpdatedAt


	return *serialized_wallet

}


func SerializeGetAdminWalletEntries(wallet_tx []models.AdminTransactionEntries) []ReadAdminWalletEntrySerializer {

	serialized_wallet_tx := new(ReadAdminWalletEntrySerializer)
	serialized_wallet_txs := []ReadAdminWalletEntrySerializer{}

	for _, tx := range wallet_tx {

		serialized_wallet_tx.Id = tx.Id
		serialized_wallet_tx.Amount = tx.Amount
		serialized_wallet_tx.Reference = tx.Reference
		serialized_wallet_tx.EntryType = tx.EntryType
		serialized_wallet_tx.Description = tx.Description
		serialized_wallet_tx.CreatedAt = tx.CreatedAt
		serialized_wallet_tx.UpdatedAt = tx.UpdatedAt

		serialized_wallet_txs = append(serialized_wallet_txs, *serialized_wallet_tx)
	}

	return serialized_wallet_txs

}
