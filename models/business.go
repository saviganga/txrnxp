package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Business struct {
	Id               uuid.UUID    `gorm:"type:uuid;primaryKey;not null" json:"id"`
	UserId           uuid.UUID    `gorm:"type:uuid;not null" json:"user_id"`
	Name string       `gorm:"type:varchar(50)" json:"name"`
	Country    string       `gorm:"type:varchar(50)" json:"country"`
	CreatedAt        time.Time    `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt        time.Time    `gorm:"type:timestamp with time zone" json:"updated_at"`
	User             Xuser `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

func (business *Business) BeforeCreate(*gorm.DB) (err error) {

	business.Id = uuid.New()
	business.Country = "NGA"
	return
}