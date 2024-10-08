package ticket_utils

import (
	"errors"
	"strings"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/event_serializers"
	"txrnxp/serializers/ticket_serializers"
	"txrnxp/utils"
	"txrnxp/utils/db_utils"
	"txrnxp/utils/wallets_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func CreateEventTicket(c *fiber.Ctx) (*event_serializers.ReadCreateEventTicketSerializer, error) {

	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	eventTicket := new(models.EventTicket)
	event := []models.Event{}

	serialized_event_ticket := new(event_serializers.ReadCreateEventTicketSerializer)

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

	entity := c.Get("Entity")

	// improve this guy to be more specific on the user business
	if strings.ToUpper(entity) == "BUSINESS" {
		businesses := []models.Business{}
		err = db.First(&businesses, "user_id = ?", authenticated_user["id"]).Error
		if err != nil {
			return nil, errors.New("oops! you are not the event organiser")
		}
		business_id := businesses[0].Id.String()
		if event[0].OrganiserId != business_id {
			return nil, errors.New("oops! you are not the event organiser")
		}
	} else {
		organiser_id := authenticated_user["id"].(string)
		if event[0].OrganiserId != organiser_id {
			return nil, errors.New("oops! you are not the event organiser")
		}
	}

	err = db.Create(&eventTicket).Error

	if err != nil {
		return nil, errors.New(err.Error())
	}

	serialized_event_ticket.Id = eventTicket.Id
	serialized_event_ticket.Price = eventTicket.Price
	serialized_event_ticket.Reference = eventTicket.Reference
	serialized_event_ticket.IsPaid = eventTicket.IsPaid
	serialized_event_ticket.IsInviteOnly = eventTicket.IsInviteOnly
	serialized_event_ticket.TicketType = eventTicket.TicketType
	serialized_event_ticket.Description = eventTicket.Description
	serialized_event_ticket.Perks = eventTicket.Perks
	serialized_event_ticket.PurchaseLimit = eventTicket.PurchaseLimit
	serialized_event_ticket.IsLimitedStock = eventTicket.IsLimitedStock
	serialized_event_ticket.StockNumber = eventTicket.StockNumber
	serialized_event_ticket.SoldTickets = eventTicket.SoldTickets
	serialized_event_ticket.CreatedAt = eventTicket.CreatedAt
	serialized_event_ticket.UpdatedAt = eventTicket.UpdatedAt

	return serialized_event_ticket, nil
}

func GetEventTickets(user_id string, entity string, c *fiber.Ctx) ([]event_serializers.ReadEventTicketSerializer, error) {
	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	event_tickets := []models.EventTicket{}
	serialized_event_ticket := new(event_serializers.ReadEventTicketSerializer)
	serialized_event_tickets := []event_serializers.ReadEventTicketSerializer{}
	privilege := authenticated_user["privilege"].(string)

	if strings.ToUpper(privilege) == "ADMIN" {
		result := db.Model(&models.EventTicket{}).Joins("Event").Order("created_at desc").Find(&event_tickets)
		if result.Error != nil {
			return nil, result.Error
		}
	} else {

		if strings.ToUpper(entity) == "BUSINESS" {
			businesses := []models.Business{}
			err := db.First(&businesses, "user_id = ?", user_id).Error
			if err != nil {
				return nil, errors.New("oops! this user is not a business")
			}
			business_id := businesses[0].Id.String()
			result := db.Model(&models.EventTicket{}).
				Joins("LEFT JOIN events ON events.id = event_tickets.event_id").
				Where("events.organiser_id = ?", business_id).
				Preload("Event").
				Order("event_tickets.created_at desc").
				Find(&event_tickets)
			if result.Error != nil {
				return nil, result.Error
			}
		} else {
			organiser_id := user_id
			result := db.Model(&models.EventTicket{}).
				Joins("LEFT JOIN events ON events.id = event_tickets.event_id").
				Where("events.organiser_id = ?", organiser_id).
				Preload("Event").
				Order("event_tickets.created_at desc").
				Find(&event_tickets)
			if result.Error != nil {
				return nil, result.Error
			}

		}
	}

	for _, ticket := range event_tickets {

		ticket_event := ticket.Event
		serialized_event := event_serializers.ReadCreateEventSerializer{
			EventId:     ticket_event.Id,
			Reference:   ticket_event.Reference,
			IsBusiness:  ticket_event.IsBusiness,
			Name:        ticket_event.Name,
			EventType:   ticket_event.EventType,
			Description: ticket_event.Description,
			Address:     ticket_event.Address,
			Category:    ticket_event.Category,
			Duration:    ticket_event.Duration,
			StartTime:   ticket_event.StartTime,
			EndTime:     ticket_event.EndTime,
			CreatedAt:   ticket_event.CreatedAt,
			UpdatedAt:   ticket_event.UpdatedAt,
		}

		serialized_event_ticket.Id = ticket.Id
		serialized_event_ticket.Event = serialized_event
		serialized_event_ticket.Price = ticket.Price
		serialized_event_ticket.Reference = ticket.Reference
		serialized_event_ticket.IsPaid = ticket.IsPaid
		serialized_event_ticket.IsInviteOnly = ticket.IsInviteOnly
		serialized_event_ticket.TicketType = ticket.TicketType
		serialized_event_ticket.Description = ticket.Description
		serialized_event_ticket.Perks = ticket.Perks
		serialized_event_ticket.PurchaseLimit = ticket.PurchaseLimit
		serialized_event_ticket.IsLimitedStock = ticket.IsLimitedStock
		serialized_event_ticket.StockNumber = ticket.StockNumber
		serialized_event_ticket.SoldTickets = ticket.SoldTickets
		serialized_event_ticket.CreatedAt = ticket.CreatedAt
		serialized_event_ticket.UpdatedAt = ticket.UpdatedAt

		serialized_event_tickets = append(serialized_event_tickets, *serialized_event_ticket)

	}

	return serialized_event_tickets, nil
}

func GetUserTickets(user_id string, entity string) ([]event_serializers.ReadUserTicketSerializer, error) {
	db := initialisers.ConnectDb().Db
	user_tickets := []models.UserTicket{}

	if strings.ToUpper(entity) == "BUSINESS" {
		return nil, errors.New("oops! this feature is not available for businesses")
	} else {
		result := db.Model(&models.UserTicket{}).
			Preload("EventTicket.Event").
			Joins("User").
			Where("user_tickets.user_id = ?", user_id).
			Order("created_at desc").
			Find(&user_tickets)

		if result.Error != nil {
			return nil, result.Error
		}
	}

	serialized_user_tickets, err := event_serializers.SerializeReadUserTickets(user_tickets)
	if err != nil {
		return nil, err
	}

	return serialized_user_tickets, nil
}

func CreateUserTicket(c *fiber.Ctx) (*event_serializers.ReadCreateUserTicketSerializer, error) {

	// initialise niggas
	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	userTicket := new(models.UserTicket)
	eventTicket := []models.EventTicket{}
	events := []models.Event{}
	privilege := authenticated_user["privilege"].(string)
	entity := c.Get("Entity")
	user_request := new(ticket_serializers.CreateUserTicketSerializer)
	businesses := []models.Business{}
	users := []models.Xuser{}

	// validate request user
	if strings.ToUpper(privilege) == "ADMIN" {
		return nil, errors.New("oops! this feature is not available for admins")
	}
	if strings.ToUpper(entity) != "" {
		return nil, errors.New("oops! this feature is not available to businesses")
	}

	// validate request body
	err := c.BodyParser(user_request)
	if err != nil {
		return nil, errors.New("invalid request body")
	}

	if user_request.Count <= 0 {
		user_request.Count = 1
	}

	// get the eventTicket
	result := db.Find(&eventTicket, "id = ?", user_request.EventTicketId)
	if result.Error != nil {
		return nil, errors.New("invalid eventTicket")
	}
	// fill in the eventId field
	userTicket.EventId = eventTicket[0].EventId
	userTicket.EventTicketId = eventTicket[0].Id

	// validate user
	userTicket_userId, err := utils.ConvertStringToUUID(authenticated_user["id"].(string))
	if err != nil {
		return nil, errors.New("invalid parsed id")
	}
	userTicket.UserId = userTicket_userId

	// check event ticket conditions and create user ticket
	userTicket, err = ValidateCreateUserTicketConditions(userTicket, &eventTicket[0], user_request.Count)

	if err != nil {
		return nil, err
	}

	// credit event organiser wallet
	entry_description := "ticket purchase commission"
	db.Find(&events, "id = ?", eventTicket[0].EventId)
	if events[0].IsBusiness {
		err = db.First(&businesses, "id = ?", events[0].OrganiserId).Error
		if err != nil {
			return nil, errors.New("oops! nable to find event organiser")
		}
		is_credited, credited_wallet := wallets_utils.CreditUserWallet(businesses[0].UserId, eventTicket[0].Price, entry_description)
		if !is_credited {
			return nil, errors.New(credited_wallet)
		}
	} else {
		err = db.First(&users, "id = ?", events[0].OrganiserId).Error
		if err != nil {
			return nil, errors.New("oops! nable to find event organiser")
		}
		is_credited, credited_wallet := wallets_utils.CreditUserWallet(users[0].Id, eventTicket[0].Price, entry_description)
		if !is_credited {
			return nil, errors.New(credited_wallet)
		}
	}

	serialized_user_ticket, err := event_serializers.SerializeCreateUserTickets(*userTicket)
	if err != nil {
		return nil, err
	}

	return &serialized_user_ticket, nil

}

func ValidateCreateUserTicketConditions(userTicket *models.UserTicket, eventTicket *models.EventTicket, ticket_count int) (*models.UserTicket, error) {

	db := initialisers.ConnectDb().Db
	db.Where("user_id = ? AND event_ticket_id = ?", userTicket.UserId, userTicket.EventTicketId).Find(&userTicket)

	if userTicket.Id != uuid.Nil {

		// event ticket limited stock/stock number
		if eventTicket.IsLimitedStock {

			// compare stock number and sold tickets
			if eventTicket.SoldTickets >= eventTicket.StockNumber {
				return nil, errors.New("oops! ticket out of stock")
			}
		}

		// purchase limit
		userTicket.Count += ticket_count
		if eventTicket.PurchaseLimit > 0 {
			if userTicket.Count > eventTicket.PurchaseLimit {
				return nil, errors.New("oops! purchase limit reached")
			}
		}

		// calculate ticket price based on count
		// convert amount to float
		amount_float, err := utils.ConvertStringToFloat(eventTicket.Price)
		if err != nil || amount_float == 0.0 {
			return nil, errors.New("oops! an error occured")
		}
		amount := float64(ticket_count) * amount_float

		// debit user wallet
		entry_description := "ticket purchase"
		is_debited, debited_wallet := wallets_utils.DebitUserWallet(userTicket.UserId, amount, entry_description)
		if !is_debited {
			return nil, errors.New(debited_wallet)
		}

		// update sold tickets - event ticket
		eventTicket.SoldTickets += ticket_count
		is_updated_event_ticket, updated_event_ticket := db_utils.UpdateEventTicket(eventTicket)
		if !is_updated_event_ticket {
			return nil, errors.New(updated_event_ticket)
		}

		// increase userticket count
		is_updated_user_ticket, updated_user_ticket := db_utils.UpdateUserTicket(userTicket)
		if !is_updated_user_ticket {
			return nil, errors.New(updated_user_ticket)
		}

		return userTicket, nil

	} else {

		// create reference
		reference := utils.CreateEventReference()
		userTicket.Reference = reference
		userTicket.EventTicketId = eventTicket.Id

		// event ticket limited stock/stock number
		if eventTicket.IsLimitedStock {

			// compare stock number and sold tickets
			if eventTicket.SoldTickets >= eventTicket.StockNumber {
				return nil, errors.New("oops! ticket out of stock")
			}
		}

		// purchase limit
		userTicket.Count = ticket_count
		if eventTicket.PurchaseLimit > 0 {
			if userTicket.Count > eventTicket.PurchaseLimit {
				return nil, errors.New("oops! purchase limit reached")
			}
		}

		amount_float, err := utils.ConvertStringToFloat(eventTicket.Price)
		if err != nil || amount_float == 0.0 {
			return nil, errors.New("oops! an error occured")
		}

		amount := amount_float * float64(ticket_count)

		// debit user wallet
		entry_description := "ticket purchase"
		is_debited, debited_wallet := wallets_utils.DebitUserWallet(userTicket.UserId, amount, entry_description)
		if !is_debited {
			return nil, errors.New(debited_wallet)
		}

		// update sold tickets - event ticket
		eventTicket.SoldTickets += ticket_count
		is_updated_event_ticket, updated_event_ticket := db_utils.UpdateEventTicket(eventTicket)
		if !is_updated_event_ticket {
			return nil, errors.New(updated_event_ticket)
		}

		err = db.Create(&userTicket).Error

		if err != nil {
			return nil, errors.New(err.Error())
		}

		return userTicket, nil
	}

}

func TransferUserTicket(c *fiber.Ctx) (bool, string) {
	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	transfer_request := new(ticket_serializers.TransferUserTicketSerializer)
	privilege := authenticated_user["privilege"].(string)
	users := []models.Xuser{}
	user_tickets := []models.UserTicket{}
	userTicket := new(models.UserTicket)
	userTickets := []models.UserTicket{}

	if strings.ToUpper(privilege) != "USER" {
		return false, "Oops! this feature is only available for users"
	}

	err := c.BodyParser(transfer_request)
	if err != nil {
		return false, err.Error()
	}

	if transfer_request.ReceiverEmail == authenticated_user["email"].(string) {
		return false, "Oops! you cannot transfer to yourself."
	}

	// get the sender id
	sender_id_string := authenticated_user["id"].(string)
	sender_id_uuid, err := utils.ConvertStringToUUID(sender_id_string)
	if err != nil {
		return false, err.Error()
	}

	// get the receiver id
	db.First(&users, "email = ?", transfer_request.ReceiverEmail)
	receiver_id_uuid := users[0].Id

	// find the userticket
	db.First(&user_tickets, "reference = ?", transfer_request.UserTicketReference)

	// validate user ticket values
	user_ticket := user_tickets[0]

	// user ticket owner
	if user_ticket.UserId != sender_id_uuid {
		return false, "Oops! You do not have permission to transfer this ticket"
	}

	// user ticket count
	if user_ticket.Count < transfer_request.Count {
		return false, "Oops! Insufficient tickets"
	}

	// check if the receiver already has usertickets for the event ticket
	db.Where("user_id = ? AND event_ticket_id = ?", receiver_id_uuid, userTicket.EventTicketId).Find(&userTickets)

	if userTickets[0].Id != uuid.Nil {
		userTickets[0].EventId = user_ticket.EventId
		userTickets[0].EventTicketId = user_ticket.EventTicketId
		userTickets[0].UserId = receiver_id_uuid
		userTickets[0].Count += transfer_request.Count
		is_updated_user_ticket, updated_user_ticket := db_utils.UpdateUserTicket(&userTickets[0])
		if !is_updated_user_ticket {
			return false, updated_user_ticket
		}
	} else {
		// create user ticket reference
		reference := utils.CreateEventReference()

		// fill new userTicket values
		userTicket.EventId = user_ticket.EventId
		userTicket.EventTicketId = user_ticket.EventTicketId
		userTicket.UserId = receiver_id_uuid
		userTicket.Reference = reference
		userTicket.Count += transfer_request.Count
		err = db.Create(&userTicket).Error
		if err != nil {
			return false, err.Error()
		}
	}

	// decrease userticket count
	user_ticket.Count -= transfer_request.Count
	is_updated_user_ticket, updated_user_ticket := db_utils.UpdateUserTicket(&user_ticket)
	if !is_updated_user_ticket {
		return false, updated_user_ticket
	}

	return true, "Transfer successful"

}
