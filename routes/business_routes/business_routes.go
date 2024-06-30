package business_routes

import (
	"fmt"
	"os"

	"txrnxp/views/business_views"
	"txrnxp/utils/auth_utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {

	version := os.Getenv("VERSION")
	pathPrefix := fmt.Sprintf("/api/%v/business/", version)
	routes := app.Group(pathPrefix, logger.New())

	routes.Get("", auth_utils.ValidateAuth, business_views.GetBusiness)
	routes.Post("", auth_utils.ValidateAuth, business_views.CreateBusiness)

	_ = routes
}
