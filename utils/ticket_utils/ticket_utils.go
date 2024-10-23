package ticket_utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/event_serializers"
	"txrnxp/serializers/ticket_serializers"
	"txrnxp/utils"
	"txrnxp/utils/admin_utils"
	"txrnxp/utils/business_utils"
	"txrnxp/utils/db_utils"
	"txrnxp/utils/wallets_utils"

	"image/png"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

func CreateEventTicket(c *fiber.Ctx) (*event_serializers.ReadCreateEventTicketSerializer, error) {

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

	entity := c.Get("Entity")

	// improve this guy to be more specific on the user business
	if strings.ToUpper(entity) == "BUSINESS" {
		businesses := []models.Business{}
		business_member := models.BusinessMember{}
		business_reference := c.Get("Business")
		if business_reference == "" {
			return nil, errors.New("oops! this is a business event, please pass in the business reference")
		}
		err = db.Find(&businesses, "reference = ?", business_reference).Error
		if err != nil {
			return nil, errors.New("oops! you are not the event organiser")
		}

		business_id := businesses[0].Id.String()

		err = db.Model(&models.BusinessMember{}).First(&business_member, "user_id = ? AND business_id = ?", authenticated_user["id"], business_id).Error
		if err != nil || business_member.Id == uuid.Nil {
			return nil, errors.New("oop! this user is not a business member")
		}

		// // validate the organiser id
		// if businesses[0].UserId.String() != authenticated_user["id"] {
		// 	return nil, errors.New("oops! you are not the business owner")
		// }

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

	serialized_event_ticket := event_serializers.SerializeCreateEventTicket(*eventTicket)

	return &serialized_event_ticket, nil

}

func GetEventTickets(user_id string, entity string, c *fiber.Ctx) error {
	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	privilege := authenticated_user["privilege"].(string)
	repo := utils.NewGenericDB[models.EventTicket](db)

	limit := c.Locals("size").(int)
	page := c.Locals("page").(int)
	joins := []string{"LEFT JOIN events as event ON event_tickets.event_id = event.id"}
	preloads := []string{"Event"}

	if strings.ToUpper(privilege) == "ADMIN" {

		filters := c.Locals("filters").(map[string]interface{})
		event_tickets, err := repo.GetPagedAndFiltered(limit, page, filters, preloads, joins)
		if err != nil {
			return utils.BadRequestResponse(c, "Unable to get event tickets")
		}
		serialized_event_tickets := event_serializers.SerializeGetEventTickets(event_tickets.Data)
		event_tickets.SerializedData = serialized_event_tickets
		event_tickets.Status = "Success"
		event_tickets.Message = "Successfully fetched event tickets"
		event_tickets.Type = "OK"
		return utils.PaginatedSuccessResponse(c, event_tickets, "success")

	} else {

		if strings.ToUpper(entity) == "BUSINESS" {
			filters := make(map[string]interface{})
			business_reference := c.Get("Business")
			if business_reference == "" {
				return utils.BadRequestResponse(c, "oops! this is a business event, please pass in the business reference")
			}

			businesses := []models.Business{}

			err := db.First(&businesses, "reference = ?", business_reference).Error
			if err != nil {
				return utils.BadRequestResponse(c, "Oops! This user is not a business")
			}

			filters["event__id"] = c.Params("id")

			event_tickets, err := repo.GetPagedAndFiltered(limit, page, filters, preloads, joins)
			if err != nil {
				return utils.BadRequestResponse(c, "Unable to get businesses")
			}
			serialized_event_tickets := event_serializers.SerializeGetEventTickets(event_tickets.Data)
			event_tickets.SerializedData = serialized_event_tickets
			event_tickets.Status = "Success"
			event_tickets.Message = "Successfully fetched event_tickets"
			event_tickets.Type = "OK"
			return utils.PaginatedSuccessResponse(c, event_tickets, "success")

		} else {

			filters := make(map[string]interface{})
			organiser_id := user_id
			business_utils.RemoveBusinessKeys(filters)
			filters["event__organiser_id"] = organiser_id

			event_tickets, err := repo.GetPagedAndFiltered(limit, page, filters, preloads, joins)
			if err != nil {
				return utils.BadRequestResponse(c, "Unable to get businesses")
			}
			serialized_event_tickets := event_serializers.SerializeGetEventTickets(event_tickets.Data)
			event_tickets.SerializedData = serialized_event_tickets
			event_tickets.Status = "Success"
			event_tickets.Message = "Successfully fetched event_tickets"
			event_tickets.Type = "OK"
			return utils.PaginatedSuccessResponse(c, event_tickets, "success")

		}
	}

}

func GetEventTicketById(c *fiber.Ctx) (*event_serializers.ReadEventTicketSerializer, error) {

	db := initialisers.ConnectDb().Db
	event_ticket := models.EventTicket{}
	event_ticket_id := c.Params("id")
	entity := c.Get("Entity")
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	user_id := authenticated_user["id"].(string)

	if strings.ToUpper(entity) == "BUSINESS" {
		business_reference := c.Get("Business")
		if business_reference == "" {
			return nil, errors.New("oops! this is a business event, please pass in the business reference")
		}
		businesses := []models.Business{}
		business_member := models.BusinessMember{}

		err := db.Find(&businesses, "reference = ?", business_reference).Error
		if err != nil {
			return nil, errors.New("oops! this user is not a business")
		}
		business_id := businesses[0].Id.String()

		err = db.Model(&models.BusinessMember{}).First(&business_member, "user_id = ? AND business_id = ?", authenticated_user["id"], business_id).Error
		if err != nil || business_member.Id == uuid.Nil {
			return nil, errors.New("oop! this user is not a business member")
		}

		// // validate the organiser id
		// if businesses[0].UserId.String() != authenticated_user["id"] {
		// 	return nil, errors.New("oops! this user is not the business owner")
		// }

		result := db.Model(&models.EventTicket{}).
			Preload("Event").
			// Joins("Event").
			Where("id = ?", event_ticket_id).
			Order("created_at desc").
			First(&event_ticket)

		if result.Error != nil {
			return nil, result.Error
		}

		if event_ticket.Event.OrganiserId != business_id {
			return nil, errors.New("oops! this feature is only available to event organisers")
		}

	} else {
		if strings.ToUpper(authenticated_user["privilege"].(string)) != "ADMIN" {
			result := db.Model(&models.EventTicket{}).
				Preload("Event").
				Where("id = ?", event_ticket_id).
				Order("created_at desc").
				First(&event_ticket)

			if result.Error != nil {
				return nil, result.Error
			}

			if event_ticket.Event.OrganiserId != user_id {
				return nil, errors.New("oops! you do not have permission to view this resource")
			}
		} else {

			result := db.Model(&models.EventTicket{}).
				Preload("Event").
				Where("id = ?", event_ticket_id).
				Order("created_at desc").
				First(&event_ticket)

			if result.Error != nil {
				return nil, result.Error
			}

			if event_ticket.Event.IsBusiness {
				return nil, errors.New("oops! you do not have permission to view this resource")
			}

		}

	}

	serialized_event_ticket := event_serializers.SerializeGetEventTicket(event_ticket)
	return &serialized_event_ticket, nil

}

func GetUserTickets(user_id string, entity string, c *fiber.Ctx) error {
	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	privilege := authenticated_user["privilege"].(string)

	repo := utils.NewGenericDB[models.UserTicket](db)

	limit := c.Locals("size").(int)
	page := c.Locals("page").(int)
	joins := []string{"LEFT JOIN event_tickets as event_ticket ON user_tickets.event_ticket_id = event_ticket.id", "LEFT JOIN xusers as u ON user_tickets.user_id = u.id"}
	preloads := []string{"EventTicket.Event", "User"}
	filters := c.Locals("filters").(map[string]interface{})

	if strings.ToUpper(privilege) == "ADMIN" {
		user_tickets, err := repo.GetPagedAndFiltered(limit, page, filters, preloads, joins)
		if err != nil {
			return utils.BadRequestResponse(c, "Unable to get event tickets")
		}

		serialized_user_tickets, err := event_serializers.SerializeReadUserTickets(user_tickets.Data)
		if err != nil {
			return utils.BadRequestResponse(c, err.Error())
		}
		user_tickets.SerializedData = serialized_user_tickets
		user_tickets.Status = "Success"
		user_tickets.Message = "Successfully fetched event tickets"
		user_tickets.Type = "OK"
		return utils.PaginatedSuccessResponse(c, user_tickets, "success")
	}

	if strings.ToUpper(entity) == "BUSINESS" {
		return utils.BadRequestResponse(c, "This feature is not available for businesses")
	} else {

		filters["u__id"] = user_id

		user_tickets, err := repo.GetPagedAndFiltered(limit, page, filters, preloads, joins)
		if err != nil {
			return utils.BadRequestResponse(c, "Unable to get event tickets")
		}

		serialized_user_tickets, err := event_serializers.SerializeReadUserTickets(user_tickets.Data)
		if err != nil {
			return utils.BadRequestResponse(c, err.Error())
		}
		user_tickets.SerializedData = serialized_user_tickets
		user_tickets.Status = "Success"
		user_tickets.Message = "Successfully fetched event tickets"
		user_tickets.Type = "OK"
		return utils.PaginatedSuccessResponse(c, user_tickets, "success")

	}

}


func GetEventAttendees(user_id string, entity string, c *fiber.Ctx) error {
	db := initialisers.ConnectDb().Db
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	privilege := authenticated_user["privilege"].(string)

	repo := utils.NewGenericDB[models.UserTicket](db)

	limit := c.Locals("size").(int)
	page := c.Locals("page").(int)
	joins := []string{"LEFT JOIN event_tickets as event_ticket ON user_tickets.event_ticket_id = event_ticket.id", "LEFT JOIN xusers as u ON user_tickets.user_id = u.id", "LEFT JOIN events as event ON user_tickets.event_id = event.id"}
	preloads := []string{"EventTicket.Event", "User", "Event"}
	filters := c.Locals("filters").(map[string]interface{})

	if strings.ToUpper(privilege) == "ADMIN" {
		user_tickets, err := repo.GetPagedAndFiltered(limit, page, filters, preloads, joins)
		if err != nil {
			return utils.BadRequestResponse(c, "Unable to get event tickets")
		}

		serialized_user_tickets, err := event_serializers.SerializeReadUserTickets(user_tickets.Data)
		if err != nil {
			return utils.BadRequestResponse(c, err.Error())
		}
		user_tickets.SerializedData = serialized_user_tickets
		user_tickets.Status = "Success"
		user_tickets.Message = "Successfully fetched event tickets"
		user_tickets.Type = "OK"
		return utils.PaginatedSuccessResponse(c, user_tickets, "success")
	}

	if strings.ToUpper(entity) == "BUSINESS" {
		// business := models.Business{}
		// entity := c.Get("Entity")
		// business_reference := c.Get("Business")
		// if business_reference == "" || strings.ToUpper(entity) != "BUSINESS" {
		// 	return utils.BadRequestResponse(c, "this is a business event, please pass in the business header")
		// }
		// err := db.Model(&models.Business{}).First(&business, "reference = ? AND user_id = ?", business_reference, authenticated_user["id"]).Error
		// if err != nil {
		// 	return utils.BadRequestResponse(c, "you do not have permission to perform this action")
		// }
		user_tickets, err := repo.GetPagedAndFiltered(limit, page, filters, preloads, joins)
		if err != nil {
			return utils.BadRequestResponse(c, "Unable to get event tickets")
		}

		serialized_user_tickets, err := event_serializers.SerializeReadUserTickets(user_tickets.Data)
		if err != nil {
			return utils.BadRequestResponse(c, err.Error())
		}
		user_tickets.SerializedData = serialized_user_tickets
		user_tickets.Status = "Success"
		user_tickets.Message = "Successfully fetched event tickets"
		user_tickets.Type = "OK"
		return utils.PaginatedSuccessResponse(c, user_tickets, "success")
	} else {

		// filters["u__id"] = user_id

		user_tickets, err := repo.GetPagedAndFiltered(limit, page, filters, preloads, joins)
		if err != nil {
			return utils.BadRequestResponse(c, "Unable to get event tickets")
		}

		serialized_user_tickets, err := event_serializers.SerializeReadUserTickets(user_tickets.Data)
		if err != nil {
			return utils.BadRequestResponse(c, err.Error())
		}
		user_tickets.SerializedData = serialized_user_tickets
		user_tickets.Status = "Success"
		user_tickets.Message = "Successfully fetched event tickets"
		user_tickets.Type = "OK"
		return utils.PaginatedSuccessResponse(c, user_tickets, "success")

	}

}

func GetUserTicketByReference(c *fiber.Ctx) (*event_serializers.ReadUserTicketSerializer, error) {

	db := initialisers.ConnectDb().Db
	user_tickets := []models.UserTicket{}
	reference := c.Params("reference")
	entity := c.Get("Entity")
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	user_id := authenticated_user["id"].(string)

	if strings.ToUpper(entity) == "BUSINESS" {
		return nil, errors.New("oops! this feature is not available for businesses")
	} else {
		if strings.ToUpper(authenticated_user["privilege"].(string)) != "ADMIN" {
			result := db.Model(&models.UserTicket{}).
				Preload("EventTicket.Event").
				Joins("User").
				Where("user_tickets.reference = ?", reference).
				Order("created_at desc").
				First(&user_tickets)

			if result.Error != nil {
				return nil, result.Error
			}
			user_id_uuid, err := utils.ConvertStringToUUID(user_id)
			if err != nil {
				return nil, err
			}
			if user_tickets[0].EventTicket.Event.OrganiserId != user_id && user_tickets[0].UserId != user_id_uuid {
				return nil, errors.New("oops! you do not have permission to view this resource")
			}
		} else {

			result := db.Model(&models.UserTicket{}).
				Preload("EventTicket.Event").
				Joins("User").
				Where("user_tickets.reference = ?", reference).
				Order("created_at desc").
				First(&user_tickets)

			if result.Error != nil {
				return nil, result.Error
			}

		}

	}

	serialized_user_tickets, err := event_serializers.SerializeReadUserTickets(user_tickets)
	if err != nil {
		return nil, err
	}

	return &serialized_user_tickets[0], nil

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
	result := db.Model(&models.EventTicket{}).
		Joins("LEFT JOIN events ON events.id = event_tickets.event_id").
		Preload("Event").
		Order("event_tickets.created_at desc").
		First(&eventTicket, "event_tickets.id = ?", user_request.EventTicketId)
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
		return nil, errors.New(err.Error())
	}

	// admin commission
	is_paid, commission_amount, err := admin_utils.PayAdminCommission("event_tickets", eventTicket[0].Price, eventTicket[0].Reference, user_request.Count)
	if !is_paid {
		return nil, errors.New(err.Error())
	}

	amount_str := eventTicket[0].Price
	amount_float, err := utils.ConvertStringToFloat(amount_str)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	commission_amount_float, err := utils.ConvertStringToFloat(commission_amount)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	organiser_amount := (amount_float * float64(user_request.Count)) - commission_amount_float
	organiser_amount_str := fmt.Sprintf("%.2f", organiser_amount)

	// credit event organiser wallet
	entry_description := fmt.Sprintf("ticket purchase commission - %s", eventTicket[0].Reference)
	db.Find(&events, "id = ?", eventTicket[0].EventId)
	if events[0].IsBusiness {
		err = db.First(&businesses, "id = ?", events[0].OrganiserId).Error
		if err != nil {
			return nil, errors.New("oops! nable to find event organiser")
		}
		is_credited, credited_wallet := wallets_utils.CreditUserWallet(businesses[0].UserId, organiser_amount_str, entry_description)
		if !is_credited {
			return nil, errors.New(credited_wallet)
		}
	} else {
		err = db.First(&users, "id = ?", events[0].OrganiserId).Error
		if err != nil {
			return nil, errors.New("oops! nable to find event organiser")
		}
		is_credited, credited_wallet := wallets_utils.CreditUserWallet(users[0].Id, organiser_amount_str, entry_description)
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
		entry_description := fmt.Sprintf("ticket purchase - %s", eventTicket.Reference)
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

		// generate ticket barcode
		barcode_url, barcode_image, err := GenerateUserTicketBarcode(userTicket.Reference)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		barcode_img := base64.StdEncoding.EncodeToString(barcode_image)
		barcode := make(map[string]interface{})
		barcode["url"] = barcode_url
		barcode["code"] = barcode_img
		userTicket.Barcode = barcode

		err = db.Create(&userTicket).Error
		if err != nil {
			return nil, errors.New(err.Error())
		}

		// increase userticket count
		is_updated_user_ticket, updated_user_ticket := db_utils.UpdateUserTicket(userTicket)
		if !is_updated_user_ticket {
			return nil, errors.New(updated_user_ticket)
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

		// generate ticket barcode
		barcode_url, barcode_image, err := GenerateUserTicketBarcode(userTicket.Reference)
		if err != nil {
			return false, err.Error()
		}
		barcode_img := base64.StdEncoding.EncodeToString(barcode_image)
		barcode := make(map[string]interface{})
		barcode["url"] = barcode_url
		barcode["code"] = barcode_img
		userTicket.Barcode = barcode

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

func GenerateUserTicketBarcode(reference string) (string, []byte, error) {

	url := fmt.Sprintf("https://ec08-197-211-59-80.ngrok-free.app/api/v1/tickets/users/%s", reference)

	// Create a new QR code object
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		return "", nil, errors.New(err.Error())
	}

	// Convert image to bytes
	imgBuf := new(bytes.Buffer)
	err = png.Encode(imgBuf, qr.Image(256))
	if err != nil {
		return "", nil, errors.New(err.Error())
	}

	imageData := imgBuf.Bytes()

	return url, imageData, nil

}

func ValidateUserTicket(c *fiber.Ctx) (bool, string) {
	db := initialisers.ConnectDb().Db
	ticket_reference := c.Params("reference")
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	user_request := new(ticket_serializers.ValidateUserTicket)
	userTicket := models.UserTicket{}
	privilege := authenticated_user["privilege"].(string)
	business := models.Business{}
	user := models.Xuser{}

	// validate request body
	err := c.BodyParser(user_request)
	if err != nil {
		return false, "Invalid request body"
	}

	// get the ticket
	result := db.Model(&models.UserTicket{}).
		Preload("EventTicket.Event").
		Joins("User").
		Where("user_tickets.reference = ?", ticket_reference).
		Order("created_at desc").
		First(&userTicket)

	if result.Error != nil {
		return false, "Unable to get user ticket"
	}

	// validate event organiser
	organiser_id_uuid, err := utils.ConvertStringToUUID(userTicket.EventTicket.Event.OrganiserId)
	if err != nil {
		return false, "Oops! unable to validate request user"
	}

	if userTicket.EventTicket.Event.IsBusiness {
		entity := c.Get("Entity")
		business_reference := c.Get("Business")
		business_member := models.BusinessMember{}
		if business_reference == "" || strings.ToUpper(entity) != "BUSINESS" {
			return false, "oops! this is a business event, please pass in the business reference"
		}
		err := db.Model(&models.Business{}).Find(&business, "reference = ?", business_reference).Error
		if err != nil {
			return false, fmt.Sprintf("oops! unable to fetch business - reference: %s", business_reference)
		}
		err = db.Model(&models.BusinessMember{}).First(&business_member, "user_id = ? AND business_id = ?", authenticated_user["id"], business.Id.String()).Error
		if err != nil || business_member.Id == uuid.Nil {
			return false, fmt.Sprintf("oops! unable to fetch business member - business: %s", business_reference)
		}
	} else {
		db.Model(&models.Xuser{}).First(&user, "id = ?", authenticated_user["id"])
		organiser_id := user.Id
		if organiser_id != organiser_id_uuid && strings.ToUpper(privilege) != "ADMIN" {
			return false, "Oops! you do not have permission to perform this action"
		}
	}

	// check if the ticket is already valid
	if userTicket.IsValidated {
		return false, "Oops! this ticket has already been validated"
	}

	// get the ticket count and valid_count, increase ticket valid_count
	ticket_count := userTicket.Count
	ticket_valid_count := userTicket.ValidCount
	ticket_valid_count += user_request.Count
	if ticket_valid_count > ticket_count {
		return false, "Oops! invalid number of tickets to validate"
	}

	// if ticket count == valid_count: is_validated == true, increase valid count, save, return
	if ticket_valid_count == ticket_count {
		userTicket.IsValidated = true
		userTicket.ValidCount = ticket_valid_count
	} else { // if ticket count != valid_count: is_validated == false, increase valid count, save, return
		userTicket.ValidCount = ticket_valid_count
	}

	is_updated_user_ticket, updated_user_ticket := db_utils.UpdateUserTicket(&userTicket)
	if !is_updated_user_ticket {
		return false, updated_user_ticket
	}

	return true, "Successfully validated user ticket"

}
