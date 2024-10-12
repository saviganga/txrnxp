package business_validators

import (
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/business_serializers"
	"txrnxp/utils"

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


func ValidateUpdateBusinessRequestBody(c *fiber.Ctx) error {

	// validate the request body
	body := new(business_serializers.UpdateBusinessSerializer)

	if err := c.BodyParser(&body); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	c.Locals("body", body)
	return c.Next()


}
