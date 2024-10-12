package admin_routes

import (
	"fmt"
	"os"

	"txrnxp/utils"
	"txrnxp/views/wallets"
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
	routes.Get("/config", auth_utils.ValidateAuth, admin_views.GetAdminCommissionConfig)
	routes.Post("/config", auth_utils.ValidateAuth, admin_views.CreateAdminCommissionConfig)
	routes.Patch("/config/:id", auth_utils.ValidateAuth, admin_views.UpdateAdminCommissionConfig)
	routes.Get("/wallet", auth_utils.ValidateAuth, wallets.GetAdminWallet)
	routes.Get(
		"wallet/entries",
		auth_utils.ValidateAuth,
		utils.ValidateRequestLimitAndPage,
		utils.ValidateRequestFilters(func() string {
			return "admin_wallet_tx"
		}),
		wallets.GetAdminWalletTransactions,
	)
	// routes.Get("admin/admin", auth_utils.ValidateAuth, wallets.CreateAdminWallet)

	_ = routes
}
