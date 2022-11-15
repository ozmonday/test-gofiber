package main

import (
	"database/sql"
	"testfiber/api/routes"
	"testfiber/storage/activitiy"

	"github.com/gofiber/fiber/v2"
)

func connect() *sql.DB {
	conn, _ := sql.Open("mysql", "")
	return conn
}

func main() {
	app := fiber.New()
	conn := connect()
	repo := activitiy.NewRepository(conn)
	service := activitiy.NewService(repo)

	routes.ActivityRouter(app, service)

	app.Listen(":3000")
}
