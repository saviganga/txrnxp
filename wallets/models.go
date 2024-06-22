package wallets

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserWallet struct {
	Id               uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`
	AvailableBalance string    `gorm:"type:numeric(10,2);not null;default:0.00"`
	LedgerBalance    string    `gorm:"type:numeric(10,2);not null;default:0.00"`
	CreatedAt        time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
}

func (wallet *UserWallet) BeforeCreate(*gorm.DB) (err error) {

	wallet.Id = uuid.New()
	return
}
