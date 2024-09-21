package wallet_serializers

import (
	"time"
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
	Id               uuid.UUID                       `json:"id" validate:"required"`
	User             user_serializers.UserSerializer `json:"user" validate:"required"`
	Reference string                          `json:"reference" validate:"required"`
	EntryType    string                          `json:"entry_type" validate:"required"`
	Description    string                          `json:"description" validate:"required"`
	CreatedAt        time.Time                       `json:"created_at" validate:"required"`
	UpdatedAt        time.Time                       `json:"updated_at" validate:"required"`
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
