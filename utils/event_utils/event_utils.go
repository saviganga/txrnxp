package event_utils

import (
	"errors"
	"fmt"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/event_serializers"
	"txrnxp/serializers/ticket_serializers"
	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateEvent(c *fiber.Ctx) (*event_serializers.ReadCreateEventSerializer, error) {

	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	event := new(models.Event)
	event_serializer := new(event_serializers.ReadCreateEventSerializer)
	privilege := authenticated_user["privilege"]

	if privilege == "ADMIN" {
		return nil, errors.New("oops! this feature is not available for admins")
	}

	// check the entity and update the organiser id
	entity := c.Get("Entity")

	// improve this guy to be more specific on the user business
	if strings.ToUpper(entity) == "BUSINESS" {
		businesses := []models.Business{}
		db.First(&businesses, "user_id = ?", authenticated_user["id"])
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

	// fill in the serializer
	event_serializer.EventId = event.Id
	event_serializer.Reference = event.Reference
	event_serializer.Name = event.Name
	event_serializer.IsBusiness = event.IsBusiness
	event_serializer.EventType = event.EventType
	event_serializer.Description = event.Description
	event_serializer.Address = event.Address
	event_serializer.Category = event.Category
	event_serializer.Duration = event.Duration
	event_serializer.StartTime = event.StartTime
	event_serializer.EndTime = event.EndTime
	event_serializer.CreatedAt = event.CreatedAt
	event_serializer.UpdatedAt = event.UpdatedAt

	return event_serializer, nil
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
	is_business := event_list[0].IsBusiness

	if is_business {
		err := db.Model(&models.Business{}).First(&organiser_business, "id = ?", organiser_id).Error
		if err != nil {
			return nil, fmt.Errorf("oops! unable to fetch events - organiser: %s", event.Reference)
		}
		organiser_details["name"] = organiser_business[0].Name
		organiser_details["is_business"] = true
		organiser_details["id"] = organiser_id
	} else {
		err := db.Model(&models.Xuser{}).First(&organiser_user, "id = ?", organiser_id).Error
		if err != nil {
			return nil, fmt.Errorf("oops! unable to fetch events - organiser: %s", event.Reference)
		}
		organiser_details["name"] = organiser_user[0].UserName
		organiser_details["is_business"] = false
		organiser_details["id"] = organiser_id
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
