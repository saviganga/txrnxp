package event_utils

import (
	"errors"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/event_serializers"
	"txrnxp/utilities"
	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateEvent(c *fiber.Ctx) (*event_serializers.ReadCreateEventSerializer, error) {

	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	event := new(models.Event)
	privilege := authenticated_user["privilege"]

	if privilege == "ADMIN" {
		return nil, errors.New("oops! this feature is not available for admins")
	}

	// check the entity and update the organiser id
	entity := c.Get("Entity")

	// improve this guy to be more specific on the user business
	if strings.ToUpper(entity) == "BUSINESS" {
		business_reference := c.Get("Business")
		businesses := []models.Business{}
		err := db.Find(&businesses, "user_id = ? AND reference = ?", authenticated_user["id"].(string), business_reference).Error
		if err != nil {
			return nil, errors.New("oops! this user is not a business")
		}
		organiser_id := businesses[0].Id.String()
		event.OrganiserId = organiser_id
		event.IsBusiness = true
	} else {
		organiser_id := authenticated_user["id"].(string)
		event.OrganiserId = organiser_id
		event.IsBusiness = false
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

	serialized_event, err := event_serializers.SerializeCreateEvent(*event, c)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return serialized_event, nil
}

func GetEventByReference(c *fiber.Ctx) (*event_serializers.EventDetailSerializer, error) {

	db := initialisers.ConnectDb().Db
	reference := c.Params("reference")
	event_list := []models.Event{}
	// event := new(event_serializers.EventDetailSerializer)
	event_tickets := []models.EventTicket{}

	// get the event model
	err := db.Model(&models.Event{}).First(&event_list, "reference = ?", reference).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}

	organiser_details, err := utilities.GetEventOrganiser(event_list[0].OrganiserId, event_list[0].Reference, event_list[0].IsBusiness)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// get the event ticket details
	err = db.Model(&models.EventTicket{}).Find(&event_tickets, "event_id = ?", event_list[0].Id).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}

	serialized_event, err := event_serializers.SerializeGetEventByReference(c, event_list, event_tickets, organiser_details)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return serialized_event, nil

}

func GetEventById(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	event := models.Event{}
	privilege := authenticated_user["privilege"]
	err := db.First(&event, "id = ?", c.Params("id")).Error
	if err != nil {
		return utils.BadRequestResponse(c, "Unable to get event")
	}
	if privilege != "ADMIN" && authenticated_user["id"].(string) != event.OrganiserId {
		return utils.BadRequestResponse(c, "You do not have permission to view this resource")
	}

	serialized_event, err := event_serializers.SerializeCreateEvent(event, c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, serialized_event, "success")

}


func UpdateEvent(c *fiber.Ctx) error {

	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	eventRepo := utils.NewGenericDB[models.Event](db)
	privilege := authenticated_user["privilege"].(string)
	event_id := c.Params("id")

	if strings.ToUpper(privilege) == "ADMIN" {
		return utils.BadRequestResponse(c, "this feature is not available for admins")
	}

	event, err := eventRepo.UpdateEntity(c, "event", event_id)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	serialized_event, err := event_serializers.SerializeCreateEvent(event.Data, c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, serialized_event, "success")


}

func UploadEventImage(c *fiber.Ctx) error {

	authenticated_user := c.Locals("user").(jwt.MapClaims)
	db := initialisers.ConnectDb().Db
	eventRepo := utils.NewGenericDB[models.Event](db)
	privilege := authenticated_user["privilege"].(string)
	event_id := c.Params("id")

	if strings.ToUpper(privilege) == "ADMIN" {
		return utils.BadRequestResponse(c, "this feature is not available for admins")
	}

	event, err := eventRepo.UploadImage(c, "event", event_id)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	serialized_event, err := event_serializers.SerializeCreateEvent(event.Data, c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	event.SerializedData = serialized_event
	event.Status = "Success"
	event.Message = "Successfully uploaded event image"
	event.Type = "OK"
	return utils.SuccessResponse(c, serialized_event, "Successfully uploaded event image")

}
