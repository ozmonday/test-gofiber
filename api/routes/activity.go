package routes

import (
	"testfiber/api/handler"
	"testfiber/pkg/activity"

	"github.com/gofiber/fiber/v2"
)

func ActivityRouter(app fiber.Router, service activity.Service) {
	app.Get("/activity-groups", handler.GetActivities(service))
	app.Get("/activity-groups/:id", handler.GetActivity(service))
	app.Post("/activity-groups", handler.AddActivity(service))
	app.Patch("/activity-groups/:id", handler.EditActivity(service))
	app.Delete("/activity-groups/:id", handler.DeleteActivity(service))
}
