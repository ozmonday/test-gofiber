package main

import (
	"testfiber/api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.ActivityRouter(app)

	app.Listen(":3000")
}
