package event_views

import (
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/utils"
	"txrnxp/utils/event_utils"

	"github.com/gofiber/fiber/v2"
)

func GetEvents(c *fiber.Ctx) error {
	// find a way to get the event tickets
	db := initialisers.ConnectDb().Db
	events := []models.Event{}
	db.Find(&events)
	return utils.SuccessResponse(c, events, "Successfully fetched events")
}

func CreateEvents(c *fiber.Ctx) error {

	event, err := event_utils.CreateEvent(c)

	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	return utils.CreatedResponse(c, event, "Successfully created event")
}



func GetEventByReference(c *fiber.Ctx) error {

	event, err := event_utils.GetEventByReference(c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	return utils.CreatedResponse(c, event, "Successfully fetched event")
}