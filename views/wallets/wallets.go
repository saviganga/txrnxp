package wallets

import (
	"txrnxp/utils/wallets_utils"

	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.SendString("Welcome to the United States of Ganga wallets nigguhhhhh!")
}

func GetWallets(c *fiber.Ctx) error {
	return wallets_utils.GetUserWallets(c)
}
