package event_utils

import (
	"errors"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/event_serializers"
	"txrnxp/serializers/ticket_serializers"
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

func GetEventByReference(c *fiber.Ctx) (*event_serializers.EventDetailSerializer, error) {

	db := initialisers.ConnectDb().Db
	// authenticated_user := c.Locals("user").(jwt.MapClaims)
	reference := c.Params("reference")
	event_list := []models.Event{}
	event := new(event_serializers.EventDetailSerializer)
	event_tickets := []models.EventTicket{}
	organiser_user := []models.Xuser{}
	organiser_business := []models.Business{}
	is_business := false
	organiser_details := make(map[string]interface{})
	eventTicket := new(ticket_serializers.EventTicketCustomuserSerializer)
	eventTickets := []ticket_serializers.EventTicketCustomuserSerializer{}

	// get the event model
	err := db.Model(&models.Event{}).First(&event_list, "reference = ?", reference).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}

	// get the event organiser details
	organiser_id := event_list[0].OrganiserId
	err = db.Model(&models.Xuser{}).First(&organiser_user, "id = ?", organiser_id).Error
	if err == nil {
		is_business = false
	} else {
		is_business = true
		db.Model(&models.Business{}).First(&organiser_business, "id = ?", organiser_id)
	}

	if is_business {
		organiser_details["name"] = organiser_business[0].Name
	} else {
		organiser_details["name"] = organiser_user[0].UserName
	}

	// get the event ticket details
	err = db.Model(&models.EventTicket{}).Find(&event_tickets, "event_id = ?", event_list[0].Id).Error
	if err != nil {
		return nil, errors.New(err.Error())
	}

	for _, eventticket := range event_tickets {

		eventTicket.Id = eventticket.Id
		eventTicket.IsPaid = eventticket.IsPaid
		eventTicket.IsInviteOnly = eventticket.IsInviteOnly
		eventTicket.Reference = eventticket.Reference
		eventTicket.TicketType = eventticket.TicketType
		eventTicket.Description = eventticket.Description
		eventTicket.Perks = eventticket.Perks
		eventTicket.Price = eventticket.Price

		eventTickets = append(eventTickets, *eventTicket)

	}

	// fill in the serializer
	event.EventId = event_list[0].Id
	event.Reference = reference
	event.Tickets = eventTickets
	event.Organiser = organiser_details
	event.Name = event_list[0].Name
	event.EventType = event_list[0].EventType
	event.Description = event_list[0].Description
	event.Address = event_list[0].Address
	event.Category = event_list[0].Category
	event.Duration = event_list[0].Duration
	event.StartTime = event_list[0].StartTime
	event.EndTime = event_list[0].EndTime
	event.CreatedAt = event_list[0].CreatedAt
	event.UpdatedAt = event_list[0].UpdatedAt

	return event, nil

}
