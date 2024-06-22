package home

import (
	"fmt"
	"os"

	"txrnxp/views/home"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {

	version := os.Getenv("VERSION")
	pathPrefix := fmt.Sprintf("/api/%v/welcome/", version)
	routes := app.Group(pathPrefix, logger.New())

	routes.Get("", home.Home)

	_ = routes
}
