package admin_routes

import (
	"fmt"
	"os"

	"txrnxp/utils/auth_utils"
	"txrnxp/views/admin_views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {

	version := os.Getenv("VERSION")
	pathPrefix := fmt.Sprintf("/api/%v/admin/", version)
	routes := app.Group(pathPrefix, logger.New())

	routes.Get("", auth_utils.ValidateAuth, admin_views.GetAdminUsers)
	routes.Post("", admin_views.CreateAdminUsers)

	_ = routes
}
