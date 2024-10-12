package business_utils

import (
	"errors"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/business_serializers"

	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateBusiness(c *fiber.Ctx) (*business_serializers.ReadCreateBusinessSerializer, error) {

	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	business := new(models.Business)
	privilege := authenticated_user["privilege"]

	if privilege == "ADMIN" {
		return nil, errors.New("oops! this feature is not available for admins")
	}

	parsedUUID, err := utils.ConvertStringToUUID(authenticated_user["id"].(string))

	if err != nil {
		return nil, errors.New("invalid parsed id")
	}

	business.UserId = parsedUUID
	err = c.BodyParser(business)
	if err != nil {
		return nil, errors.New("invalid request body")
	}
	err = db.Create(&business).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	serialized_business, err := business_serializers.SerializeCreateBusiness(*business, c)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return serialized_business, nil
}


func GetBusinessById(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	business := models.Business{}
	privilege := authenticated_user["privilege"]
	err := db.First(&business, "id = ?", c.Params("id")).Error
		if err != nil {
			return utils.BadRequestResponse(c, "Unable to get user")
		}
	if privilege != "ADMIN" && authenticated_user["id"].(string) != business.UserId.String() {
		return utils.BadRequestResponse(c, "You do not have permission to view this resource")
	}

	serialized_user, err := business_serializers.SerializeCreateBusiness(business, c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, serialized_user, "success")

}

func UploadBusinessImage(c *fiber.Ctx) error {

	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	businessRepo := utils.NewGenericDB[models.Business](db)
	privilege := authenticated_user["privilege"].(string)
	business_id := c.Params("id")

	if strings.ToUpper(privilege) == "ADMIN" {
		return utils.BadRequestResponse(c, "this feature is not available for admins")
	}

	business, err := businessRepo.UploadImage(c, "business", business_id)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	serialized_business, err := business_serializers.SerializeCreateBusiness(business.Data, c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	business.SerializedData = serialized_business
	business.Status = "Success"
	business.Message = "Successfully uploaded business image"
	business.Type = "OK"
	return utils.SuccessResponse(c, serialized_business, "Successfully uploaded business image")

}
