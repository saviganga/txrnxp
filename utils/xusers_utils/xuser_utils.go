package xusers_utils

import (
	"errors"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/user_serializers"
	"txrnxp/utils"
	"txrnxp/utils/wallets_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateUser(c *fiber.Ctx) (*user_serializers.UserSerializer, error) {

	db := initialisers.ConnectDb().Db
	user := new(models.Xuser)
	err := c.BodyParser(user)
	if err != nil {
		return nil, errors.New("invalid request body")
	}
	err = db.Create(&user).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = wallets_utils.CreateUserWallet(user)
	if err != nil {
		return nil, errors.New("unable to create user wallet")
	}

	serialized_user := user_serializers.SerializeUserSerializer(*user)

	return &serialized_user, nil
}

func GetUsers(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	userRepo := utils.NewGenericDB[models.Xuser](db)
	privilege := authenticated_user["privilege"]
	if privilege == "ADMIN" {
		users, err := userRepo.GetPagedAndFiltered(c.Locals("size").(int), c.Locals("page").(int), c.Locals("filters").(map[string]interface{}), nil, nil)
		if err != nil {
			return utils.BadRequestResponse(c, "Unable to get users")
		}
		serialized_users := user_serializers.SerializeUsers(users.Data)
		users.SerializedData = serialized_users
		users.Status = "Success"
		users.Message = "Successfully fetched users"
		users.Type = "OK"
		return utils.PaginatedSuccessResponse(c, users, "success")
	} else {
		users := []models.Xuser{}
		db.Order("created_at desc").First(&users, "id = ?", authenticated_user["id"])
		if users[0].Image != "" {
			imageUrl, err := userRepo.GetSignedUrl(c, "xuser")
			if err != nil {
				return utils.BadRequestResponse(c, err.Error())
			}
			serialized_users := user_serializers.SerializeUsers(users)
			serialized_users[0].Image = imageUrl
			return utils.SuccessResponse(c, serialized_users[0], "Successfully fetched users")

		}
		serialized_users := user_serializers.SerializeUsers(users)
		return utils.SuccessResponse(c, serialized_users[0], "Successfully fetched users")
	}

}


func UploadUserImage(c *fiber.Ctx) error {

	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	userRepo := utils.NewGenericDB[models.Xuser](db)
	privilege := authenticated_user["privilege"].(string)

	if strings.ToUpper(privilege) == "ADMIN" {
		return utils.BadRequestResponse(c, "this feature is not available for admins")
	}

	user, err := userRepo.UploadImage(c, "xuser")
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	serialized_user := user_serializers.SerializeUserSerializer(user.Data)
	user.SerializedData = serialized_user
	user.Status = "Success"
	user.Message = "Successfully uploaded user image"
	user.Type = "OK"
	return utils.SuccessResponse(c, serialized_user, "Successfully uploaded user image")


}