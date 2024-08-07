package wallet_serializers

type WalletManualEntrySerializer struct {
	UserId string `json:"user_id" validate:"required"`
	Amount string `json:"amount" validate:"required"`
	EntryType string `json:"entry_type" validate:"required"`
}

type WalletTransferSerializer struct {
	ReceiverEmail string `json:"receiver_email" validate:"required"`
	Amount string `json:"amount" validate:"required"`
}