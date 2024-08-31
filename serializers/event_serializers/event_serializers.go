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
