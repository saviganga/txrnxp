package business_validators

import (
	"fmt"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/business_serializers"
	"txrnxp/utils"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func ValidateBusinessOwner(c *fiber.Ctx) error {

	db := initialisers.ConnectDb().Db
	business_id := c.Params("id")
	business := models.Business{}
	authenticated_user := c.Locals("user").(jwt.MapClaims)

	// get the business model
	err := db.Model(&models.Business{}).First(&business, "id = ?", business_id).Error
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	if authenticated_user["id"].(string) != business.UserId.String() {
		return utils.BadRequestResponse(c, "this feature is only available for business owners")
	}

	return c.Next()
}

func ValidateBusinessMember(c *fiber.Ctx) error {

	db := initialisers.ConnectDb().Db
	business := models.Business{}
	business_member := models.BusinessMember{}
	business_reference := c.Get("Business")
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	if business_reference == "" {
		return utils.BadRequestResponse(c, "oops! this is a business event, please pass in the business reference")
	}
	err := db.Model(&models.Business{}).Find(&business, "reference = ?", business_reference).Error
	if err != nil {
		return utils.BadRequestResponse(c, fmt.Sprintf("oops! unable to fetch business - reference: %s", business_reference))
	}
	err = db.Model(&models.BusinessMember{}).First(&business_member, "user_id = ? AND business_id = ?", authenticated_user["id"], business.Id.String()).Error
	if err != nil || business_member.Id == uuid.Nil {
		return utils.BadRequestResponse(c, fmt.Sprintf("oops! unable to fetch business member - business: %s", business_reference))
	}

	return c.Next()
}

func ValidateUpdateBusinessRequestBody(c *fiber.Ctx) error {

	// validate the request body
	body := new(business_serializers.UpdateBusinessSerializer)

	if err := c.BodyParser(&body); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	c.Locals("body", body)
	return c.Next()

}
