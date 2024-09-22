package business_utils

import (
	"errors"
	"fmt"
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
	serialized_business := new(business_serializers.ReadCreateBusinessSerializer)
	privilege := authenticated_user["privilege"]

	if privilege == "ADMIN" {
		return nil, errors.New("oops! this feature is not available for admins")
	}

	parsedUUID, err := utils.ConvertStringToUUID(authenticated_user["id"].(string))

	if err != nil {
		return nil, errors.New("invalid parsed id")
	} else {
		fmt.Println("Successfully parsed UUID:", parsedUUID)
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

	serialized_business.Id = business.Id
	serialized_business.Reference = business.Reference
	serialized_business.Name = business.Name
	serialized_business.Country = business.Country
	serialized_business.CreatedAt = business.CreatedAt
	serialized_business.UpdatedAt = business.UpdatedAt

	return serialized_business, nil
}
