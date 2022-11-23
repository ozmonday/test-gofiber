package routes

import (
	"testfiber/api/handler"
	"testfiber/pkg/todo"

	"github.com/gofiber/fiber/v2"
)

func TodoRouter(app fiber.Router, service todo.Service) {
	app.Post("/todo-items", handler.AddTodo(service))
	app.Get("/todo-items", handler.GetTodos(service))
	app.Patch("/todo-items/:id", handler.EditTodo(service))
	app.Delete("/todo-items/:id", handler.DeleteTodo(service))
	app.Get("/todo-items/:id", handler.GetTodo(service))
}
