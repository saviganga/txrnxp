package user_serializers

import (
	"time"
	"txrnxp/models"

	"github.com/google/uuid"
)

type UserSerializer struct {
	Id          uuid.UUID `json:"id" validate:"required"`
	Email       string    `json:"email" validate:"required"`
	UserName    string    `json:"username" validate:"required"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	IsActive    bool      `json:"is_active" validate:"required"`
	IsBusiness  bool      `json:"is_business"`
	LastLogin   time.Time `json:"last_login" validate:"required"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`
}


type ExportUserSerializer struct {
	Id          uuid.UUID `json:"id" validate:"required"`
	Email       string    `json:"email" validate:"required"`
	UserName    string    `json:"username" validate:"required"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
}

func SerializeUser(user models.AdminUser) (UserSerializer) {

	serialized_user := new(UserSerializer) 

	serialized_user.Id = user.Id
	serialized_user.Email = user.Email
	serialized_user.UserName = user.UserName
	serialized_user.FirstName = user.FirstName
	serialized_user.LastName = user.LastName
	serialized_user.PhoneNumber = user.PhoneNumber
	serialized_user.IsActive = user.IsActive
	serialized_user.LastLogin= user.LastLogin
	serialized_user.CreatedAt = user.CreatedAt
	serialized_user.UpdatedAt = user.UpdatedAt

	return *serialized_user
}