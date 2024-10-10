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
	eventRepo := utils.NewGenericDB[models.Event](db)

	events, err := eventRepo.GetPagedAndFiltered(c.Locals("size").(int), c.Locals("page").(int), c.Locals("filters").(map[string]interface{}), nil, nil)
	if err != nil {
		return utils.BadRequestResponse(c, "Unable to get businesses")
	}

	serialized_events, err := event_serializers.SerializeReadEventsList(events.Data)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	events.SerializedData = serialized_events
	events.Status = "Success"
	events.Message = "Successfully fetched events"
	events.Type = "OK"

	return utils.PaginatedSuccessResponse(c, events, "success")
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
