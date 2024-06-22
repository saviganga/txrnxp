package auth_utils

import (
	"errors"
	"fmt"
	"os"
	"time"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/utils"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserAuthToken(user *models.Xuser) (string, error) {

	// create JWT and auth tokens
	db := initialisers.ConnectDb().Db
	token := utils.GenerateRandomString(6)
	claims := jwt.MapClaims{
		"email": user.Email,
		"id":    user.Id,
		"token": token,
	}

	adminToken := models.XuserAuthToken{UserId: user.Id, Token: token, ExpiryDate: time.Now().Add(time.Hour * 72)}
	dbError := db.Create(&adminToken).Error

	secret_token := os.Getenv("SECRET_TOKEN")
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := generateToken.SignedString([]byte(secret_token))
	if err != nil || dbError != nil {
		return "", errors.New("unable to create user token")
	}

	return token, nil
}

func ValidateUserEmail(email string) (*models.Xuser, error) {

	db := initialisers.ConnectDb().Db
	user := new(models.Xuser)
	keyObject := "email"
	filter := fmt.Sprintf("%s = ?", keyObject)
	db.Find(&user, filter, email)
	if user.Id == uuid.Nil {
		respMessage := "invalid credentials. validate parameters again"
		return nil, errors.New(respMessage)
	}
	return user, nil

}

func ValidateUserPassword(user *models.Xuser, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Invalid Password:", err)
		return false
	}
	return true
}


