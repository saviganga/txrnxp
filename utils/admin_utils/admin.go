package admin_utils

import (
	"errors"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/admin_serializers"
	"txrnxp/serializers/user_serializers"
	"txrnxp/utils"
	"txrnxp/utils/db_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateAdminUser(c *fiber.Ctx) (*user_serializers.UserSerializer, error) {

	db := initialisers.ConnectDb().Db
	user := new(models.AdminUser)
	err := c.BodyParser(user)
	if err != nil {
		return nil, errors.New("invalid request body")
	}
	err = db.Create(&user).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	serialized_user := user_serializers.SerializeUser(*user)

	return &serialized_user, nil

}

func GetAdminUsers(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	users := []models.AdminUser{}
	db.First(&users, "id = ?", authenticated_user["id"])

	serialized_user := user_serializers.SerializeUser(users[0])

	return utils.SuccessResponse(c, serialized_user, "Successfully fetched users")

}

func CreateAdminCommission(c *fiber.Ctx) (*admin_serializers.ReadAdminCommissionConfigSerializer, error) {

	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	admin_commission_serializer := new(admin_serializers.CreateAdminCommissionConfigSerializer)
	privilege := authenticated_user["privilege"].(string)

	if strings.ToUpper(privilege) != "ADMIN" {
		return nil, errors.New("oops! you do not have permission to perform this action")
	}

	err := c.BodyParser(admin_commission_serializer)
	if err != nil {
		return nil, errors.New("invalid request body")
	}

	admin_commission := &models.AdminCommissionConfig{
		Type:       admin_commission_serializer.Type,
		Commission: admin_commission_serializer.Commission,
		Cap:        admin_commission_serializer.Cap,
	}

	err = db.Create(admin_commission).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}

	serialized_commission := admin_serializers.SerializeCreateAdminCommissionConfig(*admin_commission)

	return &serialized_commission, nil

}

func GetAdminCommission(c *fiber.Ctx) (*admin_serializers.ReadAdminCommissionConfigSerializer, error) {

	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	admin_commission := models.AdminCommissionConfig{}
	privilege := authenticated_user["privilege"].(string)

	if strings.ToUpper(privilege) != "ADMIN" {
		return nil, errors.New("oops! you do not have permission to perform this action")
	}

	db.First(&admin_commission)

	serialized_commission := admin_serializers.SerializeCreateAdminCommissionConfig(admin_commission)

	return &serialized_commission, nil

}

func UpdateAdminCommission(c *fiber.Ctx) (*admin_serializers.ReadAdminCommissionConfigSerializer, error) {

	commission_id := c.Params("id")
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	admin_commission_serializer := new(admin_serializers.CreateAdminCommissionConfigSerializer)
	admin_commission := models.AdminCommissionConfig{}
	privilege := authenticated_user["privilege"].(string)

	if strings.ToUpper(privilege) != "ADMIN" {
		return nil, errors.New("oops! you do not have permission to perform this action")
	}

	err := c.BodyParser(admin_commission_serializer)
	if err != nil {
		return nil, errors.New("invalid request body")
	}

	err = db.First(&admin_commission, "id = ?", commission_id).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}

	admin_commission.Type = admin_commission_serializer.Type
	admin_commission.Commission = admin_commission_serializer.Commission
	admin_commission.Cap = admin_commission_serializer.Cap

	is_updated_admin_commission, updated_admin_commission := db_utils.UpdateAdminCommissionConfig(&admin_commission)
	if !is_updated_admin_commission {
		return nil, errors.New(updated_admin_commission)
	}

	serialized_commission := admin_serializers.SerializeCreateAdminCommissionConfig(admin_commission)

	return &serialized_commission, nil

}
