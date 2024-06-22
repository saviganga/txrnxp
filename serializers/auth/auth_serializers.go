package auth

type UserLoginSerializer struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}