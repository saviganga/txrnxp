package event_views

import (
	"strings"
	"time"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/event_serializers"
	"txrnxp/utils"
	"txrnxp/utils/event_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func GetEvents(c *fiber.Ctx) error {
	db := initialisers.ConnectDb().Db
	eventRepo := utils.NewGenericDB[models.Event](db)

	filters := c.Locals("filters").(map[string]interface{})
	filters["start_time"] = time.Now()

	events, err := eventRepo.GetPagedAndFiltered(c.Locals("size").(int), c.Locals("page").(int), filters, nil, nil)
	if err != nil {
		return utils.BadRequestResponse(c, "Unable to get events")
	}

	serialized_events, err := event_serializers.SerializeReadEventsList(events.Data, c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	events.SerializedData = serialized_events
	events.Status = "Success"
	events.Message = "Successfully fetched events"
	events.Type = "OK"

	return utils.PaginatedSuccessResponse(c, events, "success")
}

func GetMyEvents(c *fiber.Ctx) error {

	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	privilege := authenticated_user["privilege"]
	eventRepo := utils.NewGenericDB[models.Event](db)
	filters := c.Locals("filters").(map[string]interface{})

	if strings.ToUpper(privilege.(string))  == "ADMIN" {
		return utils.BadRequestResponse(c, "oops! this feature is not available for admins")
	}

	// check if it is a user or business
	entity := c.Get("Entity")

	if strings.ToUpper(entity) == "BUSINESS" {
		// find the business
		business_reference := c.Get("Business")
		if business_reference == "" {
			return utils.BadRequestResponse(c, "oops! please pass in the business reference")
		}
		business := models.Business{}
		err := db.First(&business, "reference = ?", business_reference).Error
		if err != nil {
			return utils.BadRequestResponse(c, "oops! this business does not exist")
		}

		// validate the organiser id
		if business.UserId.String() != authenticated_user["id"].(string) {
			return utils.BadRequestResponse(c, "this user does not own this business")
		}

		filters["organiser_id"] = business.Id.String()

	} else {

		filters["organiser_id"] = authenticated_user["id"].(string)

	}

	events, err := eventRepo.GetPagedAndFiltered(c.Locals("size").(int), c.Locals("page").(int), filters, nil, nil)
	if err != nil {
		return utils.BadRequestResponse(c, "Unable to get events")
	}

	serialized_events, err := event_serializers.SerializeReadEventsList(events.Data, c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	events.SerializedData = serialized_events
	events.Status = "Success"
	events.Message = "Successfully fetched events"
	events.Type = "OK"

	return utils.PaginatedSuccessResponse(c, events, "success")

}

func GetEventById(c *fiber.Ctx) error {
	return event_utils.GetEventById(c)
}

func UpdateEvent(c *fiber.Ctx) error {
	return event_utils.UpdateEvent(c)
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

func UploadEventImage(c *fiber.Ctx) error {
	return event_utils.UploadEventImage(c)
}
