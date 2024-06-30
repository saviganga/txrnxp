package event_utils

import (
	"errors"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateEvent(c *fiber.Ctx) (*models.Event, error) {

	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	event := new(models.Event)

	// check the entity and update the organiser id
	entity := c.Get("Entity")

	// improve this guy to be more specific on the user business
	if strings.ToUpper(entity) == "BUSINESS" {
		businesses := []models.Business{}
		db.First(&businesses, "user_id = ?", authenticated_user["id"])
		organiser_id := businesses[0].Id.String()
		event.OrganiserId = organiser_id
	} else {
		organiser_id := authenticated_user["id"].(string)
		event.OrganiserId = organiser_id
	}

	reference := utils.CreateEventReference()
	event.Reference = reference

	err := c.BodyParser(event)
	if err != nil {
		return nil, errors.New("invalid request body")
	}
	err = db.Create(&event).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return event, nil
}
