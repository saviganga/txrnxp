package ticket_views

import (
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
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(201).JSON(event)
}

func GetEventTickets(c *fiber.Ctx) error {
	authenticated_user := c.Locals("user").(jwt.MapClaims)
	entity := c.Get("Entity")
	user_id := authenticated_user["id"].(string)

	event_tickets, err := ticket_utils.GetEventTickets(user_id, entity)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err,
		})
	}

	return c.Status(200).JSON(event_tickets)

}
