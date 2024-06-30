package event_routes

import (
	"fmt"
	"os"

	"txrnxp/utils/auth_utils"
	"txrnxp/views/event_views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {

	version := os.Getenv("VERSION")
	pathPrefix := fmt.Sprintf("/api/%v/events/", version)
	routes := app.Group(pathPrefix, logger.New())

	routes.Get("", event_views.GetEvents)
	routes.Post("", auth_utils.ValidateAuth, event_views.CreateEvents)

	_ = routes
}
