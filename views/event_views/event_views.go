package event_views

import (
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/event_serializers"
	"txrnxp/utils"
	"txrnxp/utils/event_utils"

	"github.com/gofiber/fiber/v2"
)

func GetEvents(c *fiber.Ctx) error {
	db := initialisers.ConnectDb().Db
	events := []models.Event{}

	db.Order("created_at desc").Find(&events)

	serialized_events, err := event_serializers.SerializeReadEventsList(events)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, serialized_events, "Successfully fetched events")
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
	return utils.SuccessResponse(c, event, "Successfully fetched event")
}
