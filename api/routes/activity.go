package routes

import (
	"testfiber/api/handler"

	"github.com/gofiber/fiber/v2"
)

func ActivityRouter(app fiber.Router) {
	app.Get("/activity-groups", handler.GetActivities())
	app.Get("/activity-groups/:id", handler.GetActivity())
	app.Post("/activity-groups", handler.AddActivity())
	app.Patch("/activity-groups/:id", handler.EditActivity())
	app.Delete("/activity-groups/:id", handler.DeleteActivity())
}
