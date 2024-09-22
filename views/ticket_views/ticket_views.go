package ticket_views

import (
	"txrnxp/utils"
	"txrnxp/utils/ticket_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the United States of Ganga tickets mehn!")
}

func CreateEventTicket(c *fiber.Ctx) error {

	event, err := ticket_utils.CreateEventTicket(c)

	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}
	return utils.CreatedResponse(c, event, "Successfully created event")
}

func GetEventTickets(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	entity := c.Get("Entity")
	user_id := authenticated_user["id"].(string)

	event_tickets, err := ticket_utils.GetEventTickets(user_id, entity, c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, event_tickets, "Successfully fetched event tickets")

}

func GetUserTickets(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	entity := c.Get("Entity")
	if entity == "" {
		return utils.BadRequestResponse(c, "Please pass in your header")
	}
	user_id := authenticated_user["id"].(string)

	user_tickets, err := ticket_utils.GetUserTickets(user_id, entity)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, user_tickets, "Successfully fetched user tickets")

}

func CreateUserTicket(c *fiber.Ctx) error {
	user_ticket, err := ticket_utils.CreateUserTicket(c)
	if err != nil {
		return utils.BadRequestResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, user_ticket, "Success")

}

// transfer tickets between users
func TransferUserTicket(c *fiber.Ctx) error {

	is_transferred, ticket_transfer := ticket_utils.TransferUserTicket(c)
	if !is_transferred {
		return utils.BadRequestResponse(c, ticket_transfer)
	}
	return utils.NoDataSuccessResponse(c, "Success")

}

// VALIDATE TICKETS FOR ENTRY
