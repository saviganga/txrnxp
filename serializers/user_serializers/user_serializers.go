package user_serializers

import (
	"time"

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
	IsBusiness  bool      `json:"is_business" validate:"required"`
	LastLogin   time.Time `json:"last_login" validate:"required"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`
}
