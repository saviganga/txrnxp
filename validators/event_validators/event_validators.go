package event_validators

import (
	"fmt"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func ValidateEventOrganiser(c *fiber.Ctx) error {

	db := initialisers.ConnectDb().Db
	organiser_user := []models.Xuser{}
	organiser_business := []models.Business{}
	// organiser_details := make(map[string]interface{})
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
		// organiser_details["name"] = organiser_business[0].Name
		// organiser_details["is_business"] = true
		// organiser_details["id"] = event.OrganiserId
		// organiser_details["business_user"] = organiser_business[0].UserId.String()
	} else {
		err := db.Model(&models.Xuser{}).First(&organiser_user, "id = ?", event.OrganiserId).Error
		if err != nil {
			return utils.BadRequestResponse(c, fmt.Sprintf("oops! unable to fetch events - organiser: %s", event.Reference))
		}
		if organiser_user[0].Id.String() != authenticated_user["id"] {
			return utils.BadRequestResponse(c, "this feature is only available for event organisers")
		}
		// organiser_details["name"] = organiser_user[0].UserName
		// organiser_details["is_business"] = false
		// organiser_details["id"] = event.OrganiserId
		// organiser_details["business_user"] = ""
	}

	return c.Next()
}
