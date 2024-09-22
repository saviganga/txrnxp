package business_serializers

import (
	"time"
	"txrnxp/serializers/user_serializers"

	"github.com/google/uuid"
)

type ReadBusinessSerializer struct {
	Id        uuid.UUID                       `json:"id" validate:"required"`
	User      user_serializers.UserSerializer `json:"user" validate:"required"`
	Reference string                          `json:"reference" validate:"required"`
	Name      string                          `json:"name" validate:"required"`
	Country   string                          `json:"country" validate:"required"`
	CreatedAt time.Time                       `json:"created_at" validate:"required"`
	UpdatedAt time.Time                       `json:"updated_at" validate:"required"`
}


type ReadCreateBusinessSerializer struct {
	Id        uuid.UUID                       `json:"id" validate:"required"`
	Reference string                          `json:"reference" validate:"required"`
	Name      string                          `json:"name" validate:"required"`
	Country   string                          `json:"country" validate:"required"`
	CreatedAt time.Time                       `json:"created_at" validate:"required"`
	UpdatedAt time.Time                       `json:"updated_at" validate:"required"`
}
