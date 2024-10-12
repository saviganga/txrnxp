package ticket_serializers

import (
	"github.com/google/uuid"
)

type CreateUserTicketSerializer struct {
	EventTicketId string `json:"event_ticket_id" validate:"required"`
	UserId        string `json:"user_id"`
	Count         int    `json:"count"`
}

type TransferUserTicketSerializer struct {
	UserTicketReference string `json:"ticket_reference" validate:"required"`
	ReceiverEmail       string `json:"receiver_email" validate:"required"`
	Count               int    `json:"count" validate:"required"`
}

type EventTicketCustomuserSerializer struct {
	Id           uuid.UUID              `json:"id" validate:"required"`
	Price        string                 `json:"price" validate:"required"`
	Reference    string                 `json:"reference" validate:"required"`
	IsPaid       bool                   `json:"is_paid" validate:"required"`
	IsInviteOnly bool                   `json:"is_invite_only"`
	TicketType   string                 `json:"ticket_type" validate:"required"`
	Description  string                 `json:"description" validate:"required"`
	Perks        map[string]interface{} `json:"perks" validate:"required"`
}

type ValidateUserTicket struct {
	Count         int    `json:"count"`
}
