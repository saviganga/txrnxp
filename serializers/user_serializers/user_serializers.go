package user_serializers

import (
	"time"
	"txrnxp/models"

	"github.com/google/uuid"
)

type ChangePasswordSerializer struct {
	OldPassword	string `json:"old_password" validate:"required"`
	NewPassword	string `json:"new_password" validate:"required"`
	ConfirmPassword	string `json:"confirm_password" validate:"required"`
}

type UserSerializer struct {
	Id          uuid.UUID `json:"id" validate:"required"`
	Email       string    `json:"email" validate:"required"`
	UserName    string    `json:"username" validate:"required"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Image       string    `json:"image"`
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


type UpdateUserSerializer struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
}

func SerializeUser(user models.AdminUser) UserSerializer {

	serialized_user := new(UserSerializer)

	serialized_user.Id = user.Id
	serialized_user.Email = user.Email
	serialized_user.UserName = user.UserName
	serialized_user.FirstName = user.FirstName
	serialized_user.LastName = user.LastName
	serialized_user.PhoneNumber = user.PhoneNumber
	serialized_user.IsActive = user.IsActive
	serialized_user.LastLogin = user.LastLogin
	serialized_user.CreatedAt = user.CreatedAt
	serialized_user.UpdatedAt = user.UpdatedAt

	return *serialized_user
}

func SerializeUsers(users []models.Xuser) []UserSerializer {

	serialized_user := new(UserSerializer)
	serialized_users := []UserSerializer{}

	for _, user := range users {
		serialized_user.Id = user.Id
		serialized_user.Email = user.Email
		serialized_user.UserName = user.UserName
		serialized_user.FirstName = user.FirstName
		serialized_user.LastName = user.LastName
		serialized_user.PhoneNumber = user.PhoneNumber
		serialized_user.IsActive = user.IsActive
		serialized_user.IsBusiness = user.IsBusiness
		serialized_user.LastLogin = user.LastLogin
		serialized_user.CreatedAt = user.CreatedAt
		serialized_user.UpdatedAt = user.UpdatedAt

		serialized_users = append(serialized_users, *serialized_user)
	}

	return serialized_users
}

func SerializeUserSerializer(user models.Xuser) UserSerializer {

	serialized_user := new(UserSerializer)

	serialized_user.Id = user.Id
	serialized_user.Email = user.Email
	serialized_user.UserName = user.UserName
	serialized_user.FirstName = user.FirstName
	serialized_user.LastName = user.LastName
	serialized_user.PhoneNumber = user.PhoneNumber
	serialized_user.Image = user.Image
	serialized_user.IsActive = user.IsActive
	serialized_user.IsBusiness = user.IsBusiness
	serialized_user.LastLogin = user.LastLogin
	serialized_user.CreatedAt = user.CreatedAt
	serialized_user.UpdatedAt = user.UpdatedAt

	return *serialized_user
}
