package ticket_serializers

type CreateUserTicketSerializer struct {
	EventTicketId string `json:"event_ticket_id" validate:"required"`
	UserId string `json:"user_id"`
	Count int `json:"count"`
}
