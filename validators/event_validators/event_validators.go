package event_validators

import (
	"fmt"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/event_serializers"
	"txrnxp/utils"

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
		business_reference := c.Get("Business")
		if business_reference == "" {
			return utils.BadRequestResponse(c, "oops! this is a business event, please pass in the business reference")
		}
		err := db.Model(&models.Business{}).Find(&organiser_business, "reference = ?", business_reference).Error
		if err != nil {
			return utils.BadRequestResponse(c, fmt.Sprintf("oops! unable to fetch events - organiser: %s", event.Reference))
		}
		// validate the organiser id
		if organiser_business[0].UserId.String() != authenticated_user["id"] {
			return utils.BadRequestResponse(c, "you are not the owner of this bsiness")
		}
		if organiser_business[0].Id.String() != event.OrganiserId {
			return utils.BadRequestResponse(c, "this feature is only available for event organisers")
		}

		c.Locals("organiser_id", organiser_business[0].Id.String())
	} else {
		if strings.ToUpper(c.Get("Entity")) == "BUSINESS" || strings.ToUpper(c.Get("Entity")) != "" {
			return utils.BadRequestResponse(c, "oops! this is not a business event")
		}
		err := db.Model(&models.Xuser{}).First(&organiser_user, "id = ?", event.OrganiserId).Error
		if err != nil {
			return utils.BadRequestResponse(c, fmt.Sprintf("oops! unable to fetch events - organiser: %s", event.Reference))
		}
		if organiser_user[0].Id.String() != authenticated_user["id"] {
			return utils.BadRequestResponse(c, "this feature is only available for event organisers")
		}
		c.Locals("organiser_id", organiser_user[0].Id.String())

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
