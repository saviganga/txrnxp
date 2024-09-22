package event_views

import (
	"fmt"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/event_serializers"
	"txrnxp/utils"
	"txrnxp/utils/event_utils"

	"github.com/gofiber/fiber/v2"
)

func GetEvents(c *fiber.Ctx) error {
	// find a way to get the event tickets
	db := initialisers.ConnectDb().Db
	events := []models.Event{}
	event_serializer := new(event_serializers.EventDetailSerializer)
	event_serializers := []event_serializers.EventDetailSerializer{}
	organiser_user := []models.Xuser{}
	organiser_business := []models.Business{}

	db.Order("created_at desc").Find(&events)

	for _, event := range events {

		organiser_details := make(map[string]interface{})

		// get the event organiser details
		organiser_id := event.OrganiserId
		is_business := event.IsBusiness

		if is_business {
			err := db.Model(&models.Business{}).First(&organiser_business, "id = ?", organiser_id).Error
			if err != nil {
				return utils.BadRequestResponse(c, fmt.Sprintf("oops! unable to fetch events - organiser: %s", event.Reference))
			}
			organiser_details["name"] = organiser_business[0].Name
			organiser_details["is_business"] = true
			organiser_details["id"] = organiser_id

		} else {
			err := db.Model(&models.Xuser{}).First(&organiser_user, "id = ?", organiser_id).Error
			if err != nil {
				return utils.BadRequestResponse(c, fmt.Sprintf("oops! unable to fetch events - organiser: %s", event.Reference))
			}
			organiser_details["name"] = organiser_user[0].UserName
			organiser_details["is_business"] = false
			organiser_details["id"] = organiser_id
		}

		// fill in the serializer
		event_serializer.EventId = event.Id
		event_serializer.Reference = event.Reference
		event_serializer.Organiser = organiser_details
		event_serializer.Name = event.Name
		event_serializer.EventType = event.EventType
		event_serializer.Description = event.Description
		event_serializer.Address = event.Address
		event_serializer.Category = event.Category
		event_serializer.Duration = event.Duration
		event_serializer.StartTime = event.StartTime
		event_serializer.EndTime = event.EndTime
		event_serializer.CreatedAt = event.CreatedAt
		event_serializer.UpdatedAt = event.UpdatedAt

		event_serializers = append(event_serializers, *event_serializer)

	}

	return utils.SuccessResponse(c, event_serializers, "Successfully fetched events")
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
