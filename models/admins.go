package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminUser struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`
	Email       string    `gorm:"type:varchar(50);not null;unique" json:"email"`
	Password    string    `gorm:"type:varchar(128);not null" json:"password"`
	FirstName   string    `gorm:"type:varchar(50)" json:"first_name"`
	LastName    string    `gorm:"type:varchar(50)" json:"last_name"`
	UserName    string    `gorm:"type:varchar(50);not null;unique" json:"username" `
	PhoneNumber string    `gorm:"type:varchar(15)" json:"phone_number"`
	IsActive    bool      `gorm:"type:boolean;default:true" json:"is_active"`
	IsVerified  bool      `gorm:"type:boolean;default:false" json:"is_verified"`
	LastLogin   time.Time `gorm:"type:timestamp with time zone" json:"last_login"`
	CreatedAt   time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
}

func (user *AdminUser) BeforeCreate(*gorm.DB) (err error) {
	if len(user.Email) == 0 {
		return errors.New("email cannot be empty")
	}
	if len(user.Password) == 0 {
		return errors.New("password cannot be empty")
	}
	if len(user.UserName) == 0 {
		return errors.New("username cannot be empty")
	}
	pass := []byte(user.Password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	user.Id = uuid.New()
	user.Password = string(hash)
	return
}

type AdminUserAuthToken struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`
	UserId     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Token      string    `gorm:"type:varchar(50);not null;" json:"token"`
	ExpiryDate time.Time `gorm:"type:timestamp with time zone" json:"expiry_date"`
	User       AdminUser `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE"`
}

func (authToken *AdminUserAuthToken) BeforeCreate(*gorm.DB) (err error) {
	authToken.Id = uuid.New()
	return
}

type AdminCommissionConfig struct {
	Id         uuid.UUID `gorm:"type:uuid;primaryKey;not null" json:"id"`
	Type       string    `gorm:"type:varchar(50);unique" json:"type"`
	Commission string    `gorm:"type:numeric(10,2);not null;default:0.00" json:"commission"`
	Cap        string    `gorm:"type:numeric(10,2);not null;default:0.00" json:"cap"`
	CreatedAt  time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at"`
	UpdatedAt  time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
}

func (tx *AdminCommissionConfig) BeforeCreate(*gorm.DB) (err error) {

	tx.Id = uuid.New()
	return
}
