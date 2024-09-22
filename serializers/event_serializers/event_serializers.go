package event_serializers

import (
	"time"
	"txrnxp/models"
	"txrnxp/serializers/ticket_serializers"
	"txrnxp/serializers/user_serializers"

	"github.com/google/uuid"
)

type EventDetailSerializer struct {
	EventId     uuid.UUID                                            `json:"id" validate:"required"`
	Reference   string                                               `json:"reference" validate:"required"`
	Organiser   map[string]interface{}                               `json:"organiser" validate:"required"`
	IsBusiness  bool                                                 `json:"is_business" validate:"required"`
	Name        string                                               `json:"name" validate:"required"`
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

type ReadCreateEventSerializer struct {
	EventId     uuid.UUID `json:"id" validate:"required"`
	Reference   string    `json:"reference" validate:"required"`
	IsBusiness  bool      `json:"is_business" validate:"required"`
	Name        string    `json:"name" validate:"required"`
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
	Id          uuid.UUID `json:"id" validate:"required"`
	Reference   string    `json:"reference" validate:"required"`
	Count       int       `json:"count" validate:"required"`
	IsValidated bool      `json:"is_validated" validate:"required"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`
}

type ReadUserTicketSerializer struct {
	Id          uuid.UUID                                          `json:"id" validate:"required"`
	Reference   string                                             `json:"reference" validate:"required"`
	Event       ReadCreateEventSerializer                          `json:"event" validate:"required"`
	EventTicket ticket_serializers.EventTicketCustomuserSerializer `json:"event_ticket" validate:"required"`
	User        user_serializers.ExportUserSerializer              `json:"user" validate:"required"`
	Count       int                                                `json:"count" validate:"required"`
	IsValidated bool                                               `json:"is_validated" validate:"required"`
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


func SerializeReadUserTickets(user_tickets []models.UserTicket) ([]ReadUserTicketSerializer, error) {
	serialized_user_tickets := []ReadUserTicketSerializer{}
	for _, ticket := range user_tickets {
		event := ticket.EventTicket.Event
		event_ticket := ticket.EventTicket
		user := ticket.User

		serialized_event := ReadCreateEventSerializer{
			EventId: event.Id,
			Reference: event.Reference,
			IsBusiness: event.IsBusiness,
			Name: event.Name,
			EventType: event.EventType,
			Description: event.Description,
			Address: event.Address,
			Category: event.Category,
			Duration: event.Duration,
			StartTime: event.StartTime,
			EndTime: event.EndTime,
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
		}

		serialized_event_ticket := ticket_serializers.EventTicketCustomuserSerializer{
			Id: event_ticket.Id,
			Price: event_ticket.Price,
			Reference: event_ticket.Reference,
			IsPaid: event_ticket.IsPaid,
			IsInviteOnly: event_ticket.IsInviteOnly,
			TicketType: event_ticket.TicketType,
			Description: event_ticket.Description,
			Perks: event_ticket.Perks,
		}

		serialized_user := user_serializers.ExportUserSerializer{
			Id: user.Id,
			Email: user.Email,
			UserName: user.UserName,
			FirstName: user.FirstName,
			LastName: user.LastName,
			PhoneNumber: user.PhoneNumber,
		}

		serialized_user_ticket := ReadUserTicketSerializer{
			Id: ticket.Id,
			Reference: ticket.Reference,
			Event: serialized_event,
			EventTicket: serialized_event_ticket,
			User: serialized_user,
			Count: ticket.Count,
			IsValidated: ticket.IsValidated,
			CreatedAt: ticket.CreatedAt,
			UpdatedAt: ticket.UpdatedAt,
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
		IsValidated: user_ticket.IsValidated,
		CreatedAt:   user_ticket.CreatedAt,
		UpdatedAt:   user_ticket.UpdatedAt,
	}

	return serialized_user_ticket, nil
}