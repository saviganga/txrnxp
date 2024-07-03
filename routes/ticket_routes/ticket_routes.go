package ticket_routes

import (
	"fmt"
	"os"

	"txrnxp/utils/auth_utils"
	"txrnxp/views/ticket_views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {

	version := os.Getenv("VERSION")
	pathPrefix := fmt.Sprintf("/api/%v/tickets/", version)
	routes := app.Group(pathPrefix, logger.New())

	routes.Get("events/", auth_utils.ValidateAuth, ticket_views.GetEventTickets)
	routes.Post("events/", auth_utils.ValidateAuth, ticket_views.CreateEventTicket)
	routes.Get("users/", auth_utils.ValidateAuth, ticket_views.GetUserTickets)

	_ = routes
}
