package xusers

import (
	"fmt"
	"os"

	"txrnxp/views/xusers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {

	version := os.Getenv("VERSION")
	pathPrefix := fmt.Sprintf("/api/%v/users/", version)
	routes := app.Group(pathPrefix, logger.New())

	routes.Get("", xusers.Home)

	_ = routes
}
