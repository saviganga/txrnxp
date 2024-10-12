package event_validators

import (
	"fmt"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/utils"
	"txrnxp/serializers/event_serializers"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func ValidateEventOrganiser(c *fiber.Ctx) error {

	db := initialisers.ConnectDb().Db
	organiser_user := []models.Xuser{}
	organiser_business := []models.Business{}
	event_id := c.Params("id")
	event := models.Event{}
	authenticated_user := c.Locals("user").(jwt.MapClaims)

	// get the event model
	err := db.Model(&models.Event{}).First(&event, "id = ?", event_id).Error
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	if event.IsBusiness {
		err := db.Model(&models.Business{}).First(&organiser_business, "id = ?", event.OrganiserId).Error
		if err != nil {
			return utils.BadRequestResponse(c, fmt.Sprintf("oops! unable to fetch events - organiser: %s", event.Reference))
		}
		if organiser_business[0].UserId.String() != authenticated_user["id"] {
			return utils.BadRequestResponse(c, "this feature is only available for event organisers")
		}
	} else {
		err := db.Model(&models.Xuser{}).First(&organiser_user, "id = ?", event.OrganiserId).Error
		if err != nil {
			return utils.BadRequestResponse(c, fmt.Sprintf("oops! unable to fetch events - organiser: %s", event.Reference))
		}
		if organiser_user[0].Id.String() != authenticated_user["id"] {
			return utils.BadRequestResponse(c, "this feature is only available for event organisers")
		}

	}

	return c.Next()
}

func ValidateUpdateEventRequestBody(c *fiber.Ctx) error {

	// validate the request body
	body := new(event_serializers.UpdateEventSerializer)

	if err := c.BodyParser(&body); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body")
	}

	c.Locals("body", body)
	return c.Next()

}
