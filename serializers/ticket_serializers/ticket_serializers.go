package ticket_serializers

type CreateUserTicketSerializer struct {
	EventTicketId string `json:"event_ticket_id" validate:"required"`
	UserId string `json:"user_id"`
	Count int `json:"count"`
}

type TransferUserTicketSerializer struct {
	UserTicketReference string `json:"ticket_reference" validate:"required"`
	ReceiverEmail string `json:"receiver_email" validate:"required"`
	Count int `json:"count" validate:"required"`
}
