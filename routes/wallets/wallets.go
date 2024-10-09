package wallets

import (
	"fmt"
	"os"
	"txrnxp/utils"
	"txrnxp/utils/auth_utils"
	"txrnxp/views/wallets"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Routes(app *fiber.App) {

	version := os.Getenv("VERSION")
	pathPrefix := fmt.Sprintf("/api/%v/wallets/", version)
	routes := app.Group(pathPrefix, logger.New())
	routes.Get(
		"",
		auth_utils.ValidateAuth,
		utils.ValidateRequestLimitAndPage,
		utils.ValidateRequestFilters(func() string {
			return "wallets"
		}),
		wallets.GetWallets,
	)
	routes.Get("/entries", auth_utils.ValidateAuth, wallets.GetUserWalletTransactions)
	routes.Post("/topup/admin", auth_utils.ValidateAuth, wallets.AdminTopupWallet)
	routes.Post("/transfer", auth_utils.ValidateAuth, wallets.WalletTransfer)

	_ = routes
}
