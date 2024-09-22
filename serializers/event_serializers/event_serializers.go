package event_serializers

import (
	"time"
	"txrnxp/serializers/ticket_serializers"

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
	EventId     uuid.UUID                                            `json:"id" validate:"required"`
	Reference   string                                               `json:"reference" validate:"required"`
	IsBusiness  bool                                                 `json:"is_business" validate:"required"`
	Name        string                                               `json:"name" validate:"required"`
	EventType   string                                               `json:"type" validate:"required"`
	Description string                                               `json:"description" validate:"required"`
	Address     string                                               `json:"address" validate:"required"`
	Category    string                                               `json:"category" validate:"required"`
	Duration    string                                               `json:"duration" validate:"required"`
	StartTime   time.Time                                            `json:"start_time" validate:"required"`
	EndTime     time.Time                                            `json:"end_time" validate:"required"`
	CreatedAt   time.Time                                            `json:"created_at" validate:"required"`
	UpdatedAt   time.Time                                            `json:"updated_at" validate:"required"`
}

type ReadEventTicketSerializer struct {
	Id             uuid.UUID              `json:"id" validate:"required"`
	Event          EventDetailSerializer  `json:"event" validate:"required"`
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


func PopulateEventOrganiserDetails(name string, is_business bool, id string) map[string]interface{} {
	return map[string]interface{}{
        "name":        name,
        "is_business": is_business,
        "id":          id,
    }
}