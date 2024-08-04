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

func AdminTopupWallet(c *fiber.Ctx) error {
	is_manual_entry, manual_entry := wallets_utils.AdminWalletManualEntry(c)
	if !is_manual_entry {
		return c.Status(400).JSON(fiber.Map{
			"message": manual_entry,
		})
	}
	return c.Status(200).JSON(manual_entry)
}

func GetUserWalletTransactions(c *fiber.Ctx) error {
	return wallets_utils.GetUserWalletTransactions(c)
}


func WalletTransfer(c *fiber.Ctx) error {
	is_transferred, wallet_transfer := wallets_utils.WalletTransfer(c)
	if !is_transferred {
		return c.Status(400).JSON(fiber.Map{
			"message": wallet_transfer,
		})
	}
	return c.Status(200).JSON(wallet_transfer)
}