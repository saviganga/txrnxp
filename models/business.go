package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"txrnxp/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Business struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`
	UserId    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Name      string    `gorm:"type:varchar(50)" json:"name"`
	Image     string    `gorm:"type:varchar(150)" json:"image"`
	Reference string    `gorm:"type:varchar(50);not null" json:"reference"`
	Country   string    `gorm:"type:varchar(50)" json:"country"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
	User      Xuser     `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

func (business *Business) BeforeCreate(*gorm.DB) (err error) {

	
	business.Id = uuid.New()
	business.Country = "NGA"
	if len(business.Name) < 3 {
		return errors.New("oops! business name must be greater than 2")
	}
	business.Reference = fmt.Sprintf("%s-%s", strings.ToUpper(business.Name[:3]), utils.GenerateRandomString(4))
	return
}

type BusinessMember struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`
	UserId     uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_business" json:"user_id"`
	BusinessId uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_user_business" json:"business_id"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
	User      Xuser     `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
	Business  Business  `gorm:"foreignKey:BusinessId;constraint:OnDelete:CASCADE"`
}

func (business_member *BusinessMember) BeforeCreate(*gorm.DB) (err error) {

	business_member.Id = uuid.New()
	return
}
