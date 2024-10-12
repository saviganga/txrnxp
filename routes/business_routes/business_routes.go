package business_routes

import (
	"fmt"
	"os"
	"txrnxp/utils"
	"txrnxp/utils/auth_utils"
	"txrnxp/validators/business_validators"
	"txrnxp/views/business_views"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {

	version := os.Getenv("VERSION")
	pathPrefix := fmt.Sprintf("/api/%v/business/", version)
	routes := app.Group(pathPrefix, logger.New())

	routes.Get(
		"",
		auth_utils.ValidateAuth,
		utils.ValidateRequestLimitAndPage,
		utils.ValidateRequestFilters(func() string {
			return "business"
		}),
		business_views.GetBusiness,
	)
	routes.Get(
		":id/",
		auth_utils.ValidateAuth,
		business_views.GetBusinessById,
	)
	routes.Post("", auth_utils.ValidateAuth, business_views.CreateBusiness)
	routes.Post(":id/upload-image/", auth_utils.ValidateAuth, business_validators.ValidateBusinessOwner, business_views.UploadBusinessImage)

	_ = routes
}
