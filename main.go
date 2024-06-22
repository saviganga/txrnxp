package main

import (
	"os"
	"txrnxp/home"
	"txrnxp/initialisers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	initialisers.LoadEnv()
}

func main() {

	port := os.Getenv("PORT")

	app := fiber.New()
	app.Use(logger.New())

	home.Routes(app)

	app.Listen(":" + port)
}
