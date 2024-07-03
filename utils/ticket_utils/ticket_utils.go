package ticket_utils

import (
	"errors"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func CreateEventTicket(c *fiber.Ctx) (*models.EventTicket, error) {

	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	eventTicket := new(models.EventTicket)
	event := []models.Event{}

	// create reference
	reference := utils.CreateEventReference()
	eventTicket.Reference = reference

	// pparse the request body
	err := c.BodyParser(eventTicket)
	if err != nil {
		return nil, errors.New("invalid request body")
	}

	// get the event
	result := db.First(&event, "id = ?", eventTicket.EventId)
	if result.Error != nil {
		return nil, errors.New("invalid event")
	}

	// check the entity and update the organiser id
	entity := c.Get("Entity")

	// improve this guy to be more specific on the user business
	if strings.ToUpper(entity) == "BUSINESS" {
		businesses := []models.Business{}
		db.First(&businesses, "user_id = ?", authenticated_user["id"])
		business_id := businesses[0].Id.String()
		if event[0].OrganiserId != business_id {
			return nil, errors.New("not event organiser")
		}
	} else {
		organiser_id := authenticated_user["id"].(string)
		if event[0].OrganiserId != organiser_id {
			return nil, errors.New("not event organiser")
		}
	}

	err = db.Create(&eventTicket).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return eventTicket, nil
}


func GetEventTickets(user_id string, entity string) ([]models.EventTicket, error) {
	db := initialisers.ConnectDb().Db
	event_tickets := []models.EventTicket{}

	if strings.ToUpper(entity) == "BUSINESS" {
		businesses := []models.Business{}
		db.First(&businesses, "user_id = ?", user_id)
		business_id := businesses[0].Id.String()
		// result := db.Model(&models.EventTicket{}).Joins("Event").Where("Event.organiser_id = ?", business_id).Find(&event_tickets)
		result := db.Joins("JOIN events ON event_tickets.event_id = events.id").Where("events.organiser_id = ?", business_id).Find(&event_tickets)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		organiser_id := user_id
		// result := db.Model(&models.EventTicket{}).Joins("Event").Where("Event.organiser_id = ?", organiser_id).Find(&event_tickets)
		result := db.Joins("JOIN events ON event_tickets.event_id = events.id").Where("events.organiser_id = ?", organiser_id).Find(&event_tickets)
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return event_tickets, nil
}


func GetUserTickets(user_id string, entity string) ([]models.UserTicket, error) {
	db := initialisers.ConnectDb().Db
	user_tickets := []models.UserTicket{}

	if strings.ToUpper(entity) == "BUSINESS" {
		businesses := []models.Business{}
		db.First(&businesses, "user_id = ?", user_id)
		business_id := businesses[0].Id.String()
		result := db.Model(&models.UserTicket{}).Joins("User").Joins("Event").Joins("EventTicket").Where("events.organiser_id = ?", business_id).Find(&user_tickets)
		// result := db.Joins("JOIN events ON user_ticket.event_id = events.id").Where("events.organiser_id = ?", business_id).Find(&user_tickets)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		organiser_id := user_id
		result := db.Model(&models.UserTicket{}).Joins("User").Find(&user_tickets, "user_tickets.user_id = ?", organiser_id)
		// result := db.Table("user_tickets").Joins("User").Find(&user_tickets, "user_ticket.user_id = ?", organiser_id)
		// result := db.Joins("User").Find(&user_tickets, "user_ticket.user_id == ?", organiser_id) //.Joins("Manager").Joins("Account").Find(&users, "users.id IN ?", []int{1,2,3,4,5})
		// result := db.Joins("user").First()
		if result.Error != nil {
			return nil, result.Error
		}
	}

	return user_tickets, nil
}