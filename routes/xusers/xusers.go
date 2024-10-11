package xusers

import (
	"fmt"
	"os"

	"txrnxp/utils"
	"txrnxp/utils/auth_utils"
	"txrnxp/views/xusers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {

	version := os.Getenv("VERSION")
	pathPrefix := fmt.Sprintf("/api/%v/users/", version)
	routes := app.Group(pathPrefix, logger.New())

	routes.Get(
		"",
		auth_utils.ValidateAuth,
		utils.ValidateRequestLimitAndPage,
		utils.ValidateRequestFilters(func() string {
			return "xuser"
		}),
		xusers.GetUsers,
	)
	routes.Post("", xusers.CreateUsers)
	routes.Post("upload-image/", auth_utils.ValidateAuth, xusers.UploadUserImage)

	_ = routes
}
