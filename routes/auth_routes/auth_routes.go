package auth_routes

import (
	"fmt"
	"os"

	"txrnxp/views/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {

	version := os.Getenv("VERSION")
	pathPrefix := fmt.Sprintf("/api/%v/auth/", version)
	routes := app.Group(pathPrefix, logger.New())

	routes.Get("", auth.Home)
	routes.Post("/login", auth.Login)

	_ = routes
}
