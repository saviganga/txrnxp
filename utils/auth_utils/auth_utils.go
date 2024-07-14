package auth_utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserAuthToken(user *models.Xuser) (string, error) {

	// create JWT and auth tokens
	db := initialisers.ConnectDb().Db
	token := utils.GenerateRandomString(6)
	claims := jwt.MapClaims{
		"email":     user.Email,
		"id":        user.Id,
		"token":     token,
		"privilege": "USER",
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

func CreateAdminUserAuthToken(adminUser *models.AdminUser) (string, error) {

	// create JWT and auth tokens
	db := initialisers.ConnectDb().Db
	token := utils.GenerateRandomString(6)
	claims := jwt.MapClaims{
		"email":     adminUser.Email,
		"id":        adminUser.Id,
		"token":     token,
		"privilege": "ADMIN",
	}

	adminToken := models.AdminUserAuthToken{UserId: adminUser.Id, Token: token, ExpiryDate: time.Now().Add(time.Hour * 72)}
	dbError := db.Create(&adminToken).Error

	secret_token := os.Getenv("SECRET_TOKEN")
	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := generateToken.SignedString([]byte(secret_token))
	if err != nil || dbError != nil {
		return "", errors.New("unable to create user token")
	}

	return token, nil
}

func ValidateUserEmail(email string, platform string) (*models.Xuser, *models.AdminUser, error) {

	db := initialisers.ConnectDb().Db
	keyObject := "email"
	filter := fmt.Sprintf("%s = ?", keyObject)
	if platform == "ADMIN" {
		user := new(models.AdminUser)
		db.Find(&user, filter, email)
		if user.Id == uuid.Nil {
			respMessage := "invalid admin credentials. validate parameters again"
			return nil, nil, errors.New(respMessage)
		}
		return nil, user, nil
	} else {
		user := new(models.Xuser)
		db.Find(&user, filter, email)
		if user.Id == uuid.Nil {
			respMessage := "invalid credentials. validate parameters again"
			return nil, nil, errors.New(respMessage)
		}
		return user, nil, nil
	}

}

func ValidateUserPassword(user *models.Xuser, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Invalid Password:", err)
		return false
	}
	return true
}

func ValidateAdminUserPassword(user *models.AdminUser, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Invalid Password:", err)
		return false
	}
	return true
}

func ValidateAuth(ctx *fiber.Ctx) error {
	authorization := ctx.Get("Authorization")
	if authorization == "" {
		respMessage := "oops! please pass in authorization header"
		return errors.New(respMessage)
	}

	jwtSecret := os.Getenv("SECRET_TOKEN")
	tokenString := JWTSplit(authorization)
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		respMessage := "oops! authentication error"
		return errors.New(respMessage)
	}

	privilege := claims["privilege"]
	if privilege == "USER" {
		var userAuth []models.XuserAuthToken
		db := initialisers.ConnectDb().Db
		db.Where("token = ? AND user_id = ?", claims["token"], claims["id"]).Find(&userAuth)
		if len(userAuth) == 0 || userAuth[0].ExpiryDate.String() < time.Now().String() {
			respMessage := "oop! invalid token. please log in"
			return errors.New(respMessage)
		}
	} else {
		var userAuth []models.AdminUserAuthToken
		db := initialisers.ConnectDb().Db
		db.Where("token = ? AND user_id = ?", claims["token"], claims["id"]).Find(&userAuth)
		if len(userAuth) == 0 || userAuth[0].ExpiryDate.String() < time.Now().String() {
			respMessage := "oop! invalid token. please log in"
			return errors.New(respMessage)
		}
	}

	ctx.Locals("user", claims)
	return ctx.Next()
}

func JWTSplit(headerToken string) string {
	if headerToken == "" {
		return ""
	}

	parts := strings.Split(headerToken, "JWT")
	if len(parts) != 2 {
		return ""
	}

	token := strings.TrimSpace(parts[1])
	if len(token) < 1 {
		return ""
	}
	return token
}
