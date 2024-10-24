package ticket_routes

import (
	"fmt"
	"os"

	"txrnxp/utils"
	"txrnxp/utils/auth_utils"
	"txrnxp/views/ticket_views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {

	version := os.Getenv("VERSION")
	pathPrefix := fmt.Sprintf("/api/%v/tickets/", version)
	routes := app.Group(pathPrefix, logger.New())

	// only for admins now
	routes.Get(
		"events/",
		auth_utils.ValidateAuth,
		utils.ValidateRequestLimitAndPage,
		utils.ValidateRequestFilters(func() string {
			return "event_ticket"
		}),
		ticket_views.GetEventTickets,
	)

	routes.Get(
		"events/:id",
		auth_utils.ValidateAuth,
		utils.ValidateRequestLimitAndPage,
		utils.ValidateRequestFilters(func() string {
			return "event_ticket"
		}),
		ticket_views.GetEventTicketById,
	)
	routes.Post("events/", auth_utils.ValidateAuth, ticket_views.CreateEventTicket)
	routes.Get(
		"users/",
		auth_utils.ValidateAuth,
		utils.ValidateRequestLimitAndPage,
		utils.ValidateRequestFilters(func() string {
			return "user_ticket"
		}),
		ticket_views.GetUserTickets,
	)
	routes.Get("users/:reference/", ticket_views.GetUserTicketByReference)
	routes.Post("users/:reference/", auth_utils.ValidateAuth, ticket_views.ValidateUserTicket)
	routes.Post("buy/wallet/", auth_utils.ValidateAuth, ticket_views.CreateUserTicket)
	routes.Post("transfer/", auth_utils.ValidateAuth, ticket_views.TransferUserTicket)

	_ = routes
}
