package ticket_views

import (
	"txrnxp/utils/ticket_utils"

	"github.com/gofiber/fiber/v2"
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


// func GetEventTickets(c *fiber.Ctx) error {
// 	db := initialisers.ConnectDb().Db
// 	event_tickets := []models.EventTicket{}
// 	db.Find(&event_tickets)
// 	return c.Status(200).JSON(event_tickets)
// }
