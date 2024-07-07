package ticket_utils

import (
	"errors"
	"strconv"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/ticket_serializers"
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

	// parse the request body
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
		result := db.Joins("JOIN events ON event_tickets.event_id = events.id").Where("events.organiser_id = ?", business_id).Find(&event_tickets)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {
		organiser_id := user_id
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

func CreateUserTicket(c *fiber.Ctx) (*models.UserTicket, error) {

	// initialise niggas
	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	userTicket := new(models.UserTicket)
	eventTicket := []models.EventTicket{}
	userwallets := []models.UserWallet{}

	entity := c.Get("Entity")
	if strings.ToUpper(entity) != "" {
		return nil, errors.New("oops! this feature is not available to businesses")
	}

	user_request := new(ticket_serializers.CreateUserTicketSerializer)
	err := c.BodyParser(user_request)
	if err != nil {
		return nil, errors.New("invalid request body")
	}

	// get the eventTicket
	result := db.First(&eventTicket, "id = ?", user_request.EventTicketId)
	if result.Error != nil {
		return nil, errors.New("invalid eventTicket")
	}
	// fill in the eventId field
	userTicket.EventId = eventTicket[0].EventId

	// validate user
	if user_request.UserId == "" {
		userTicket_userId, err := utils.ConvertStringToUUID(authenticated_user["id"].(string))
		if err != nil {
			return nil, errors.New("invalid parsed id")
		}
		userTicket.UserId = userTicket_userId
	} else {
		userTicket_userId, err := utils.ConvertStringToUUID(user_request.UserId)
		if err != nil {
			return nil, errors.New("invalid parsed id")
		}
		userTicket.UserId = userTicket_userId
	}

	// create reference
	reference := utils.CreateEventReference()
	userTicket.Reference = reference


	// get user wallet
	db.Model(&models.UserWallet{}).Joins("User").First(&userwallets, "user_wallets.user_id = ?", userTicket.UserId)

	// check event ticket conditions
	// price
	eventTicket_float_price, err := strconv.ParseFloat(eventTicket[0].Price, 64)
	if err != nil {
		return nil, errors.New("error converting event ticket price")
	}
	userWallet_available_balance, err := strconv.ParseFloat(userwallets[0].AvailableBalance, 64)
	if err != nil {
		return nil, errors.New("error converting wallet available balance")
	}
	if userWallet_available_balance < eventTicket_float_price {
		return nil, errors.New("oops! insufficient wallet funds")
	}

	// limited stock/stock number

	// purchase limit

	// count

	err = db.Create(&userTicket).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return userTicket, nil
}
