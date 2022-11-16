package handler

import (
	"errors"
	"net/http"
	"strconv"
	"testfiber/api/payload"
	"testfiber/storage/entities"
	"testfiber/storage/todo"

	"github.com/gofiber/fiber/v2"
)

func AddTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todo entities.Todo
		if err := c.BodyParser(&todo); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		if todo.Title == "" {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, errors.New("title cannot be null")))
		}

		if todo.ActivityID == 0 {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, errors.New("activity_group_id cannot be null")))
		}

		if err := service.Repo.Create(&todo); err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(payload.ErrorResponse(http.StatusInternalServerError, err))
		}

		c.Status(http.StatusCreated)
		return c.JSON(payload.SuccessResponse(todo.Map()))
	}
}

func EditTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todo entities.Todo

		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		todo.ID = int64(id)

		if err := c.BodyParser(&todo); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(payload.ErrorResponse(http.StatusBadRequest, err))
		}

		if err := service.Repo.Update(&todo); err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(todo.Map()))
	}
}

func DeleteTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		where := map[string]string{"id": c.Params("id")}

		if err := service.Repo.Delete(where); err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(&fiber.Map{}))
	}
}

func GetTodo(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		where := make(map[string]string)

		todo, err := service.Repo.Read(where)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SuccessResponse(todo.Map()))
	}
}

func GetTodos(service todo.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var todos []fiber.Map
		where := map[string]string{}
		where["activity_group_id"] = c.Query("activity_group_id")

		rows, err := service.Repo.Reads(where)
		if err != nil {
			c.Status(http.StatusNotFound)
			return c.JSON(payload.ErrorResponse(http.StatusNotFound, err))
		}

		for _, v := range *rows {
			todos = append(todos, *v.Map())
		}

		c.Status(http.StatusOK)
		return c.JSON(payload.SliceSuccessResponse(&todos))
	}
}
