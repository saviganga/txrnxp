package models

import (
	"time"
	"txrnxp/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserWallet struct {
	Id               uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`
	UserId           uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	AvailableBalance string    `gorm:"type:numeric(10,2);not null;default:0.00"`
	LedgerBalance    string    `gorm:"type:numeric(10,2);not null;default:0.00"`
	CreatedAt        time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt        time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
	User             Xuser     `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

func (wallet *UserWallet) BeforeCreate(*gorm.DB) (err error) {

	wallet.Id = uuid.New()
	return
}

type TransactionEntries struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`
	Reference   string    `gorm:"type:varchar(50);unique" json:"reference"`
	UserId      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Amount      string    `gorm:"type:numeric(10,2);not null;default:0.00"`
	Description string    `gorm:"type:varchar(200)" json:"description"`
	CreatedAt   time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
	User        Xuser     `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

func (tx *TransactionEntries) BeforeCreate(*gorm.DB) (err error) {

	tx.Id = uuid.New()
	tx.Reference = utils.CreateEventReference()
	return
}
