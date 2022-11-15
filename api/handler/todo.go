package handler

import (
	"net/http"
	"testfiber/api/payload"
	"testfiber/storage/entities"
	"testfiber/storage/todo"

	"github.com/gofiber/fiber/v2"
)

func AddTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todo entities.Todo

		c.Status(http.StatusCreated)
		return c.JSON(payload.SuccessResponse(todo.Map()))
	}
}

func EditTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todo entities.Todo

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(todo.Map()))
	}
}

func DeleteTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(&fiber.Map{}))
	}
}

func GetTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		where := make(map[string]string)

		todo, err := service.Repo.Read(where)
		if err != nil {
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(todo.Map()))
	}
}

func GetTodos(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todos []fiber.Map

		c.Status(http.StatusOK)
		return c.JSON(payload.SliceSuccessResponse(&todos))
	}
}
