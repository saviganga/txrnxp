package event_routes

import (
	"fmt"
	"os"
	"txrnxp/utils"
	"txrnxp/utils/auth_utils"
	"txrnxp/validators/event_validators"
	"txrnxp/views/event_views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {

	version := os.Getenv("VERSION")
	pathPrefix := fmt.Sprintf("/api/%v/events/", version)
	routes := app.Group(pathPrefix, logger.New())

	routes.Get(
		"",
		utils.ValidateRequestLimitAndPage,
		utils.ValidateRequestFilters(func() string {
			return "event"
		}),
		event_views.GetEvents,
	)
	routes.Get(
		"my-events/",
		auth_utils.ValidateAuth,
		utils.ValidateRequestLimitAndPage,
		utils.ValidateRequestFilters(func() string {
			return "event"
		}),
		event_views.GetMyEvents,
	)
	routes.Get(
		"history/",
		auth_utils.ValidateAuth,
		utils.ValidateRequestLimitAndPage,
		utils.ValidateRequestFilters(func() string {
			return "event"
		}),
		event_views.GetEventHistory,
	)
	// routes.Get(
	// 	":id/",
	// 	auth_utils.ValidateAuth,
	// 	event_views.GetEventById,
	// )
	routes.Post("", auth_utils.ValidateAuth, event_views.CreateEvents)
	routes.Patch(":id/", auth_utils.ValidateAuth, event_validators.ValidateEventOrganiser, event_validators.ValidateUpdateEventRequestBody, event_views.UpdateEvent)
	routes.Get(":reference/", event_views.GetEventByReference)
	routes.Get(
		":id/tickets/",
		auth_utils.ValidateAuth,
		event_validators.ValidateEventOrganiser,
		utils.ValidateRequestLimitAndPage,
		utils.ValidateRequestFilters(func() string {
			return "event_ticket"
		}),
		event_views.EventTickets,
	)
	routes.Post(":id/upload-image/", auth_utils.ValidateAuth, event_validators.ValidateEventOrganiser, event_views.UploadEventImage)

	_ = routes
}
