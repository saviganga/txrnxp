package event_views

import (
	"txrnxp/initialisers"
	"txrnxp/models"
	"txrnxp/utils/event_utils"

	"github.com/gofiber/fiber/v2"
)

func GetEvents(c *fiber.Ctx) error {
	// find a way to get the event tickets
	db := initialisers.ConnectDb().Db
	events := []models.Event{}
	db.Find(&events)
	return c.Status(200).JSON(events)
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



func GetEventByReference(c *fiber.Ctx) error {

	event, err := event_utils.GetEventByReference(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(event)
}