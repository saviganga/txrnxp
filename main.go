package main

import (
	"os"
	"txrnxp/home"
	"txrnxp/initialisers"
	"txrnxp/wallets"
	"txrnxp/xusers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	initialisers.LoadEnv()
	initialisers.ConnectDb()
}

func main() {

	port := os.Getenv("PORT")

	app := fiber.New()
	app.Use(logger.New())

	home.Routes(app)
	xusers.Routes(app)
	wallets.Routes(app)

	app.Listen(":" + port)
}
