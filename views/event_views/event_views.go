package event_views

import (
	"txrnxp/utils/event_utils"

	"github.com/gofiber/fiber/v2"
)

func GetEvents(c *fiber.Ctx) error {
	return c.SendString("Welcome to the United States of Ganga events nigguuhhhhh!")
}

func CreateEvents(c *fiber.Ctx) error {

	event, err := event_utils.CreateEvent(c)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(event)
}
