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

	serialized_business := business_serializers.SerializeCreateBusiness(*business)

	return &serialized_business, nil
}
