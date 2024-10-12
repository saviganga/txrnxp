package event_serializers

import (
	"errors"
	"time"
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/serializers/ticket_serializers"
	"txrnxp/serializers/user_serializers"
	"txrnxp/utilities"
	"txrnxp/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)


type UpdateEventSerializer struct {
	Name string `json:"name"`
}


type EventDetailSerializer struct {
	EventId     uuid.UUID                                            `json:"id" validate:"required"`
	Reference   string                                               `json:"reference" validate:"required"`
	Organiser   map[string]interface{}                               `json:"organiser" validate:"required"`
	IsBusiness  bool                                                 `json:"is_business" validate:"required"`
	Name        string                                               `json:"name" validate:"required"`
	Image       string                                               `json:"image" validate:"required"`
	EventType   string                                               `json:"type" validate:"required"`
	Description string                                               `json:"description" validate:"required"`
	Address     string                                               `json:"address" validate:"required"`
	Category    string                                               `json:"category" validate:"required"`
	Duration    string                                               `json:"duration" validate:"required"`
	Tickets     []ticket_serializers.EventTicketCustomuserSerializer `json:"tickets" validate:"required"`
	StartTime   time.Time                                            `json:"start_time" validate:"required"`
	EndTime     time.Time                                            `json:"end_time" validate:"required"`
	CreatedAt   time.Time                                            `json:"created_at" validate:"required"`
	UpdatedAt   time.Time                                            `json:"updated_at" validate:"required"`
}

type EventListSerializer struct {
	EventId     uuid.UUID              `json:"id" validate:"required"`
	Reference   string                 `json:"reference" validate:"required"`
	Organiser   map[string]interface{} `json:"organiser" validate:"required"`
	IsBusiness  bool                   `json:"is_business" validate:"required"`
	Name        string                 `json:"name" validate:"required"`
	Image       string                 `json:"image" validate:"required"`
	EventType   string                 `json:"type" validate:"required"`
	Description string                 `json:"description" validate:"required"`
	Address     string                 `json:"address" validate:"required"`
	Category    string                 `json:"category" validate:"required"`
	Duration    string                 `json:"duration" validate:"required"`
	StartTime   time.Time              `json:"start_time" validate:"required"`
	EndTime     time.Time              `json:"end_time" validate:"required"`
	CreatedAt   time.Time              `json:"created_at" validate:"required"`
	UpdatedAt   time.Time              `json:"updated_at" validate:"required"`
}

type ReadCreateEventSerializer struct {
	EventId     uuid.UUID `json:"id" validate:"required"`
	Reference   string    `json:"reference" validate:"required"`
	IsBusiness  bool      `json:"is_business" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	EventType   string    `json:"type" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Address     string    `json:"address" validate:"required"`
	Category    string    `json:"category" validate:"required"`
	Duration    string    `json:"duration" validate:"required"`
	StartTime   time.Time `json:"start_time" validate:"required"`
	EndTime     time.Time `json:"end_time" validate:"required"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`
}

type ReadEventTicketSerializer struct {
	Id             uuid.UUID                 `json:"id" validate:"required"`
	Event          ReadCreateEventSerializer `json:"event" validate:"required"`
	Price          string                    `json:"price" validate:"required"`
	Reference      string                    `json:"reference" validate:"required"`
	IsPaid         bool                      `json:"is_paid" validate:"required"`
	IsInviteOnly   bool                      `json:"is_invite_only"`
	TicketType     string                    `json:"ticket_type" validate:"required"`
	Description    string                    `json:"description" validate:"required"`
	Perks          map[string]interface{}    `json:"perks" validate:"required"`
	PurchaseLimit  int                       `json:"purchase_limit" validate:"required"`
	IsLimitedStock bool                      `json:"is_limited_stock" validate:"required"`
	StockNumber    int                       `json:"stock_number" validate:"required"`
	SoldTickets    int                       `json:"sold_tickets" validate:"required"`
	CreatedAt      time.Time                 `json:"created_at" validate:"required"`
	UpdatedAt      time.Time                 `json:"updated_at" validate:"required"`
}

type ReadCreateEventTicketSerializer struct {
	Id             uuid.UUID              `json:"id" validate:"required"`
	Price          string                 `json:"price" validate:"required"`
	Reference      string                 `json:"reference" validate:"required"`
	IsPaid         bool                   `json:"is_paid" validate:"required"`
	IsInviteOnly   bool                   `json:"is_invite_only"`
	TicketType     string                 `json:"ticket_type" validate:"required"`
	Description    string                 `json:"description" validate:"required"`
	Perks          map[string]interface{} `json:"perks" validate:"required"`
	PurchaseLimit  int                    `json:"purchase_limit" validate:"required"`
	IsLimitedStock bool                   `json:"is_limited_stock" validate:"required"`
	StockNumber    int                    `json:"stock_number" validate:"required"`
	SoldTickets    int                    `json:"sold_tickets" validate:"required"`
	CreatedAt      time.Time              `json:"created_at" validate:"required"`
	UpdatedAt      time.Time              `json:"updated_at" validate:"required"`
}

type ReadCreateUserTicketSerializer struct {
	Id          uuid.UUID              `json:"id" validate:"required"`
	Reference   string                 `json:"reference" validate:"required"`
	Count       int                    `json:"count" validate:"required"`
	Barcode     map[string]interface{} `json:"barcode" validate:"required"`
	IsValidated bool                   `json:"is_validated" validate:"required"`
	ValidCount  int                    `json:"valid_count" validate:"required"`
	CreatedAt   time.Time              `json:"created_at" validate:"required"`
	UpdatedAt   time.Time              `json:"updated_at" validate:"required"`
}

type ReadUserTicketSerializer struct {
	Id          uuid.UUID                                          `json:"id" validate:"required"`
	Reference   string                                             `json:"reference" validate:"required"`
	Event       ReadCreateEventSerializer                          `json:"event" validate:"required"`
	EventTicket ticket_serializers.EventTicketCustomuserSerializer `json:"event_ticket" validate:"required"`
	User        user_serializers.ExportUserSerializer              `json:"user" validate:"required"`
	Count       int                                                `json:"count" validate:"required"`
	Barcode     map[string]interface{}                             `json:"barcode" validate:"required"`
	IsValidated bool                                               `json:"is_validated" validate:"required"`
	ValidCount  int                                                `json:"valid_count" validate:"required"`
	CreatedAt   time.Time                                          `json:"created_at" validate:"required"`
	UpdatedAt   time.Time                                          `json:"updated_at" validate:"required"`
}

func PopulateEventOrganiserDetails(name string, is_business bool, id string) map[string]interface{} {
	return map[string]interface{}{
		"name":        name,
		"is_business": is_business,
		"id":          id,
	}
}

func SerializeReadEventsList(events []models.Event, c *fiber.Ctx) ([]EventListSerializer, error) {

	db := initialisers.ConnectDb().Db
	eventRepo := utils.NewGenericDB[models.Event](db)
	event_serializer := new(EventListSerializer)
	event_serializers := []EventListSerializer{}
	var imageUrl string

	for _, event := range events {

		// get the event organiser details
		organiser_id := event.OrganiserId

		organiser_details, err := utilities.GetEventOrganiser(organiser_id, event.Reference, event.IsBusiness)
		if err != nil {
			return nil, errors.New(err.Error())
		}

		if event.Image != "" {
			imageUrl, err = eventRepo.GetSignedUrl(c, "event", event.Id.String())
			if err != nil {
				return nil, errors.New(err.Error())
			}
		} else {
			imageUrl = ""
		}

		// fill in the serializer
		event_serializer.EventId = event.Id
		event_serializer.Reference = event.Reference
		event_serializer.Organiser = organiser_details
		event_serializer.Name = event.Name
		event_serializer.Image = imageUrl
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

	return event_serializers, nil
}

func SerializeReadUserTickets(user_tickets []models.UserTicket) ([]ReadUserTicketSerializer, error) {
	serialized_user_tickets := []ReadUserTicketSerializer{}
	for _, ticket := range user_tickets {
		event := ticket.EventTicket.Event
		event_ticket := ticket.EventTicket
		user := ticket.User

		serialized_event := ReadCreateEventSerializer{
			EventId:     event.Id,
			Reference:   event.Reference,
			IsBusiness:  event.IsBusiness,
			Name:        event.Name,
			EventType:   event.EventType,
			Description: event.Description,
			Address:     event.Address,
			Category:    event.Category,
			Duration:    event.Duration,
			StartTime:   event.StartTime,
			EndTime:     event.EndTime,
			CreatedAt:   event.CreatedAt,
			UpdatedAt:   event.UpdatedAt,
		}

		serialized_event_ticket := ticket_serializers.EventTicketCustomuserSerializer{
			Id:           event_ticket.Id,
			Price:        event_ticket.Price,
			Reference:    event_ticket.Reference,
			IsPaid:       event_ticket.IsPaid,
			IsInviteOnly: event_ticket.IsInviteOnly,
			TicketType:   event_ticket.TicketType,
			Description:  event_ticket.Description,
			Perks:        event_ticket.Perks,
		}

		serialized_user := user_serializers.ExportUserSerializer{
			Id:          user.Id,
			Email:       user.Email,
			UserName:    user.UserName,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			PhoneNumber: user.PhoneNumber,
		}

		serialized_user_ticket := ReadUserTicketSerializer{
			Id:          ticket.Id,
			Reference:   ticket.Reference,
			Event:       serialized_event,
			EventTicket: serialized_event_ticket,
			User:        serialized_user,
			Count:       ticket.Count,
			Barcode:     ticket.Barcode,
			IsValidated: ticket.IsValidated,
			ValidCount:  ticket.ValidCount,
			CreatedAt:   ticket.CreatedAt,
			UpdatedAt:   ticket.UpdatedAt,
		}

		serialized_user_tickets = append(serialized_user_tickets, serialized_user_ticket)
	}

	return serialized_user_tickets, nil
}

func SerializeCreateUserTickets(user_ticket models.UserTicket) (ReadCreateUserTicketSerializer, error) {
	serialized_user_ticket := ReadCreateUserTicketSerializer{
		Id:          user_ticket.Id,
		Reference:   user_ticket.Reference,
		Count:       user_ticket.Count,
		Barcode:     user_ticket.Barcode,
		IsValidated: user_ticket.IsValidated,
		ValidCount:  user_ticket.ValidCount,
		CreatedAt:   user_ticket.CreatedAt,
		UpdatedAt:   user_ticket.UpdatedAt,
	}

	return serialized_user_ticket, nil
}

func SerializeGetEventByReference(c *fiber.Ctx, event_list []models.Event, event_tickets []models.EventTicket, organiser_details map[string]interface{}) (*EventDetailSerializer, error) {

	db := initialisers.ConnectDb().Db
	event := new(EventDetailSerializer)
	eventTicket := new(ticket_serializers.EventTicketCustomuserSerializer)
	eventTickets := []ticket_serializers.EventTicketCustomuserSerializer{}
	var imageUrl string
	var err error
	eventRepo := utils.NewGenericDB[models.Event](db)

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

	if event_list[0].Image != "" {
		imageUrl, err = eventRepo.GetSignedUrl(c, "event", event_list[0].Id.String())
		if err != nil {
			return nil, errors.New(err.Error())
		}
	} else {
		imageUrl = ""
	}

	// fill in the serializer
	event.EventId = event_list[0].Id
	event.Reference = event_list[0].Reference
	event.Tickets = eventTickets
	event.Organiser = organiser_details
	event.Name = event_list[0].Name
	event.Image = imageUrl
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

func SerializeCreateEvent(event models.Event, c *fiber.Ctx) (*ReadCreateEventSerializer, error) {

	db := initialisers.ConnectDb().Db
	var imageUrl string
	var err error
	eventRepo := utils.NewGenericDB[models.Event](db)

	serialized_event := new(ReadCreateEventSerializer)

	if event.Image != "" {
		imageUrl, err = eventRepo.GetSignedUrl(c, "event", event.Id.String())
		if err != nil {
			return nil, errors.New(err.Error())
		}
	} else {
		imageUrl = ""
	}

	serialized_event.EventId = event.Id
	serialized_event.Reference = event.Reference
	serialized_event.Name = event.Name
	serialized_event.Image = imageUrl
	serialized_event.IsBusiness = event.IsBusiness
	serialized_event.EventType = event.EventType
	serialized_event.Description = event.Description
	serialized_event.Address = event.Address
	serialized_event.Category = event.Category
	serialized_event.Duration = event.Duration
	serialized_event.StartTime = event.StartTime
	serialized_event.EndTime = event.EndTime
	serialized_event.CreatedAt = event.CreatedAt
	serialized_event.UpdatedAt = event.UpdatedAt

	return serialized_event, nil

}

func SerializeCreateEventTicket(eventTicket models.EventTicket) ReadCreateEventTicketSerializer {

	serialized_event_ticket := new(ReadCreateEventTicketSerializer)

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

	return *serialized_event_ticket
}

func SerializeGetEventTickets(event_tickets []models.EventTicket) []ReadEventTicketSerializer {

	serialized_event_ticket := new(ReadEventTicketSerializer)
	serialized_event_tickets := []ReadEventTicketSerializer{}

	for _, ticket := range event_tickets {

		ticket_event := ticket.Event
		serialized_event := ReadCreateEventSerializer{
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

	return serialized_event_tickets
}
